// +build !appengine
package main

import (
	"code.google.com/p/gorilla/mux"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/config"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/goauth2/oauth"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
	"sync"
	"time"
)

var dataBrowserTemplate = template.Must(template.ParseFiles("view/templates/databrowser.html"))
var reportTemplate = template.Must(template.ParseFiles("view/templates/report.html"))
var landingTemplate = template.Must(template.ParseFiles("view/templates/landing.html"))
var muxRouter = mux.NewRouter()
var initOnce sync.Once

const (
	DEMO_PATH_PREFIX        = "demo."
	DEMO_PICTURE_URL        = "https://farm8.staticflickr.com/7389/10813078553_ab4e1397f4_b_d.jpg"
	GLUCOSE_UNIT_PARAMETER  = "unit"
	CONTENT_SECURITY_POLICY = "default-src 'self'; img-src *; style-src 'unsafe-inline' 'self'; font-src 'self' https://fonts.gstatic.com; connect-src 'self' https://api.stripe.com; frame-src 'self' https://js.stripe.com; script-src 'self' 'unsafe-inline' js.stripe.com *.googleapis.com"
)

// Some variables that are used during rendering of templates
type RenderVariables struct {
	PathPrefix           string
	StripePublishableKey string
	SSLHost              string
	GlucoseUnit          apimodel.GlucoseUnit
}

// init initializes the routes and global initialization
func main() {
	appConfig = config.NewAppConfig()

	http.Handle("/", muxRouter)

	// Create user Glukit Bernstein as a fallback for comparisons
	muxRouter.HandleFunc("/_ah/warmup", warmUp)
	muxRouter.HandleFunc("/initpower", warmUp)

	// GAE Json endpoints
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"data", demoContent)
	muxRouter.HandleFunc("/data", personalData)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"steadySailor", demoSteadySailorData)
	muxRouter.HandleFunc("/steadySailor", steadySailorData)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"dashboard", demoDashboard)
	muxRouter.HandleFunc("/dashboard", dashboard)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"glukitScores", glukitScoresForDemo)
	muxRouter.HandleFunc("/glukitScores", glukitScores)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"a1cs", a1cEstimatesForDemo)
	muxRouter.HandleFunc("/a1cs", a1cEstimates)
	muxRouter.HandleFunc("/donation", handleDonation)

	// "main"-page for both demo and real users
	muxRouter.HandleFunc("/demo", renderDemo)
	muxRouter.HandleFunc("/browse", renderRealUser)
	muxRouter.HandleFunc("/"+DEMO_PATH_PREFIX+"report", demoReport)
	muxRouter.HandleFunc("/report", report)

	// Static pages
	muxRouter.HandleFunc("/", landing)

	muxRouter.HandleFunc("/googleauth", handleRealUser)
	muxRouter.HandleFunc("/oauth2callback", oauthCallback)

	// Client API endpoints
	muxRouter.HandleFunc("/v1/calibrations", initializeAndHandleRequest).Methods("POST").Name(CALIBRATIONS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/injections", initializeAndHandleRequest).Methods("POST").Name(INJECTIONS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/meals", initializeAndHandleRequest).Methods("POST").Name(MEALS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/glucosereads", initializeAndHandleRequest).Methods("POST").Name(GLUCOSEREADS_V1_ROUTE)
	muxRouter.HandleFunc("/v1/exercises", initializeAndHandleRequest).Methods("POST").Name(EXERCISES_V1_ROUTE)

	// Register oauth endpoints to warmup which will initilize the oauth server and replace the routes with the actual oauth handlers
	muxRouter.HandleFunc("/token", initializeAndHandleRequest).Methods("POST").Name(TOKEN_ROUTE)
	muxRouter.HandleFunc("/authorize", initializeAndHandleRequest).Methods("GET").Name(AUTHORIZE_ROUTE)

	// Initialize task functions that would otherwise be prone to initialization loops
	refreshUserData = delay.Func(REFRESH_USER_DATA_FUNCTION_NAME, updateUserData)
	processFile = delay.Func(PROCESS_FILE_FUNCTION_NAME, processSingleFile)
	engine.RunGlukitScoreCalculationChunk = delay.Func(engine.GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, engine.RunGlukitScoreBatchCalculation)
	engine.RunA1CCalculationChunk = delay.Func(engine.A1C_BATCH_CALCULATION_FUNCTION_NAME, engine.RunA1CBatchCalculation)

	appengine.Main()
}

