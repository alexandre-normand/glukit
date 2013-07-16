// The glukit package is the main package for the application. This is where it all starts.
package glukit

import (
	"app/engine"
	"app/importer"
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
	"appengine/urlfetch"
	"appengine/user"
	"bufio"
	"fmt"
	"html/template"
	"lib/drive"
	"lib/goauth2/oauth"
	"net/http"
	"os"
	"time"
)

const (
	DEMO_EMAIL = "demo@glukit.com"
)

// config returns the configuration information for OAuth and Drive.
func config() *oauth.Config {
	host, clientId, clientSecret := getEnvSettings()
	config := oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile " + drive.DriveReadonlyScope,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		AccessType:   "offline",
		RedirectURL:  fmt.Sprintf("http://%s/oauth2callback", host),
	}

	return &config
}

type DataSeries struct {
	Name string            `json:"name"`
	Data []model.DataPoint `json:"data"`
	Type string            `json:"type"`
}

type RenderVariables struct {
	DataPath     string
	ChannelToken string
}

var graphTemplate = template.Must(template.ParseFiles("view/templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("view/templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("view/templates/nodata.html"))

var processFile = delay.Func("processSingleFile", processSingleFile)
var processDemoFile = delay.Func("processDemoFile", processStaticDemoFile)
var refreshUserData = delay.Func("refreshUserData", func(context appengine.Context, userEmail string, autoScheduleNextRun bool) {
	context.Criticalf("This function purely exists as a workaround to the \"initialization loop\" error that ",
		"shows up because the function calls itself. This implementation defines the same signature as the ",
		"real one which we define in init() to override this implementation!")
})

var emptyDataPointSlice []model.DataPoint

// handleLoggedInUser is responsible for directing the user to the graph page after optionally:
//   1. Storing the GlukitUser entry if it's the first access
//   2. Refreshing the glukit oauth token
//   3. Kick off the processing of background import of files
//
// TODO: This is a big function, this should be split up into smaller ones
func handleLoggedInUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	transport := new(oauth.Transport)
	var oauthToken oauth.Token
	glukitUser, _, _, _, err := store.GetUserData(context, user.Email)
	scheduleAutoRefresh := false
	if err == datastore.ErrNoSuchEntity {
		oauthToken, transport = getOauthToken(request)

		context.Infof("No data found for user [%s], creating it", user.Email)
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		// We store the refresh token separately from the rest. This token is long-lived, meaning that if we have a glukit user with no refresh token, we need to force getting a new one (which is to be avoided)
		_, err = store.StoreUserProfile(context, time.Now(), model.GlukitUser{user.Email, "", "", time.Now(), "", "", util.BEGINNING_OF_TIME, util.BEGINNING_OF_TIME, oauthToken, oauthToken.RefreshToken, model.UNDEFINED_SCORE})
		if err != nil {
			util.Propagate(err)
		}
		// We only schedule the auto refresh on first access since all subsequent runs of scheduled tasks will also
		// reschedule themselve a new run
		scheduleAutoRefresh = true
	} else if err != nil {
		util.Propagate(err)
	} else {
		oauthToken = glukitUser.Token

		context.Debugf("Initializing transport from token [%s]", oauthToken)
		transport = &oauth.Transport{
			Config: config(),
			Transport: &urlfetch.Transport{
				Context: context,
			},
			Token: &oauthToken,
		}

		if !oauthToken.Expired() && len(glukitUser.RefreshToken) > 0 {
			context.Debugf("Token [%s] still valid, reusing it...", oauthToken)
		} else {
			if oauthToken.Expired() {
				context.Infof("Token [%v] expired on [%s], refreshing with refresh token [%s]...", oauthToken, oauthToken.Expiry, glukitUser.RefreshToken)
			} else if len(glukitUser.RefreshToken) == 0 {
				context.Warningf("No refresh token stored, getting a new one and saving it...")
			}

			// We lost the refresh token, we need to force approval and get a new one
			if len(glukitUser.RefreshToken) == 0 {
				context.Criticalf("We lost the refresh token for user [%s], getting a new one with the force approval.", user.Email)
				oauthToken, transport = getOauthToken(request)
				glukitUser.RefreshToken = oauthToken.RefreshToken
				store.StoreUserProfile(context, time.Now(), *glukitUser)
			} else {
				transport.Token.RefreshToken = glukitUser.RefreshToken

				err := transport.Refresh(context)
				util.Propagate(err)

				context.Debugf("Storing new refreshed token [%s] in datastore...", oauthToken)
				glukitUser.LastUpdated = time.Now()
				glukitUser.Token = oauthToken
				_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
				if err != nil {
					util.Propagate(err)
				}
			}

		}
	}

	task, err := refreshUserData.Task(user.Email, scheduleAutoRefresh)
	if err != nil {
		context.Criticalf("Could not schedule execution of the data refresh for user [%s]: %v", user.Email, err)
	}
	taskqueue.Add(context, task, "refresh")
	context.Infof("Kicked off data update for user [%s]...", user.Email)

	// Render the graph view, it might take some time to show something but it will as soon as a file import
	// completes
	renderRealUser(writer, request)
}

