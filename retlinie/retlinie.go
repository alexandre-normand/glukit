package retlinie

import (
	"parser"
	"net/http"
	"fmt"
	"time"
	"log"
	"encoding/json"
	"html/template"
	"goauth2/oauth"
	"appengine"
	"appengine/user"
	"appengine/urlfetch"
	"appengine/datastore"
	"models"
	"drive"
	"utils"
	"store"
	"fetcher"
	"strings"
	"sort"
	"datautils"
)

// Appengine
//const (
//	// Created at http://code.google.com/apis/console, these identify
//	// our app for the OAuth protocol.
//	CLIENT_ID     = "414109645872-adbrmoh7te4mgbvr9f7rnj26j66bverl.apps.googleusercontent.com"
//	CLIENT_SECRET = "IcbtRurZqPa2PV6NnSIgay73"
//)

// Local
const (
	// Created at http://code.google.com/apis/console, these identify
	// our app for the OAuth protocol.
	CLIENT_ID     = "***REMOVED***"
	CLIENT_SECRET = "***REMOVED***"
)

// config returns the configuration information for OAuth and Drive.
func config(host string) *oauth.Config {
	return &oauth.Config{
		ClientId:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile " + drive.DriveReadonlyScope,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  fmt.Sprintf("http://%s/oauth2callback", host),
	}
}

type IndividualReads struct {
	Name        string              `json:"name"`
	Reads       []models.MeterRead  `json:"data"`
}

type IndividualCarbIntakes struct {
	Name        string              `json:"name"`
	CarbIntakes []models.CarbIntake `json:"data"`
}

type IndividualInjections struct {
	Name        string              `json:"name"`
	Injections  []models.Injection  `json:"data"`
}

type RenderVariables struct {
	DataPath string
}

var graphTemplate = template.Must(template.ParseFiles("templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("templates/nodata.html"))

func init() {
	http.HandleFunc("/json.demo", demoContent)
	http.HandleFunc("/json.demo.injections", demoInjections)
	http.HandleFunc("/json.demo.carbs", demoCarbs)
	http.HandleFunc("/json", content)
	http.HandleFunc("/json.injections", injections)
	http.HandleFunc("/json.carbs", carbs)

	http.HandleFunc("/demo", renderDemo)
	http.HandleFunc("/graph", renderRealUser)
	http.HandleFunc("/", landing)
	http.HandleFunc("/nodata", nodata)
	http.HandleFunc("/realuser", updateData)
	http.HandleFunc("/oauth2callback", callback)
}

func callback(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	// Exchange code for an access token at OAuth provider.
	code := request.FormValue("code")
	t := &oauth.Transport{
		Config: config(request.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(request),
		},
	}

	// TODO: save the token to the memcache/datastore?!
	_, err := t.Exchange(code)
	utils.Propagate(err)

	readData, _, _, _, err := store.GetUserData(context, user)
	if err == datastore.ErrNoSuchEntity {
		log.Printf("No data found for user [%s]", user.Email)
	} else {
		utils.Propagate(err)
	}

	lastUpdate := time.Unix(0, 0)
	if readData != nil {
		lastUpdate = readData.LastUpdated
	}

	thisUpdate := time.Now()
	files, err := fetcher.SearchDataFiles(t.Client(), lastUpdate)
	if err != nil {
		utils.Propagate(err)
	}

	switch {
	case len(files) == 0 && readData == nil:
		log.Printf("No files found and user [%s] has no previous data stored", user.Email)
		http.Redirect(w, request, "/nodata", 303)
	case len(files) == 0 && readData != nil:
		log.Printf("No new or updated data found for existing user [%s]", user.Email)
		http.Redirect(w, request, "/graph", 303)
	case len(files) > 0:
		log.Printf("Found new data files for user [%s], downloading and storing...", user.Email)
		reads, injections, carbIntakes := getAllData(t, files)

		if key, err := store.StoreUserData(thisUpdate, user, w, context, reads, injections, carbIntakes); err == nil {
			log.Printf("Stored user data with key: %s", key.String())
			http.Redirect(w, request, "/graph", 303)
		} else {
			utils.Propagate(err)
		}
	}
}

func getAllData(t http.RoundTripper, files []*drive.File) (readData []models.MeterRead, injectionData []models.Injection, carbIntakeData []models.CarbIntake) {
	var reads []models.MeterRead
	var carbIntakes []models.CarbIntake
	var injections []models.Injection

	for i := range files {
		content, err := fetcher.DownloadFile(t, files[i])
		if err != nil {
			log.Printf("Error reading file %s, skipping...", files[i].OriginalFilename)
		} else {
			fileReads, fileCarbIntakes, _, fileInjections := parser.ParseContent(strings.NewReader(content))
			reads = utils.MergeReadArrays(reads, fileReads)
			carbIntakes = utils.MergeCarbIntakeArrays(carbIntakes, fileCarbIntakes)
			injections = utils.MergeInjectionArrays(injections, fileInjections)
		}
	}
	sort.Sort(models.MeterReadSlice(reads))
	sort.Sort(models.CarbIntakeSlice(carbIntakes))
	sort.Sort(models.InjectionSlice(injections))

	return reads, injections, carbIntakes
}

func updateData(w http.ResponseWriter, r *http.Request) {
	url := config(r.Host).AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound)
}

func renderDemo(w http.ResponseWriter, r *http.Request) {
	renderVariables := &RenderVariables{DataPath: "/json.demo"}
	render(w, r, renderVariables)
}

func renderRealUser(w http.ResponseWriter, r *http.Request) {
	renderVariables := &RenderVariables{DataPath: "/json"}
	render(w, r, renderVariables)
}

func render(w http.ResponseWriter, request *http.Request, renderVariables *RenderVariables) {
	if err := graphTemplate.Execute(w, renderVariables); err != nil {
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
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	meterReads, carbIntakes, _, injections := parser.Parse("data.xml")
	meterReads, injections, carbIntakes = datautils.GetLastDayOfData(meterReads, injections, carbIntakes)

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualReads, 3)
	individuals[0] = IndividualReads{"You", meterReads}
	individuals[1] = IndividualReads{"Perfection", buildPerfectBaseline(meterReads)}
	individuals[2] = IndividualReads{"Scale", buildScaleValues(meterReads)}
	enc.Encode(individuals)
}

func demoInjections(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	meterReads, carbIntakes, _, injections := parser.Parse("data.xml")
	_, injections, _ = datautils.GetLastDayOfData(meterReads, injections, carbIntakes)

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualInjections, 1)
	individuals[0] = IndividualInjections{"You", injections}
	enc.Encode(individuals)
}