// landing executes the landing page template
func landing(w http.ResponseWriter, request *http.Request) {
	if err := landingTemplate.Execute(w, nil); err != nil {
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
		log.Infof(context, "No data found for demo user [%s], creating it", DEMO_EMAIL)
		dummyToken := oauth.Token{"", "", util.GLUKIT_EPOCH_TIME}
		// TODO: Populate GlukitUser correctly, this will likely require
		// getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(),
			model.GlukitUser{DEMO_EMAIL, "Demo", "OfMe", time.Now(), model.DIABETES_TYPE_1, "", time.Now(),
				apimodel.UNDEFINED_GLUCOSE_READ, dummyToken, "", model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, true, DEMO_PICTURE_URL, time.Now(),
				model.UNDEFINED_A1C_ESTIMATE})
		if err != nil {
			util.Propagate(err)
		}

		task, err := processDemoFile.Task(key)
		if err != nil {
			util.Propagate(err)
		}
		taskqueue.Add(context, task, DATASTORE_WRITES_QUEUE_NAME)

	} else if err != nil {
		util.Propagate(err)
	} else {
		log.Infof(context, "Data already stored for demo user [%s], continuing...", DEMO_EMAIL)
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
	user := user.Current(context)
	unitValue, err := resolveGlucoseUnit(user.Email, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderVariables := &RenderVariables{PathPrefix: DEMO_PATH_PREFIX, StripePublishableKey: appConfig.StripePublishableKey, SSLHost: appConfig.SSLHost, GlucoseUnit: *unitValue}

	if err := reportTemplate.Execute(w, renderVariables); err != nil {
		log.Criticalf(context, "Error executing template [%s]", dataBrowserTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// report executes the report page template
func report(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)
	unitValue, err := resolveGlucoseUnit(user.Email, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderVariables := &RenderVariables{PathPrefix: "", StripePublishableKey: appConfig.StripePublishableKey, SSLHost: appConfig.SSLHost, GlucoseUnit: *unitValue}
	w.Header().Set("Content-Security-Policy", CONTENT_SECURITY_POLICY)

	if err := reportTemplate.Execute(w, renderVariables); err != nil {
		log.Criticalf(context, "Error executing template [%s]", dataBrowserTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// render executed the graph page template
func render(email string, datapath string, w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	unitValue, err := resolveGlucoseUnit(email, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderVariables := &RenderVariables{PathPrefix: datapath, StripePublishableKey: appConfig.StripePublishableKey, SSLHost: appConfig.SSLHost, GlucoseUnit: *unitValue}

	w.Header().Set("Content-Security-Policy", CONTENT_SECURITY_POLICY)
	if err := dataBrowserTemplate.Execute(w, renderVariables); err != nil {
		log.Criticalf(context, "Error executing template [%s]", dataBrowserTemplate.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func resolveGlucoseUnit(email string, request *http.Request) (unit *apimodel.GlucoseUnit, err error) {
	rawUnitValue := request.FormValue(GLUCOSE_UNIT_PARAMETER)
	if rawUnitValue != apimodel.MMOL_PER_L && rawUnitValue != apimodel.MG_PER_DL {
		context := appengine.NewContext(request)
		glukitUser, _, _, err := store.GetUserData(context, email)
		if err != nil {
			return nil, err
		}

		unitValue := glukitUser.MostRecentRead.Unit
		return &unitValue, nil
	} else {
		unitValue := apimodel.GlucoseUnit(rawUnitValue)
		return &unitValue, nil
	}
}

// handleRealUser handles the flow for a real non-demo user. It will redirect to authorization if required
func handleRealUser(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	glukitUser, _, _, err := store.GetUserData(context, user.Email)
	if _, ok := err.(store.StoreError); err != nil && !ok || len(glukitUser.RefreshToken) == 0 {
		log.Infof(context, "Redirecting [%s], glukitUser [%v] for authorization. Error: [%v]", user.Email, glukitUser, err)

		configuration := configuration()
		log.Debugf(context, "We don't current have a refresh token (either lost or it's "+
			"the first access). Let's set the ApprovalPrompt to force to get a new one...")

		configuration.ApprovalPrompt = "force"

		url := configuration.AuthCodeURL(request.URL.RawQuery)
		http.Redirect(writer, request, url, http.StatusFound)
	} else {
		log.Infof(context, "User [%s] already exists with a valid refresh token [%s], skipping authorization step...",
			glukitUser.RefreshToken, user.Email)
		oauthCallback(writer, request)
	}
}

func warmUp(writer http.ResponseWriter, request *http.Request) {
	initOnce.Do(func() {
		c := appengine.NewContext(request)
		log.Infof(c, "Initializing application...")
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