// Async task that searches on Google Drive for dexcom files. It handles some high watermark of the last import
// to avoid downloading already imported files (unless they've been updated). It also schedules itself to run again
// the next day unless the token is invalid.
func updateUserData(context appengine.Context, userEmail string, autoScheduleNextRun bool) {
	glukitUser, userProfileKey, _, _, err := store.GetUserData(context, userEmail)
	if err != nil {
		context.Errorf("We're trying to run an update data task for user [%s] that doesn't exist. Got error: %v", userEmail, err)
		return
	}

	transport := &oauth.Transport{
		Config: config(),
		Transport: &urlfetch.Transport{
			Context: context,
		},
		Token: &glukitUser.Token,
	}

	// If the token is expired, try to get a fresh one by doing a refresh (which should use the refresh_token
	if glukitUser.Token.Expired() {
		transport.Token.RefreshToken = glukitUser.RefreshToken
		err := transport.Refresh(context)
		if err != nil {
			context.Errorf("Error updating token for user [%s], let's hope he comes back soon so we can get a fresh token: %v", userEmail, err)
			return
		}

		// Update the user with the new token
		context.Infof("Token refreshed, updating user [%s] with token [%v]", userEmail, glukitUser.Token)
		store.StoreUserProfile(context, time.Now(), *glukitUser)
	}

	nextUpdate := time.Now().AddDate(0, 0, 1)
	files, err := importer.SearchDataFiles(transport.Client(), glukitUser.LastUpdated)
	if err != nil {
		context.Warningf("Error while searching for files on google drive for user [%s]: %v", userEmail, err)
	} else {
		switch {
		case len(files) == 0:
			context.Infof("No new or updated data found for existing user [%s]", userEmail)
		case len(files) > 0:
			context.Infof("Found new data files for user [%s], downloading and storing...", userEmail)
			processData(&glukitUser.Token, files, context, userEmail, userProfileKey)
		}
	}

	if autoScheduleNextRun {
		task, err := refreshUserData.Task(userEmail, autoScheduleNextRun)
		if err != nil {
			context.Criticalf("Couldn't schedule the next execution of the data refresh for user [%s]. This breaks background updating of user data!: %v", userEmail, err)
		}
		task.ETA = nextUpdate
		taskqueue.Add(context, task, "refresh")

		context.Infof("Scheduled next data update for user [%s] at [%s]", userEmail, nextUpdate.Format(util.TIMEFORMAT))
	} else {
		context.Infof("Not scheduling a the next refresh as requested by autoScheduleNextRun [%t]", autoScheduleNextRun)
	}
}

func getOauthToken(request *http.Request) (oauthToken oauth.Token, transport *oauth.Transport) {
	context := appengine.NewContext(request)

	// Exchange code for an access token at OAuth provider.
	code := request.FormValue("code")
	configuration := config()
	context.Debugf("Getting token with configuration [%v]...", configuration)

	t := &oauth.Transport{
		Config: configuration,
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(request),
		},
	}

	token, err := t.Exchange(context, code)
	util.Propagate(err)

	context.Infof("Got brand new oauth token [%v] with refresh token [%s]", token, token.RefreshToken)
	return *token, t
}

func getEnvSettings() (host, clientId, clientSecret string) {
	if appengine.IsDevAppServer() {
		host = "localhost:8080"
		clientId = "***REMOVED***"
		clientSecret = "***REMOVED***"

	} else {
		host = "glukit.appspot.com"
		clientId = "414109645872-d6igmhnu0loafu53uphf8j67ou8ngjiu.apps.googleusercontent.com"
		clientSecret = "IYbOW0Aha34xMqTaPVO-_ar5"
	}

	return host, clientId, clientSecret
}

