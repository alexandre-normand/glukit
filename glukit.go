// The glukit package is the main package for the application. This is where it all starts.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/config"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleuser "google.golang.org/api/oauth2/v2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

const (
	DEMO_EMAIL = "demo@glukit.com"
)

var emptyDataPointSlice []apimodel.DataPoint
var appConfig *config.AppConfig

// config returns the configuration information for OAuth.
func configuration() *oauth2.Config {
	configuration := oauth2.Config{
		ClientID:     appConfig.GoogleClientId,
		ClientSecret: appConfig.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("https://%s/oauth2callback", appConfig.Host),
	}

	return &configuration
}

// handleUserLogin starts the login flow with google
func handleUserLogin(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	conf := configuration()

	// Refresh and store the profile
	state, err := randomHex(20)
	if err != nil {
		util.Propagate(err)
	}

	url := conf.AuthCodeURL(state)
	log.Infof(context, "Visiting url [%s] for auth dialog...", url)

	http.Redirect(writer, request, url, http.StatusTemporaryRedirect)
}

// randomHex generates a random state string
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// handleLoggedInUser is responsible for directing the user to the graph page after optionally:
//   Storing the GlukitUser entry if it's the first access
func handleLoggedInUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	log.Infof(context, "Oauth callback handling")

	conf := configuration()
	code := request.FormValue("code")
	tok, err := conf.Exchange(context, code)
	if err != nil {
		util.Propagate(err)
	}

	client := conf.Client(context, tok)
	if service, err := googleuser.New(client); err != nil {
		util.Propagate(err)
	} else {
		getRequest := service.Userinfo.Get()

		if userInfo, err := getRequest.Do(); err != nil {
			util.Propagate(err)

		} else {
			log.Infof(context, "User profile refreshed to %v", userInfo)

			log.Infof(context, "Got user info for logged in google account [%s]", userInfo)
			glukitUser, _, _, err := store.GetUserData(context, userInfo.Email)

			if err == datastore.ErrNoSuchEntity {
				log.Infof(context, "No data found for user [%s], creating it", userInfo.Email)

				// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when
				// this is ready
				// We store the refresh token separately from the rest. This token is long-lived, meaning that if
				// we have a glukit user with no refresh token, we need to force getting a new one (which is to be avoided)
				glukitUser = &model.GlukitUser{userInfo.Email, userInfo.GivenName, userInfo.FamilyName, time.Now(),
					model.DIABETES_TYPE_1, "", util.GLUKIT_EPOCH_TIME, apimodel.UNDEFINED_GLUCOSE_READ,
					model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, userInfo.Picture, time.Now(), model.UNDEFINED_A1C_ESTIMATE}
				_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
				if err != nil {
					util.Propagate(err)
				}
			} else if _, ok := err.(store.StoreError); err != nil && !ok {
				util.Propagate(err)
			} else {
				glukitUser.PictureUrl = userInfo.Picture
				glukitUser.FirstName = userInfo.GivenName
				glukitUser.LastName = userInfo.FamilyName

				_, err = store.StoreUserProfile(context, time.Now(), *glukitUser)
				if err != nil {
					util.Propagate(err)
				}
			}
		}
	}

	/**
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}*/

	// Render the graph view, it might take some time to show something but it will as soon as a file import
	// completes
	renderRealUser(writer, request)
}

// buildPerfectBaseline generates an array of reads that represents the target/perfection
func buildPerfectBaseline(glucoseReads []apimodel.GlucoseRead) (reads []apimodel.GlucoseRead) {
	reads = make([]apimodel.GlucoseRead, len(glucoseReads))
	for i := range glucoseReads {
		reads[i] = apimodel.GlucoseRead{glucoseReads[i].Time, apimodel.MG_PER_DL, model.TARGET_GLUCOSE_VALUE}
	}

	return reads
}
