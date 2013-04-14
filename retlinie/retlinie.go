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

type Individual struct {
	Name      string      `json:"name"`
	Reads     []models.MeterRead `json:"data"`
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

	readData, _, err := store.GetUserData(context, user)
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
		reads := getAllData(t, files)

		if key, err := store.StoreUserData(thisUpdate, user, w, context, reads); err == nil {
			log.Printf("Stored user data with key: %s", key.String())
			http.Redirect(w, request, "/graph", 303)
		} else {
			utils.Propagate(err)
		}
	}
}

func getAllData(t http.RoundTripper, files []*drive.File) (readData []models.MeterRead) {
	var reads []models.MeterRead
	for i := range files {
		content, err := fetcher.DownloadFile(t, files[i])
		if err != nil {
			log.Printf("Error reading file %s, skipping...", files[i].OriginalFilename)
		} else {
			fileReads := parser.ParseContent(strings.NewReader(content))
			reads = utils.MergeArrays(reads, fileReads)
		}
	}
	sort.Sort(models.MeterReadSlice(reads))

	return reads
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

	meterReads := parser.Parse("data.xml")
	meterReads = GetLastDayOfData(meterReads)

	enc := json.NewEncoder(writer)
	individuals := make([]Individual, 3)
	individuals[0] = Individual{"You", meterReads}
	individuals[1] = Individual{"Perfection", buildPerfectBaseline(meterReads)}
	individuals[2] = Individual{"Scale", buildScaleValues(meterReads)}
	enc.Encode(individuals)
}

func content(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, reads, err := store.GetUserData(context, user)
	if err != nil {
		utils.Propagate(err)
	}

	reads = GetLastDayOfData(reads)

	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	individuals := make([]Individual, 3)
	individuals[0] = Individual{"You", reads}
	individuals[1] = Individual{"Perfection", buildPerfectBaseline(reads)}
	individuals[2] = Individual{"Scale", buildScaleValues(reads)}
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

// Assumes reads are ordered by time
func GetLastDayOfData(meterReads []models.MeterRead) (lastDayOfReads []models.MeterRead) {
	dataSize := len(meterReads)
	startOfDayIndex := -1
	endOfDayIndex := -1

	lastValue := meterReads[dataSize - 1]
	lastTime, _ := utils.ParseTime(lastValue.LocalTime, utils.TIMEZONE)
	var upperBound time.Time;
	if (lastTime.Hour() < 6) {
		// Rewind by one more day
		previousDay := lastTime.Add(time.Duration(-24*time.Hour))
		upperBound = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, utils.TIMEZONE_LOCATION)
	} else {
		upperBound = time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 6, 0, 0, 0, utils.TIMEZONE_LOCATION)
	}
	lowerBound := upperBound.Add(time.Duration(-24*time.Hour))
	for i := dataSize - 1; i > 0; i-- {
		meter := meterReads[i]
		readTime, _ := utils.ParseTime(meter.LocalTime, utils.TIMEZONE)
		if endOfDayIndex < 0 && readTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && readTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return meterReads[startOfDayIndex:endOfDayIndex + 1]
}
