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
)

type Individual struct {
	Name      string      `json:"name"`
	Reads     []MeterRead `json:"data"`
}
type MeterRead struct {
	LocalTime string   `json:"label"`
	TimeValue int64    `json:"x"`
	Value     int      `json:"y"`
}

const (
	TIMEFORMAT = "2006-01-02 15:04:05 MST"
	TIMEZONE = "PST"
)

var TIMEZONE_LOCATION, _ = time.LoadLocation("America/Los_Angeles")
var pageTemplate = template.Must(template.ParseFiles("templates/landing.html"))

func init() {
	http.HandleFunc("/json", content)
	http.HandleFunc("/", render)
}


func render(w http.ResponseWriter, r *http.Request) {
	if err := pageTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("New Request\n")

}

func content(writer http.ResponseWriter, request *http.Request) {
	//	c := appengine.NewContext(request)
	//	u := user.Current(c)
	//	if u == nil {
	//		url, err := user.LoginURL(c, request.URL.String())
	//		if err != nil {
	//			http.Error(writer, err.Error(), http.StatusInternalServerError)
	//			return
	//		}
	//		writer.Header().Set("Location", url)
	//		writer.WriteHeader(http.StatusFound)
	//		return
	//	}

	fmt.Printf("New Request\n")
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	reads := parser.Parse("data.xml", writer)
	meterReads := convertAsReadsArray(getLastDayOfData(reads))
	enc := json.NewEncoder(writer)
	individuals := make([]Individual, 3)
	individuals[0] = Individual{"You", meterReads}
	individuals[1] = Individual{"Perfection", buildPerfectBaseline(meterReads)}
	individuals[2] = Individual{"Scale", buildScaleValues(meterReads)}
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
		previousDay := lastTime.Add(time.Duration(-24 * time.Hour))
		upperBound = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	} else {
		upperBound = time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	}
	lowerBound := upperBound.Add(time.Duration(-24 * time.Hour))
	for e := meterReads.Front(); e != nil; e = e.Next() {
		meter := e.Value.(parser.Meter)
		readTime, _ := parseTime(meter.DisplayTime)
		log.Printf("ReadTime is %s, lower bound is %s, upperBound is %s", readTime.String(),
			lowerBound.String(), upperBound.String())
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

		log.Printf("Parsed time from %s is %s", timeValue + " " + TIMEZONE, value.String())

	}

	return value, err;
}
