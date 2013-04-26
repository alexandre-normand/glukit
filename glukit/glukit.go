package glukit

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
	stat "github.com/grd/stat"
)

// Appengine
//const (
//	// Created at http://code.google.com/apis/console, these identify
//	// our app for the OAuth protocol.
//	CLIENT_ID     = "414109645872-d6igmhnu0loafu53uphf8j67ou8ngjiu.apps.googleusercontent.com"
//	CLIENT_SECRET = "IYbOW0Aha34xMqTaPVO-_ar5"
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

type DataSeries struct {
	Name        string              `json:"name"`
	Data        []models.DataPoint  `json:"data"`
	Type        string              `json:"type"`
}

type RenderVariables struct {
	DataPath string
}

var graphTemplate = template.Must(template.ParseFiles("templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("templates/nodata.html"))

func init() {
	http.HandleFunc("/json.demo", demoContent)
	http.HandleFunc("/json", content)
	http.HandleFunc("/json.demo.tracking", demoTracking)
	http.HandleFunc("/json.tracking", tracking)

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
	individuals := make([]DataSeries, 4)
	individuals[0] = DataSeries{"You", models.MeterReadSlice(meterReads).ToDataPointSlice(), "MeterReads"}
	individuals[1] = DataSeries{"You.Injection", models.InjectionSlice(injections).ToDataPointSlice(meterReads), "Injections"}
	individuals[2] = DataSeries{"You.Carbohydrates", models.CarbIntakeSlice(carbIntakes).ToDataPointSlice(meterReads), "CarbIntakes"}
	individuals[3] = DataSeries{"Perfection", models.MeterReadSlice(buildPerfectBaseline(meterReads)).ToDataPointSlice(), "ComparisonReads"}

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
	individuals := make([]DataSeries, 4)
	// TODO: filter events to align with the reads
	individuals[0] = DataSeries{"You", models.MeterReadSlice(reads).ToDataPointSlice(), "MeterReads"}
	individuals[1] = DataSeries{"You.Injection", models.InjectionSlice(injections).ToDataPointSlice(reads), "Injections"}
	individuals[2] = DataSeries{"You.Carbohydrates", models.CarbIntakeSlice(carbIntakes).ToDataPointSlice(reads), "CarbIntakes"}
	individuals[3] = DataSeries{"Perfection", models.MeterReadSlice(buildPerfectBaseline(reads)).ToDataPointSlice(), "ComparisonReads"}

	enc.Encode(individuals)
}

func generateTrackingData(writer http.ResponseWriter, request *http.Request, reads []models.MeterRead) {
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	var trackingData models.TrackingData
	sort.Sort(models.ReadStatsSlice(reads))
	trackingData.Mean = stat.Mean(models.ReadStatsSlice(reads))
	trackingData.Deviation = stat.AbsdevMean(models.ReadStatsSlice(reads), 83)
	trackingData.Max, _ = stat.Max(models.ReadStatsSlice(reads))
	trackingData.Min, _ = stat.Min(models.ReadStatsSlice(reads))
	trackingData.Median = stat.MedianFromSortedData(models.ReadStatsSlice(reads))
	distribution := datautils.BuildHistogram(reads)
	sort.Sort(models.CoordinateSlice(distribution))
	trackingData.Distribution = distribution

	enc := json.NewEncoder(writer)
	enc.Encode(trackingData)
}

func tracking(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, reads, injections, carbIntakes, err := store.GetUserData(context, user)
	if err != nil {
		utils.Propagate(err)
	}

	reads, _, _ = datautils.GetLastDayOfData(reads, injections, carbIntakes)
	generateTrackingData(writer, request, reads)
}

func demoTracking(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	meterReads, carbIntakes, _, injections := parser.Parse("data.xml")
    meterReads, _, _ = datautils.GetLastDayOfData(meterReads, injections, carbIntakes)

	generateTrackingData(writer, request, meterReads)
}


func buildPerfectBaseline(meterReads []models.MeterRead) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, len(meterReads))
	for i := range meterReads {
		reads[i] = models.MeterRead{meterReads[i].LocalTime, meterReads[i].TimeValue, 83}
	}

	return reads
}
