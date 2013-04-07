package landing

import (
	"parser"
	"net/http"
	"fmt"
	"time"
	"log"
	"encoding/json"
	"html/template"
	"container/list"
	"goauth2/oauth"
	"appengine"
	"appengine/user"
	"appengine/urlfetch"
	"appengine/datastore"
	"appengine/blobstore"
	"drive"
	"strings"
	"io/ioutil"
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

type Individual struct {
	Name      string      `json:"name"`
	Reads     []MeterRead `json:"data"`
}
type MeterRead struct {
	LocalTime string   `json:"label"`
	TimeValue int64    `json:"x"`
	Value     int      `json:"y"`
}

type ReadData struct {
	Email         string
	Name          string
	LastUpdated   time.Time
	ReadsBlobKey  appengine.BlobKey
}

type RenderVariables struct {
	DataPath string
}

const (
	TIMEFORMAT       = "2006-01-02 15:04:05 MST"
	DRIVE_TIMEFORMAT = "2006-01-02T15:04:05.000Z"
	TIMEZONE         = "PST"
)

var TIMEZONE_LOCATION, _ = time.LoadLocation("America/Los_Angeles")
var graphTemplate = template.Must(template.ParseFiles("templates/graph.html"))
var landingTemplate = template.Must(template.ParseFiles("templates/landing.html"))

func init() {
	http.HandleFunc("/json.demo", demoContent)
	http.HandleFunc("/json", content)
	http.HandleFunc("/demo", renderDemo)
	http.HandleFunc("/graph", renderRealUser)
	http.HandleFunc("/", landing)
	http.HandleFunc("/realuser", updateData)
	http.HandleFunc("/oauth2callback", callback)
}

func callback(w http.ResponseWriter, request *http.Request) {
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
	check(err)

	file, err := FetchDataFileLocation(t.Client())
	content, err := DownloadFile(t, file)
	context := appengine.NewContext(request)
	if user, err := user.CurrentOAuth(context, ""); err == nil {
		if key, err := StoreUserData(file, user, w, context, content); err == nil {
			log.Printf("Stored user data with key: %s", key.String())
			http.Redirect(w, request, "/graph", 200)
		} else {
			check(err)
		}
	} else {
		check(err)
	}

}

func StoreUserData(file *drive.File, user *user.User, writer http.ResponseWriter, context appengine.Context, content string) (key *datastore.Key, err error) {
	reads := convertAsReadsArray(parser.ParseContent(strings.NewReader(content)))
	blobKey, err := StoreReadsBlob(context, reads)
	if err != nil {
		check(err)
	}

	modifiedDate, error := parseGoogleDriveDate(file.ModifiedDate)
	log.Printf("Modified date is: %s", modifiedDate.String())
	if error != nil {
		check(error)
	}
	readData := ReadData{Email: user.Email,
		Name: user.String(),
		LastUpdated: modifiedDate,
		ReadsBlobKey: blobKey}


	key, error = datastore.Put(context, datastore.NewKey(context, "ReadData", user.Email, 0, nil), &readData)
	if error != nil {
		check(error)
	}

	return key, nil
}

func StoreReadsBlob(context appengine.Context, reads []MeterRead) (blobKey appengine.BlobKey, err error) {
	var k appengine.BlobKey
	w, err := blobstore.Create(context, "application/json")
	if err != nil {
		return k, err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(reads)
	if err != nil {
		return k, err
	}
	err = w.Close()
	if err != nil {
		return k, err
	}

	return w.Key()
}

func FetchReadsBlob(context appengine.Context, blobKey appengine.BlobKey) (reads []MeterRead, err error) {
	reader := blobstore.NewReader(context, blobKey)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&reads); err != nil {
		return nil, err
	}

	return reads, nil
}

func parseGoogleDriveDate(value string) (timeValue time.Time, err error) {
	return time.Parse(DRIVE_TIMEFORMAT, value)
}

func FetchDataFileLocation(client *http.Client) (file *drive.File, err error) {
	if service, err := drive.New(client); err != nil {
		return nil, err
	} else {
		call := service.Files.List().MaxResults(10).Q("fullText contains \"<Glucose\" and fullText contains \"<Patient Id=\"")
		if filelist, err := call.Do(); err != nil {
			return nil, err
		} else {
			for i := range filelist.Items {
				file := filelist.Items[i]
				if strings.HasSuffix(file.OriginalFilename, ".Export.xml") {
					log.Printf("Found match: %s\n", file)
					return file, nil
				} else {
					log.Printf("Skipping search result item: %s\n", file)
				}
			}
		}
	}

	return nil, nil
}

func DownloadFile(t http.RoundTripper, f *drive.File) (string, error) {
	// t parameter should use an oauth.Transport
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		log.Printf("An error occurred: File is not downloadable")
		return "", nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	resp, err := t.RoundTrip(req)
	// Make sure we close the Body later
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	return string(body), nil
}

// check aborts the current execution if err is non-nil.
func check(err error) {
	if err != nil {
		panic(err)
	}
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
	c := appengine.NewContext(request)
	u := user.Current(c)
	if u != nil {
		http.Redirect(w, request, "/realuser", http.StatusFound)
	} else {
		if err := landingTemplate.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func demoContent(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	reads := parser.Parse("data.xml")
	meterReads := convertAsReadsArray(getLastDayOfData(reads))
	enc := json.NewEncoder(writer)
	individuals := make([]Individual, 3)
	individuals[0] = Individual{"You", meterReads}
	individuals[1] = Individual{"Perfection", buildPerfectBaseline(meterReads)}
	individuals[2] = Individual{"Scale", buildScaleValues(meterReads)}
	enc.Encode(individuals)
}

func content(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	readData := new(ReadData)
	user, err := user.CurrentOAuth(context, "")
    if err != nil {
		check(err)
	}

	error := datastore.Get(context, datastore.NewKey(context, "ReadData", user.Email, 0, nil), &readData)
	if error != nil {
		check(error)
	}

	reads, err := FetchReadsBlob(context, readData.ReadsBlobKey)
	if err != nil {
		check(err)
	}

	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	//meterReads := convertAsReadsArray(getLastDayOfData(reads))
	enc := json.NewEncoder(writer)
	individuals := make([]Individual, 3)
	individuals[0] = Individual{"You", reads}
	individuals[1] = Individual{"Perfection", buildPerfectBaseline(reads)}
	individuals[2] = Individual{"Scale", buildScaleValues(reads)}
	enc.Encode(individuals)
}

func buildPerfectBaseline(meterReads []MeterRead) (reads []MeterRead) {
	reads = make([]MeterRead, len(meterReads))
	for i := range meterReads {
		reads[i] = MeterRead{meterReads[i].LocalTime, meterReads[i].TimeValue, 83}
	}

	return reads
}

// Stupid hack until I figure out how to set the min/max on the Y-axis
func buildScaleValues(meterReads []MeterRead) (reads []MeterRead) {
	if len(meterReads) > 0 {
		reads = make([]MeterRead, 2)
		reads[0] = MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 0}
		reads[1] = MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 300}
		return reads
	}

	return []MeterRead {};
}

func convertAsReadsArray(meterReads *list.List) (reads []MeterRead) {
	reads = make([]MeterRead, meterReads.Len())
	for e, i := meterReads.Front(), 0; e != nil; e, i = e.Next(), i + 1 {
		meter := e.Value.(parser.Meter)
		reads[i] = MeterRead{meter.DisplayTime, getTimeInSeconds(meter.DisplayTime), meter.Value}
	}

	return reads
}

func getLastDayOfData(meterReads *list.List) (lastDay *list.List) {
	lastDay = list.New()
	lastDay.Init()
	lastValue := meterReads.Back().Value.(parser.Meter);
	lastTime, _ := parseTime(lastValue.DisplayTime)
	var upperBound time.Time;
	if (lastTime.Hour() < 6) {
		// Rewind by one more day
		previousDay := lastTime.Add(time.Duration(-24*time.Hour))
		upperBound = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	} else {
		upperBound = time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	}
	lowerBound := upperBound.Add(time.Duration(-24*time.Hour))
	for e := meterReads.Front(); e != nil; e = e.Next() {
		meter := e.Value.(parser.Meter)
		readTime, _ := parseTime(meter.DisplayTime)
		if readTime.Before(upperBound) && readTime.After(lowerBound) && meter.Value > 0 {
			lastDay.PushBack(meter)
		}
	}

	return lastDay
}

func getTimeInSeconds(timeValue string) (value int64) {
	if timeValue, err := parseTime(timeValue); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}

func parseTime(timeValue string) (value time.Time, err error) {
	if value, err = time.Parse(TIMEFORMAT, timeValue + " " + TIMEZONE); err == nil {
		value = time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(),
			value.Nanosecond(), TIMEZONE_LOCATION)
	}

	return value, err;
}