func processData(token *oauth.Token, files []*drive.File, context appengine.Context, userEmail string, userProfileKey *datastore.Key) {
	// TODO : Look at recent file import log for that file and skip to the new data. It would be nice to be able to
	// use the Http Range header but that's unlikely to be possible since new event/read data is spreadout in the
	// file
	for i := range files {
		task, err := processFile.Task(token, files[i], userEmail, userProfileKey)
		if err != nil {
			util.Propagate(err)
		}
		taskqueue.Add(context, task, "store")
	}
}

func processSingleFile(context appengine.Context, token *oauth.Token, file *drive.File, userEmail string, userProfileKey *datastore.Key) {
	t := &oauth.Transport{
		Config: config(),
		Transport: &urlfetch.Transport{
			Context: context,
		},
		Token: token,
	}

	reader, err := importer.GetFileReader(context, t, file)
	if err != nil {
		context.Infof("Error reading file %s, skipping...", file.OriginalFilename)
	} else {
		// Default to beginning of time
		startTime := util.BEGINNING_OF_TIME
		if lastFileImportLog, err := store.GetFileImportLog(context, userProfileKey, file.Id); err == nil {
			startTime = lastFileImportLog.LastDataProcessed
			context.Infof("Reloading data from file [%s]-[%s] starting at date [%s]...", file.Id, file.OriginalFilename, startTime.Format(util.TIMEFORMAT))
		} else if err == datastore.ErrNoSuchEntity {
			context.Debugf("First import of file [%s]-[%s]...", file.Id, file.OriginalFilename)
		} else if err != nil {
			util.Propagate(err)
		}

		lastReadTime := importer.ParseContent(context, reader, 500, userProfileKey, startTime, store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
		store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: file.Id, Md5Checksum: file.Md5Checksum, LastDataProcessed: lastReadTime})
		reader.Close()

		if glukitUser, err := store.GetUserProfile(context, userProfileKey); err != nil {
			context.Warningf("Error getting retrieving GlukitUser [%s], this needs attention: [%v]", userEmail, err)
		} else if glukitUser.Score.UpperBound.Before(lastReadTime) {
			// Calculate Glukit Score here from the last 2 weeks of data
			glukitScore, err := engine.CalculateGlukitScore(context, glukitUser)
			if err != nil {
				context.Warningf("Error calculating a new GlukitScore for [%s], this needs attention: [%v]", userEmail, err)
			} else {
				// Store the updated GlukitScore
				context.Infof("New GlukitScore calculated for user [%s]: [%d]", userEmail, glukitScore.Value)
				glukitUser.Score = *glukitScore
				if _, err := store.StoreUserProfile(context, time.Now(), *glukitUser); err != nil {
					context.Errorf("Error persisting glukit score [%v] for user [%s]: %v", glukitScore, userEmail, err)
				}
			}
		} else {
			context.Debugf("Skipping recalculation of GlukitScore because the last calculation was with a more recent read [%s] than the most recent read from this last batch [%s]",
				glukitUser.Score.UpperBound.Format(util.TIMEFORMAT), lastReadTime.Format(util.TIMEFORMAT))
		}
	}
	channel.Send(context, userEmail, "Refresh")
}

func processStaticDemoFile(context appengine.Context, userProfileKey *datastore.Key) {

	// open input file
	fi, err := os.Open("data.xml")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if fi.Close() != nil {
			panic(err)
		}
	}()
	// make a read buffer
	reader := bufio.NewReader(fi)

	lastReadTime := importer.ParseContent(context, reader, 500, userProfileKey, util.BEGINNING_OF_TIME, store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
	store.LogFileImport(context, userProfileKey, model.FileImportLog{Id: "demo", Md5Checksum: "dummychecksum", LastDataProcessed: lastReadTime})
	channel.Send(context, DEMO_EMAIL, "Refresh")
}

// buildPerfectBaseline generates an array of reads that represents the target/perfection
func buildPerfectBaseline(glucoseReads []model.GlucoseRead) (reads []model.GlucoseRead) {
	reads = make([]model.GlucoseRead, len(glucoseReads))
	for i := range glucoseReads {
		reads[i] = model.GlucoseRead{glucoseReads[i].LocalTime, glucoseReads[i].Timestamp, 83}
	}

	return reads
}
