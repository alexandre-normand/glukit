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
        <style>
        #chart_container {
                position: relative;
                font-family: Arial, Helvetica, sans-serif;
        }
        #chart {
                position: relative;
                left: 40px;
        }
        #y_axis {
                position: absolute;
                top: 0;
                bottom: 0;
                width: 40px;
        }
        </style>
  </head>
  <body>
  <div id="chart_container">
          <div id="y_axis"></div>
          <div id="chart"></div>
  </div>
<script>
var graph = new Rickshaw.Graph.Ajax( {
	element: document.getElementById("chart"),
	width: 800,
	height: 600,
	renderer: 'line',
	dataURL: '/json',
	series: [ {
			name: 'You',
            color: 'steelblue',
            width: 2,
		},  {
			name: 'Perfection',
		    color: '#33AD33',
		    width: 2,
		},
		{
			name: 'Scale',
			color: '#ffffff',
			width: 1,
				},
		 ],
	onComplete: function(transport) {
	    var x_axis = new Rickshaw.Graph.Axis.Time({
	          graph: transport.graph
	        });
	    x_axis.graph.update();

	    var y_axis = new Rickshaw.Graph.Axis.Y( {
		            graph: transport.graph,
		            orientation: 'left',
		            tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
		            element: document.getElementById('y_axis'),
		} );

		y_axis.graph.update();
	  }
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
	reads = make([]MeterRead, 2)
	reads[0] = MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 0}
	reads[1] = MeterRead{meterReads[0].LocalTime, meterReads[0].TimeValue, 300}

	return reads
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
	lastTime, _ := time.Parse(TIMEFORMAT, lastValue.DisplayTime)
	lowerBound := lastTime.Add(time.Duration(-24 * time.Hour))
	for e := meterReads.Front(); e != nil; e = e.Next() {
		meter := e.Value.(parser.Meter)
		readTime, _ := time.Parse(TIMEFORMAT, meter.DisplayTime)
		if !readTime.Before(lowerBound) && meter.Value > 0 {
			lastDay.PushBack(meter)
		}
    }

	return lastDay
}

func getTimeInSeconds(timeValue string) (value int64) {
	if timeValue, err := time.Parse(TIMEFORMAT, timeValue); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}
