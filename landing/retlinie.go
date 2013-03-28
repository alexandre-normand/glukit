package landing

import (
	"parser"
	"net/http"
	"fmt"
	"io"
	"time"
	"log"
	"encoding/json"
	"html/template"
)

type Response struct {
	Content   []Individual
}

type Individual struct {
	Name      string      `json:"name"`
	Reads     []MeterRead `json:"data"`
}
type MeterRead struct {
	LocalTime string   `json:"label"`
	TimeValue int      `json:"x"`
	Value     int      `json:"y"`
}

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)


var pageTemplate = template.Must(template.New("book").Parse(page))

const page = `
<html>
  <head>


    <link type="text/css" rel="stylesheet" href="http://code.shutterstock.com/rickshaw/src/css/graph.css">
    	<link type="text/css" rel="stylesheet" href="http://code.shutterstock.com/rickshaw/css/lines.css">

    	<script src="http://code.shutterstock.com/rickshaw/vendor/d3.v2.js"></script>

    	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.8.1/jquery.min.js"></script>

    	<script src="http://code.shutterstock.com/rickshaw/rickshaw.js"></script>

  </head>
  <body>
  <div id="chart_container">
  	<div id="chart"></div>
  </div>
<script>
var ajaxGraph = new Rickshaw.Graph.Ajax( {
	element: document.getElementById("chart"),
	width: 400,
	height: 200,
	renderer: 'line',
	dataURL: '/json',
	onData: function(d) { d[0].data[0].y = 80; return d },
	series: [
		{
			name: 'You',
            color: '#30c020',

		}
	]
} );

</script>

  </body>
</html>
`

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

    parser.Parse("data.xml", writer, writeReadAsJson)
	fmt.Fprintf(writer, "{}] }]")
}

func writeReadAsJson(meter *parser.Meter) {
	enc := json.NewEncoder(writer)
	meterRead := MeterRead{meter.DisplayTime, getTimeAsFloat(meter.DisplayTime), meter.Value}
	enc.Encode(meterRead)
	fmt.Fprintf(writer, ",")
}

func getTimeAsFloat(timeValue string) (value int) {
	if timeValue, err := time.Parse(TIMEFORMAT, timeValue); err == nil {
		return timeValue.Hour() * 60 + timeValue.Minute()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}