func demoCarbs(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	meterReads, carbIntakes, _, injections := parser.Parse("data.xml")
	_, _, carbIntakes = datautils.GetLastDayOfData(meterReads, injections, carbIntakes)

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualCarbIntakes, 1)
	individuals[0] = IndividualCarbIntakes{"You", carbIntakes}
	enc.Encode(individuals)
}

func content(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, reads, injections, carbIntakes, err := store.GetUserData(context, user)
	if err != nil {
		utils.Propagate(err)
	}

	reads, injections, carbIntakes = datautils.GetLastDayOfData(reads, injections, carbIntakes)

	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualReads, 3)
	// TODO: filter events to align with the reads
	individuals[0] = IndividualReads{"You", reads}
	individuals[1] = IndividualReads{"Perfection", buildPerfectBaseline(reads)}
	individuals[2] = IndividualReads{"Scale", buildScaleValues(reads)}
	enc.Encode(individuals)
}

func injections(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, reads, injections, carbIntakes, err := store.GetUserData(context, user)
	if err != nil {
		utils.Propagate(err)
	}

	_, injections, _ = datautils.GetLastDayOfData(reads, injections, carbIntakes)

	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualInjections, 1)
	// TODO: filter events to align with the reads
	individuals[0] = IndividualInjections{"You", injections}
	enc.Encode(individuals)
}

func carbs(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, reads, injections, carbIntakes, err := store.GetUserData(context, user)
	if err != nil {
		utils.Propagate(err)
	}

	_, _, carbIntakes = datautils.GetLastDayOfData(reads, injections, carbIntakes)

	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	individuals := make([]IndividualCarbIntakes, 1)
	// TODO: filter events to align with the reads
	individuals[0] = IndividualCarbIntakes{"You", carbIntakes}
	enc.Encode(individuals)
}

func buildPerfectBaseline(meterReads []models.MeterRead) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, len(meterReads))
	for i := range meterReads {
		reads[i] = models.MeterRead{meterReads[i].LocalTime, meterReads[i].TimeValue, 83}
	}

	return reads
}

// Stupid hack until I figure out how to set the min/max on the Y-axis
func buildScaleValues(meterReads []models.MeterRead) (reads []models.MeterRead) {
	if len(meterReads) > 0 {
		reads = make([]models.MeterRead, 2)
		reads[0] = models.MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 0}
		reads[1] = models.MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 300}
		return reads
	}

	return []models.MeterRead {};
}
