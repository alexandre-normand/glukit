package glukit

import (
	"parser"
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"html/template"
	"goauth2/oauth"
	"appengine"
	"appengine/user"
	"appengine/urlfetch"
	"appengine/datastore"
	"appengine/channel"
	"appengine/taskqueue"
	"models"
	"drive"
	"store"
	"fetcher"
	"sort"
	"datautils"
	"bufio"
	"os"
	"timeutils"
	"sysutils"
	"appengine/delay"
	stat "github.com/grd/stat"
)

const (
	DEMO_EMAIL = "demo@glukit.com"
)

// config returns the configuration information for OAuth and Drive.
func config() *oauth.Config {
	host, clientId, clientSecret := getEnvSettings()
	return &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile " + drive.DriveReadonlyScope,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		AccessType:   "offline",
		RedirectURL:  fmt.Sprintf("http://%s/oauth2callback", host),
	}
}

type DataSeries struct {
	Name        string              `json:"name"`
	Data        []models.DataPoint  `json:"data"`
	Type        string              `json:"type"`
}

type RenderVariables struct {
	DataPath     string
	ChannelToken string
}

var graphTemplate = template.Must(template.ParseFiles("templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("templates/nodata.html"))

var processFile = delay.Func("processSingleFile", processSingleFile)
var processDemoFile = delay.Func("processDemoFile", processStaticDemoFile)
var refreshUserData = delay.Func("refreshUserData", func(context appengine.Context, userEmail string) {
		context.Criticalf("This function purely exists as a workaround to the \"initialization loop\" error that ",
				"shows up because the function calls itself. This implementation defines the same signature as the ",
				"real one which we define in init() to override this implementation!")
	})

var emptyDataPointSlice []models.DataPoint

func init() {
	http.HandleFunc("/json.demo", demoContent)
	http.HandleFunc("/json", content)
	http.HandleFunc("/json.demo.tracking", demoTracking)
	http.HandleFunc("/json.tracking", tracking)

	http.HandleFunc("/demo", renderDemo)
	http.HandleFunc("/graph", renderRealUser)

	http.HandleFunc("/", landing)
	http.HandleFunc("/nodata", nodata)
	http.HandleFunc("/realuser", handleRealUser)
	http.HandleFunc("/oauth2callback", callback)

	refreshUserData = delay.Func("refreshUserData", updateUserData)
}

func callback(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	t := new(oauth.Transport)
	var oauthToken oauth.Token
	glukitUser, _, _, _, err := store.GetUserData(context, user.Email)
	if err == datastore.ErrNoSuchEntity {
		oauthToken, t = getOauthToken(request)

		context.Infof("No data found for user [%s], creating it", user.Email)
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		_, err = store.StoreUserProfile(context, time.Now(), models.GlukitUser{user.Email, "", "", time.Now(), "", "", timeutils.BEGINNING_OF_TIME, timeutils.BEGINNING_OF_TIME, oauthToken, models.UNDEFINED_SCORE})
		if err != nil {
			sysutils.Propagate(err)
		}
	} else if err != nil {
		sysutils.Propagate(err)
	} else {
		oauthToken = glukitUser.Token

		context.Debugf("Initializing transport from token [%s]", oauthToken)
		t = &oauth.Transport{
			Config: config(),
			Transport: &urlfetch.Transport{
				Context: context,
			},
			Token: &oauthToken,
		}

		if !oauthToken.Expired() {
			context.Debugf("Token [%s] still valid, reusing it...", oauthToken)
		} else {
			context.Infof("Token expired on [%s], refreshing...", oauthToken.Expiry)
			err := t.Refresh()
			sysutils.Propagate(err)

			context.Debugf("Storing new refreshed token [%s] in datastore...", oauthToken)
			glukitUser.LastUpdated = time.Now()
			glukitUser.Token = oauthToken
			_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
			if err != nil {
				sysutils.Propagate(err)
			}
		}
	}

	context.Debugf("Found existing user: %s", user.Email)

	task, err := refreshUserData.Task(user.Email)
	if err != nil {
		context.Criticalf("Couldn't schedule execution of the data refresh for user [%s]: %v", user.Email, err)
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
func updateUserData(context appengine.Context, userEmail string) {
	glukitUser, key, _, _, err := store.GetUserData(context, userEmail)
	if err != nil {
		context.Errorf("We're trying to run an update data task for user [%s] that doesn't exist. Got error: %v", userEmail, err)
		return;
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
		err := transport.Refresh()
		if err != nil {
			context.Errorf("Error updating token for user [%s], let's hope he comes back soon so we can get a fresh token: %v", userEmail, err)
			return;
		}

		// Update the user with the new token
		context.Infof("Token refreshed, updating user [%s] with token [%v]", userEmail, glukitUser.Token)
		store.StoreUserProfile(context, time.Now(), *glukitUser)
	}

	nextUpdate := time.Now().AddDate(0, 0, 1)
	files, err := fetcher.SearchDataFiles(transport.Client(), glukitUser.LastUpdated)
	if err != nil {
		context.Warningf("Error while searching for files on google drive for user [%s]: %v", userEmail, err)
	} else {
		switch {
		case len(files) == 0:
			context.Infof("No new or updated data found for existing user [%s]", userEmail)
		case len(files) > 0:
			context.Infof("Found new data files for user [%s], downloading and storing...", userEmail)
			processData(&glukitUser.Token, files, context, userEmail, key)
		}
	}

	task, err := refreshUserData.Task(userEmail)
	if err != nil {
		context.Criticalf("Couldn't schedule the next execution of the data refresh for user [%s]. This breaks background updating of user data!: %v", userEmail, err)
	}
	task.ETA = nextUpdate
	taskqueue.Add(context, task, "refresh")

	context.Infof("Scheduled next data update for user [%s] at [%s]", userEmail, nextUpdate.Format(timeutils.TIMEFORMAT))
}

func getOauthToken(request *http.Request) (oauthToken oauth.Token, transport *oauth.Transport) {
	// Exchange code for an access token at OAuth provider.
	code := request.FormValue("code")
	t := &oauth.Transport{
		Config: config(),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(request),
		},
	}

	token, err := t.Exchange(code)
	sysutils.Propagate(err)

	return *token, t
}

func getEnvSettings() (host, clientId, clientSecret string) {
	if (appengine.IsDevAppServer()) {
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
			sysutils.Propagate(err)
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

	reader, err := fetcher.GetFileReader(context, t, file)
	if err != nil {
		context.Infof("Error reading file %s, skipping...", file.OriginalFilename)
	} else {
		// Default to beginning of time
		startTime := timeutils.BEGINNING_OF_TIME
		if lastFileImportLog, err := store.GetFileImportLog(context, userProfileKey, file.Id); err == nil {
			startTime = lastFileImportLog.LastDataProcessed
			context.Infof("Reloading data from file [%s]-[%s] starting at date [%s]...", file.Id, file.OriginalFilename, startTime.Format(timeutils.TIMEFORMAT))
		} else if err == datastore.ErrNoSuchEntity {
			context.Debugf("First import of file [%s]-[%s]...", file.Id, file.OriginalFilename)
		} else if err != nil {
			sysutils.Propagate(err)
		}

		lastReadTime := parser.ParseContent(context, reader, 500, userProfileKey, startTime, store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
		store.LogFileImport(context, userProfileKey, models.FileImportLog{Id: file.Id, Md5Checksum: file.Md5Checksum, LastDataProcessed: lastReadTime})
		reader.Close()
	}
	channel.Send(context, userEmail, "Refresh")
}

func handleRealUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	glukitUser, _, _, _, err := store.GetUserData(context, user.Email)
	if err != nil {
		context.Infof("Redirecting [%s], glukitUser [%v] for authorization", user.Email, glukitUser)
		url := config().AuthCodeURL(request.URL.RawQuery)
		http.Redirect(writer, request, url, http.StatusFound)
	} else {
		context.Infof("User [%s] already exists, skipping authorization step...", user.Email)
		callback(writer, request)
	}
}

func processStaticDemoFile(context appengine.Context, userProfileKey *datastore.Key) {

	// open input file
	fi, err := os.Open("data.xml")
	if err != nil { panic(err) }
	// close fi on exit and check for its returned error
	defer func() {
		if fi.Close() != nil {
			panic(err)
		}
	}()
	// make a read buffer
	reader := bufio.NewReader(fi)

	lastReadTime := parser.ParseContent(context, reader, 500, userProfileKey, timeutils.BEGINNING_OF_TIME, store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreDaysOfExercises)
	store.LogFileImport(context, userProfileKey, models.FileImportLog{Id: "demo", Md5Checksum: "dummychecksum", LastDataProcessed: lastReadTime})
	channel.Send(context, DEMO_EMAIL, "Refresh")
}

func renderDemo(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	_, key, _, _, err := store.GetUserData(context, DEMO_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for demo user [%s], creating it", DEMO_EMAIL)
		dummyToken := oauth.Token{"", "", timeutils.BEGINNING_OF_TIME}
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(), models.GlukitUser{DEMO_EMAIL, "", "", time.Now(), "", "", time.Now(), timeutils.BEGINNING_OF_TIME, dummyToken, models.UNDEFINED_SCORE})
		if err != nil {
			sysutils.Propagate(err)
		}

		task, err := processDemoFile.Task(key)
		if err != nil {
			sysutils.Propagate(err)
		}
		taskqueue.Add(context, task, "store")
	} else if err != nil {
		sysutils.Propagate(err)
	} else {
		context.Infof("Data already stored for demo user [%s], continuing...", DEMO_EMAIL)
	}

	render(DEMO_EMAIL, "/json.demo", w, request)
}

func renderRealUser(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)
	render(user.Email, "/json", w, request)
}

func render(email string, datapath string, w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	token, err := channel.Create(context, email)
	if err != nil {
		context.Criticalf("Error creating channel for user [%s]", email)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderVariables := &RenderVariables{DataPath: datapath, ChannelToken: token}

	if err := graphTemplate.Execute(w, renderVariables); err != nil {
		context.Criticalf("Error executing template [%s]", graphTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func landing(w http.ResponseWriter, request *http.Request) {
	if err := landingTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nodata(w http.ResponseWriter, request *http.Request) {
	if err := nodataTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func demoContent(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	value := writer.Header()
	value.Add("Content-type", "application/json")

	_, _, lowerBound, upperBound, err := store.GetUserData(context, DEMO_EMAIL)
	if err != nil {
		sysutils.Propagate(err)
	}

	glucoseReads, err := store.GetUserReads(context, DEMO_EMAIL, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}
	injections, err := store.GetUserInjections(context, DEMO_EMAIL, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}
	carbs, err := store.GetUserCarbs(context, DEMO_EMAIL, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}

	context.Infof("Got %d reads", len(glucoseReads))
	writeUserJsonData(writer, glucoseReads, injections, carbs)
}

func writeUserJsonData(writer http.ResponseWriter, reads []models.GlucoseRead, injections []models.Injection, carbs []models.Carb) {
	enc := json.NewEncoder(writer)
	individuals := make([]DataSeries, 4)

	if (len(reads) == 0) {
		individuals[0] = DataSeries{"You", emptyDataPointSlice, "GlucoseReads"}
		individuals[1] = DataSeries{"You.Injection", emptyDataPointSlice, "Injections"}
		individuals[2] = DataSeries{"You.Carbohydrates", emptyDataPointSlice, "Carbs"}
		individuals[3] = DataSeries{"Perfection", emptyDataPointSlice, "ComparisonReads"}
	} else {
		individuals[0] = DataSeries{"You", models.GlucoseReadSlice(reads).ToDataPointSlice(), "GlucoseReads"}
		individuals[1] = DataSeries{"You.Injection", models.InjectionSlice(injections).ToDataPointSlice(reads), "Injections"}
		individuals[2] = DataSeries{"You.Carbohydrates", models.CarbSlice(carbs).ToDataPointSlice(reads), "Carbs"}
		individuals[3] = DataSeries{"Perfection", models.GlucoseReadSlice(buildPerfectBaseline(reads)).ToDataPointSlice(), "ComparisonReads"}
	}

	enc.Encode(individuals)
}

func content(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, _, lowerBound, upperBound, err := store.GetUserData(context, user.Email)
	if err != nil {
		sysutils.Propagate(err)
	}

	reads, err := store.GetUserReads(context, user.Email, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}
	injections, err := store.GetUserInjections(context, user.Email, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}
	carbs, err := store.GetUserCarbs(context, user.Email, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}

	value := writer.Header()
	value.Add("Content-type", "application/json")

	writeUserJsonData(writer, reads, injections, carbs)
}

func generateTrackingData(writer http.ResponseWriter, request *http.Request, reads []models.GlucoseRead) {
	value := writer.Header()
	value.Add("Content-type", "application/json")

	var trackingData models.TrackingData
	if (len(reads) > 0) {
		sort.Sort(models.ReadStatsSlice(reads))
		trackingData.Mean = stat.Mean(models.ReadStatsSlice(reads))
		trackingData.Deviation = stat.AbsdevMean(models.ReadStatsSlice(reads), 83)
		trackingData.Max, _ = stat.Max(models.ReadStatsSlice(reads))
		trackingData.Min, _ = stat.Min(models.ReadStatsSlice(reads))
		trackingData.Median = stat.MedianFromSortedData(models.ReadStatsSlice(reads))
		distribution := datautils.BuildHistogram(reads)
		sort.Sort(models.CoordinateSlice(distribution))
		trackingData.Distribution = distribution
	}

	enc := json.NewEncoder(writer)
	enc.Encode(trackingData)
}

func tracking(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, _, lowerBound, upperBound, err := store.GetUserData(context, user.Email)
	if err != nil {
		sysutils.Propagate(err)
	}

	reads, err := store.GetUserReads(context, user.Email, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}

	generateTrackingData(writer, request, reads)
}

func demoTracking(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	value := writer.Header()
	value.Add("Content-type", "application/json")

	_, _, lowerBound, upperBound, err := store.GetUserData(context, DEMO_EMAIL)
	if err != nil {
		sysutils.Propagate(err)
	}

	reads, err := store.GetUserReads(context, DEMO_EMAIL, lowerBound, upperBound)
	if err != nil {
		sysutils.Propagate(err)
	}

	generateTrackingData(writer, request, reads)
}


func buildPerfectBaseline(glucoseReads []models.GlucoseRead) (reads []models.GlucoseRead) {
	reads = make([]models.GlucoseRead, len(glucoseReads))
	for i := range glucoseReads {
		reads[i] = models.GlucoseRead{glucoseReads[i].LocalTime, glucoseReads[i].TimeValue, 83}
	}

	return reads
}
