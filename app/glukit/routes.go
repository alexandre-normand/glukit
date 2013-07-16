package glukit

import (
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
	"appengine/user"
	"html/template"
	"lib/goauth2/oauth"
	"net/http"
	"time"
)

var graphTemplate = template.Must(template.ParseFiles("view/templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("view/templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("view/templates/nodata.html"))

// init initializes the routes and global initialization
func init() {
	// Json endpoints
	http.HandleFunc("/json.demo", demoContent)
	http.HandleFunc("/json", content)
	http.HandleFunc("/json.demo.tracking", demoTracking)
	http.HandleFunc("/json.tracking", tracking)

	// "main"-page for both demo and real users
	http.HandleFunc("/demo", renderDemo)
	http.HandleFunc("/graph", renderRealUser)

	// Static pages
	http.HandleFunc("/", landing)
	http.HandleFunc("/nodata", nodata)

	http.HandleFunc("/realuser", handleRealUser)
	http.HandleFunc("/oauth2callback", oauthCallback)

	refreshUserData = delay.Func("refreshUserData", updateUserData)
}

// landing executes the landing page template
func landing(w http.ResponseWriter, request *http.Request) {
	if err := landingTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// nodata renders the no-data available page that shows up when a user
// first accesses the app and doesn't have any dexcom files on Google Drive
func nodata(w http.ResponseWriter, request *http.Request) {
	if err := nodataTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// oauthCallback is invoked on return of the google oauth flow after being
// the user comes back from being redirected to google oauth for authorization
func oauthCallback(writer http.ResponseWriter, request *http.Request) {
	handleLoggedInUser(writer, request)
}

// renderDemo executes the graph template for the demo user
func renderDemo(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	_, key, _, _, err := store.GetUserData(context, DEMO_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for demo user [%s], creating it", DEMO_EMAIL)
		dummyToken := oauth.Token{"", "", util.BEGINNING_OF_TIME}
		// TODO: Populate GlukitUser correctly, this will likely require
		// getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(),
			model.GlukitUser{DEMO_EMAIL, "", "", time.Now(), "", "", time.Now(),
				util.BEGINNING_OF_TIME, dummyToken, "", model.UNDEFINED_SCORE})
		if err != nil {
			util.Propagate(err)
		}

		task, err := processDemoFile.Task(key)
		if err != nil {
			util.Propagate(err)
		}
		taskqueue.Add(context, task, "store")
	} else if err != nil {
		util.Propagate(err)
	} else {
		context.Infof("Data already stored for demo user [%s], continuing...", DEMO_EMAIL)
	}

	render(DEMO_EMAIL, "/json.demo", w, request)
}

// renderRealUser executes the graph page template for a real user
func renderRealUser(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)
	render(user.Email, "/json", w, request)
}

// render executed the graph page template
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

// handleRealUser handles the flow for a real non-demo user. It will redirect to authorization if required
func handleRealUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	glukitUser, _, _, _, err := store.GetUserData(context, user.Email)
	if err != nil || len(glukitUser.RefreshToken) == 0 {
		context.Infof("Redirecting [%s], glukitUser [%v] for authorization", user.Email, glukitUser)

		configuration := config()
		context.Debugf("We don't current have a refresh token (either lost or it's " +
			"the first access). Let's set the ApprovalPrompt to force to get a new one...")
		configuration.ApprovalPrompt = "force"

		url := configuration.AuthCodeURL(request.URL.RawQuery)
		http.Redirect(writer, request, url, http.StatusFound)
	} else {
		context.Infof("User [%s] already exists with a valid refresh token [%s], skipping authorization step...",
			glukitUser.RefreshToken, user.Email)
		oauthCallback(writer, request)
	}
}
