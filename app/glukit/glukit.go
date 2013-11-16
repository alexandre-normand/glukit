// The glukit package is the main package for the application. This is where it all starts.
package glukit

import (
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"appengine/urlfetch"
	"appengine/user"
	"fmt"
	"lib/drive"
	"lib/goauth2/oauth"
	"lib/oauth2"
	"net/http"
	"time"
)

const (
	DEMO_EMAIL = "demo@glukit.com"
)

var emptyDataPointSlice []model.DataPoint

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
	glukitUser, _, _, err := store.GetUserData(context, user.Email)
	scheduleAutoRefresh := false
	if err == datastore.ErrNoSuchEntity {
		oauthToken, transport = getOauthToken(request)

		context.Infof("No data found for user [%s], creating it", user.Email)

		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when
		// this is ready
		// We store the refresh token separately from the rest. This token is long-lived, meaning that if
		// we have a glukit user with no refresh token, we need to force getting a new one (which is to be avoided)
		glukitUser = &model.GlukitUser{user.Email, "", "", time.Now(),
			"", "", util.GLUKIT_EPOCH_TIME, model.UNDEFINED_GLUCOSE_READ, oauthToken, oauthToken.RefreshToken,
			model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, "", time.Now()}
		_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
		if err != nil {
			util.Propagate(err)
		}
		// We only schedule the auto refresh on first access since all subsequent runs of scheduled tasks will also
		// reschedule themselve a new run
		scheduleAutoRefresh = true
	} else if _, ok := err.(store.StoreError); err != nil && !ok {
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
				context.Infof("Token [%v] expired on [%s], refreshing with refresh token [%s]...",
					oauthToken, oauthToken.Expiry, glukitUser.RefreshToken)
			} else if len(glukitUser.RefreshToken) == 0 {
				context.Warningf("No refresh token stored, getting a new one and saving it...")
			}

			// We lost the refresh token, we need to force approval and get a new one
			if len(glukitUser.RefreshToken) == 0 {
				context.Criticalf("We lost the refresh token for user [%s], getting a new one "+
					"with the force approval.", user.Email)
				oauthToken, transport = getOauthToken(request)
				glukitUser.RefreshToken = oauthToken.RefreshToken
			} else {
				transport.Token.RefreshToken = glukitUser.RefreshToken

				err := transport.Refresh(context)
				util.Propagate(err)

				context.Debugf("Storing new refreshed token [%s] in datastore...", oauthToken)
				glukitUser.LastUpdated = time.Now()
				glukitUser.Token = oauthToken
			}
		}
	}

	// Refresh and store the profile
	if service, err := oauth2.New(transport.Client()); err != nil {
		util.Propagate(err)
	} else {
		getRequest := service.Userinfo.Get()
		if userInfo, err := getRequest.Do(); err != nil {
			util.Propagate(err)

		} else {
			context.Infof("User profile refreshed to %v", userInfo)

			glukitUser.PictureUrl = userInfo.Picture
			glukitUser.FirstName = userInfo.Given_name
			glukitUser.LastName = userInfo.Family_name
		}
	}

	_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
	if err != nil {
		util.Propagate(err)
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

// getOauthToken deals with getting an oauth token from a oauth authorization code
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

// getEnvSettings returns either the production of the dev environment settings.
// TODO: check what best practices are surrounding the maintenance of clientIds and secrets.
func getEnvSettings() (host, clientId, clientSecret string) {
	if appengine.IsDevAppServer() {
		host = "localhost:8080"
		clientId = "***REMOVED***"
		clientSecret = "***REMOVED***"

	} else {
		host = "www.mygluk.it"
		clientId = "***REMOVED***"
		clientSecret = "***REMOVED***"
	}

	return host, clientId, clientSecret
}

// buildPerfectBaseline generates an array of reads that represents the target/perfection
func buildPerfectBaseline(glucoseReads []model.GlucoseRead) (reads []model.GlucoseRead) {
	reads = make([]model.GlucoseRead, len(glucoseReads))
	for i := range glucoseReads {
		reads[i] = model.GlucoseRead{glucoseReads[i].LocalTime, glucoseReads[i].Timestamp, model.TARGET_GLUCOSE_VALUE}
	}

	return reads
}
