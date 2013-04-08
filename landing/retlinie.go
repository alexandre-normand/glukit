package landing

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
	"models"
	"drive"
	"strings"
	"io/ioutil"
	"utils"
	"store"
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
	Reads     []models.MeterRead `json:"data"`
}

type RenderVariables struct {
	DataPath string
}

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
	utils.Propagate(err)

	file, err := FetchDataFileLocation(t.Client())
	content, err := DownloadFile(t, file)
	context := appengine.NewContext(request)
	// TODO Find a better way to get the user's email, not working on GAE
	if user, err := user.CurrentOAuth(context, ""); err == nil {
		if key, err := store.StoreUserData(file, user, w, context, content); err == nil {
			log.Printf("Stored user data with key: %s", key.String())
			http.Redirect(w, request, "/graph", 303)
		} else {
			utils.Propagate(err)
		}
	} else {
		utils.Propagate(err)
	}

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
	// TODO Find a better way to get the user's email, not working on GAE
	user, err := user.CurrentOAuth(context, "")
	if err != nil {
		utils.Propagate(err)
	}

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
	lastTime, _ := utils.ParseTime(lastValue.LocalTime)
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
		readTime, _ := utils.ParseTime(meter.LocalTime)
		if endOfDayIndex < 0 && readTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && readTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return meterReads[startOfDayIndex:endOfDayIndex + 1]
}
