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
	"models"
	"drive"
	"store"
	"fetcher"
	"sort"
	"datautils"
	"bufio"
	"os"
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
//)

// Local
const (
	// Created at http://code.google.com/apis/console, these identify
	// our app for the OAuth protocol.
	CLIENT_ID     = "414109645872-g5og4q7pmua0na6sod0jtnvt16mdl4fh.apps.googleusercontent.com"
	CLIENT_SECRET = "U3KV6G8sYqxa-qtjoxRnk6tX"
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

var processFile = delay.Func("processFile", processData)

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

	context.Debugf("Called back")

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
	sysutils.Propagate(err)

	readData, key, _, _, _, err := store.GetUserData(context, user.Email)
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
		processData(t, files, context, key)

		context.Infof("Storing user data with key: %s", key.String())

		http.Redirect(w, request, "/graph", 303)
	}

}

func processData(t http.RoundTripper, files []*drive.File, context appengine.Context, userProfileKey *datastore.Key) {
	for i := range files {
		// TODO: Make this stream the content
		reader, err := fetcher.GetFileReader(t, files[i])
		if err != nil {
			context.Infof("Error reading file %s, skipping...", files[i].OriginalFilename)
		} else {
			lastReadTime := parser.ParseContent(context, reader, 500, userProfileKey, store.StoreDaysOfReads, store.StoreCarbs, store.StoreInjections, store.StoreExerciseData)
			store.LogFileImport(context, userProfileKey, models.FileImportLog{Id: files[i].Id, Md5Checksum: files[i].Md5Checksum, LastDataProcessed: lastReadTime})
			reader.Close()
		}
	}
}

func updateData(w http.ResponseWriter, r *http.Request) {
	url := config(r.Host).AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound)
}

func renderDemo(w http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	_, key, _, _, _, err := store.GetUserData(context, "demo@glukit.com")
	if err == datastore.ErrNoSuchEntity {
		context.Infof("No data found for demo user [%s], creating it", "demo@glukit.com")
		// TODO: Populate GlukitUser correctly, this will likely require getting rid of all data from the store when this is ready
		key, err = store.StoreUserProfile(context, time.Now(), models.GlukitUser{"demo@glukit.com", "", "", time.Now(), "", "", time.Now(), time.Unix(0, 0)})
		if err != nil {
			sysutils.Propagate(err)
		}

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

		parser.ParseContent(context, reader, 500, key, store.StoreDaysOfReads, store.StoreCarbs, store.StoreInjections, store.StoreExerciseData)
	} else {
		sysutils.Propagate(err)
	}

	render("demo@glukit.com", "/json.demo", w, request)
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

	_, _, meterReads, injections, carbIntakes, err := store.GetUserData(context, "demo@glukit.com")
	if err != nil {
		sysutils.Propagate(err)
	}

	context.Infof("Got %d reads", len(meterReads))

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

	_, _, reads, injections, carbIntakes, err := store.GetUserData(context, user.Email)
	if err != nil {
		sysutils.Propagate(err)
	}

	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	individuals := make([]DataSeries, 4)

	individuals[0] = DataSeries{"You", models.MeterReadSlice(reads).ToDataPointSlice(), "MeterReads"}
	individuals[1] = DataSeries{"You.Injection", models.InjectionSlice(injections).ToDataPointSlice(reads), "Injections"}
	individuals[2] = DataSeries{"You.Carbohydrates", models.CarbIntakeSlice(carbIntakes).ToDataPointSlice(reads), "CarbIntakes"}
	individuals[3] = DataSeries{"Perfection", models.MeterReadSlice(buildPerfectBaseline(reads)).ToDataPointSlice(), "ComparisonReads"}

	enc.Encode(individuals)
}

func generateTrackingData(writer http.ResponseWriter, request *http.Request, reads []models.MeterRead) {
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

	_, _, reads, _, _, err := store.GetUserData(context, user.Email)
	if err != nil {
		sysutils.Propagate(err)
	}

	generateTrackingData(writer, request, reads)
}

func demoTracking(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)

	value := writer.Header()
	value.Add("Content-type", "application/json")

	_, _, reads, _, _, err := store.GetUserData(context, "demo@glukit.com")
	if (err != nil) {
		sysutils.Propagate(err)
	}

	generateTrackingData(writer, request, reads)
}


func buildPerfectBaseline(meterReads []models.MeterRead) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, len(meterReads))
	for i := range meterReads {
		reads[i] = models.MeterRead{meterReads[i].LocalTime, meterReads[i].TimeValue, 83}
	}

	return reads
}
