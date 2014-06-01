package glukit

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
	"appengine/user"
	"code.google.com/p/gorilla/mux"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/goauth2/oauth"
	"html/template"
	"net/http"
	"sync"
	"time"
)

var graphTemplate = template.Must(template.ParseFiles("view/templates/graph.html"))
var reportTemplate = template.Must(template.ParseFiles("view/templates/report.html"))
var landingTemplate = template.Must(template.ParseFiles("view/templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("view/templates/nodata.html"))
var muxRouter = mux.NewRouter()
var initOnce sync.Once

const (
	DEMO_PATH_PREFIX = "demo."
	DEMO_PICTURE_URL = "http://farm8.staticflickr.com/7389/10813078553_ab4e1397f4_b_d.jpg"
)

// Some variables that are used during rendering of templates
type RenderVariables struct {
	PathPrefix   string
	ChannelToken string
}

// init initializes the routes and global initialization
func init() {
	http.Handle("/", muxRouter)

	// Create user Glukit Bernstein as a fallback for comparisons
	muxRouter.HandleFunc("/_ah/warmup", warmUp)
	muxRouter.HandleFunc("/initpower", warmUp)

	// Json endpoints
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"data", demoContent)
	muxRouter.HandleFunc("/data", personalData)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"steadySailor", demoSteadySailorData)
	muxRouter.HandleFunc("/steadySailor", steadySailorData)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"dashboard", demoDashboard)
	muxRouter.HandleFunc("/dashboard", dashboard)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"glukitScores", glukitScoresForDemo)
	muxRouter.HandleFunc("/glukitScores", glukitScores)

	// "main"-page for both demo and real users
	muxRouter.HandleFunc("/demo", renderDemo)
	muxRouter.HandleFunc("/graph", renderRealUser)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"report", demoReport)
	muxRouter.HandleFunc("/report", report)

	// Static pages
	muxRouter.HandleFunc("/", landing)
	muxRouter.HandleFunc("/nodata", nodata)

	muxRouter.HandleFunc("/realuser", handleRealUser)
	muxRouter.HandleFunc("/oauth2callback", oauthCallback)

	// Client API endpoints
	muxRouter.HandleFunc("/v1/calibrations", initializeAndHandleRequest).Methods("POST").Name(CALIBRATIONS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/injections", initializeAndHandleRequest).Methods("POST").Name(INJECTIONS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/meals", initializeAndHandleRequest).Methods("POST").Name(MEALS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/glucosereads", initializeAndHandleRequest).Methods("POST").Name(GLUCOSEREADS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/exercises", initializeAndHandleRequest).Methods("POST").Name(EXERCISES_V1_ROUTE)

	// Register oauth endpoints to warmup which will initilize the oauth server and replace the routes with the actual oauth handlers
	muxRouter.HandleFunc("/token", initializeAndHandleRequest).Name(TOKEN_ROUTE)
	muxRouter.HandleFunc("/authorize", initializeAndHandleRequest).Name(AUTHORIZE_ROUTE)

	// Initialize task functions that would otherwise be prone to initialization loops
	refreshUserData = delay.Func(REFRESH_USER_DATA_FUNCTION_NAME, updateUserData)
	processFile = delay.Func(PROCESS_FILE_FUNCTION_NAME, processSingleFile)
	engine.RunGlukitScoreCalculationChunk = delay.Func(engine.GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, engine.RunGlukitScoreBatchCalculation)
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

	_, key, _, err := store.GetUserData(context, DEMO_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for demo user [%s], creating it", DEMO_EMAIL)
		dummyToken := oauth.Token{"", "", util.GLUKIT_EPOCH_TIME}
		// TODO: Populate GlukitUser correctly, this will likely require
		// getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(),
			model.GlukitUser{DEMO_EMAIL, "Demo", "OfMe", time.Now(), model.DIABETES_TYPE_1, "", time.Now(),
				model.UNDEFINED_GLUCOSE_READ, dummyToken, "", model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, true, DEMO_PICTURE_URL, time.Now()})
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

	render(DEMO_EMAIL, DEMO_PATH_PREFIX, w, request)
}

// renderRealUser executes the graph page template for a real user
func renderRealUser(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)
	render(user.Email, "", w, request)
}

// report executes the report page template
func demoReport(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	renderVariables := &RenderVariables{PathPrefix: DEMO_PATH_PREFIX, ChannelToken: "none"}

	if err := reportTemplate.Execute(w, renderVariables); err != nil {
		context.Criticalf("Error executing template [%s]", graphTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// report executes the report page template
func report(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	renderVariables := &RenderVariables{PathPrefix: "", ChannelToken: "none"}

	if err := reportTemplate.Execute(w, renderVariables); err != nil {
		context.Criticalf("Error executing template [%s]", graphTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	renderVariables := &RenderVariables{PathPrefix: datapath, ChannelToken: token}

	if err := graphTemplate.Execute(w, renderVariables); err != nil {
		context.Criticalf("Error executing template [%s]", graphTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleRealUser handles the flow for a real non-demo user. It will redirect to authorization if required
func handleRealUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	glukitUser, _, _, err := store.GetUserData(context, user.Email)
	if _, ok := err.(store.StoreError); err != nil && !ok || len(glukitUser.RefreshToken) == 0 {
		context.Infof("Redirecting [%s], glukitUser [%v] for authorization. Error: [%v]", user.Email, glukitUser, err)

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

func warmUp(writer http.ResponseWriter, request *http.Request) {
	initOnce.Do(func() {
		c := appengine.NewContext(request)
		c.Infof("Initializing application...")
		initializeApp(writer, request)
	})
}

func initializeAndHandleRequest(writer http.ResponseWriter, request *http.Request) {
	warmUp(writer, request)

	muxRouter.ServeHTTP(writer, request)
}

func initializeApp(writer http.ResponseWriter, request *http.Request) {
	initOauthProvider(writer, request)
	initApiEndpoints(writer, request)
	initializeGlukitBernstein(writer, request)
}
