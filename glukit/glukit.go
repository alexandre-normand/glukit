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

// Appengine
//const (
//	// Created at http://code.google.com/apis/console, these identify
//	// our app for the OAuth protocol.
//	CLIENT_ID     = "414109645872-d6igmhnu0loafu53uphf8j67ou8ngjiu.apps.googleusercontent.com"
//	CLIENT_SECRET = "IYbOW0Aha34xMqTaPVO-_ar5"
//	DEMO_EMAIL    = "demo@glukit.com"
//
//)

// Local
const (
	// Created at http://code.google.com/apis/console, these identify
	// our app for the OAuth protocol.
	CLIENT_ID     = "***REMOVED***"
	CLIENT_SECRET = "***REMOVED***"
	DEMO_EMAIL    = "demo@glukit.com"
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
	DataPath     string
	ChannelToken string
}

var graphTemplate = template.Must(template.ParseFiles("templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("templates/landing.html"))
var nodataTemplate = template.Must(template.ParseFiles("templates/nodata.html"))

var processFile = delay.Func("processSingleFile", processSingleFile)
var processDemoFile = delay.Func("processDemoFile", processStaticDemoFile)
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
	token, err := t.Exchange(code)
	sysutils.Propagate(err)

	readData, key, _, _, err := store.GetUserData(context, user.Email)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for user [%s], creating it", user.Email)
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(), models.GlukitUser{user.Email, "", "", time.Now(), "", "", time.Now(), time.Unix(0, 0)})
		if err != nil {
			sysutils.Propagate(err)
		}
	} else {
		sysutils.Propagate(err)
	}

	context.Debugf("Found existing user: %s", user.Email)
	lastUpdate := time.Unix(0, 0)
	if readData != nil {
		lastUpdate = readData.LastUpdated
	}

	context.Debugf("Key %s, lastUpdate: %s", key, lastUpdate)
	files, err := fetcher.SearchDataFiles(t.Client(), lastUpdate)
	if err != nil {
		sysutils.Propagate(err)
	}

	switch {
	case len(files) == 0 && readData == nil:
		context.Infof("No files found and user [%s] has no previous data stored", user.Email)
		http.Redirect(w, request, "/nodata", 303)
	case len(files) == 0 && readData != nil:
		context.Infof("No new or updated data found for existing user [%s]", user.Email)
		http.Redirect(w, request, "/graph", 303)
	case len(files) > 0:
		context.Infof("Found new data files for user [%s], downloading and storing...", user.Email)
		processData(token, files, context, user.Email, key)

		context.Infof("Storing user data with key: %s", key.String())

		http.Redirect(w, request, "/graph", 303)
	}

}

func getHost() (host string) {
	if (appengine.IsDevAppServer()) {
		host = "localhost:8080"
	} else {
		host = "glukit.appspot.com"
	}

	return host
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
		Config: config(getHost()),
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
		startTime := time.Unix(0, 0)
		if lastFileImportLog, err := store.GetFileImportLog(context, userProfileKey, file.Id); err == nil {
			startTime = lastFileImportLog.LastDataProcessed
			context.Infof("Reloading data from file [%s]-[%s] starting at date [%s]...", file.Id, file.OriginalFilename, startTime.Format(timeutils.TIMEFORMAT))
		} else if err == datastore.ErrNoSuchEntity {
			context.Debugf("First import of file [%s]-[%s]...", file.Id, file.OriginalFilename)
		} else if err != nil {
			sysutils.Propagate(err)
		}

		lastReadTime := parser.ParseContent(context, reader, 500, userProfileKey, startTime, store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreExerciseData)
		store.LogFileImport(context, userProfileKey, models.FileImportLog{Id: file.Id, Md5Checksum: file.Md5Checksum, LastDataProcessed: lastReadTime})
		reader.Close()
	}
	channel.Send(context, userEmail, "Refresh")
}

func updateData(w http.ResponseWriter, r *http.Request) {
	// TODO: use https://developers.google.com/appengine/docs/go/reference#AccessToken instead?
	url := config(r.Host).AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound)
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

	lastReadTime := parser.ParseContent(context, reader, 500, userProfileKey, time.Unix(0, 0), store.StoreDaysOfReads, store.StoreDaysOfCarbs, store.StoreDaysOfInjections, store.StoreExerciseData)
	store.LogFileImport(context, userProfileKey, models.FileImportLog{Id: "demo", Md5Checksum: "dummychecksum", LastDataProcessed: lastReadTime})
	channel.Send(context, DEMO_EMAIL, "Refresh")
}

func renderDemo(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	_, key, _, _, err := store.GetUserData(context, DEMO_EMAIL)
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for demo user [%s], creating it", DEMO_EMAIL)
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(), models.GlukitUser{DEMO_EMAIL, "", "", time.Now(), "", "", time.Now(), time.Unix(0, 0)})
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
