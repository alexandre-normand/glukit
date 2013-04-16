package parser

import (
	"encoding/xml"
	"io"
	"os"
	"bufio"
	"container/list"
	"models"
	"utils"
	"fmt"
	"strings"
)

type Meter struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	Value        int    `xml:"Value,attr"`
}

type Event struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	EventTime    string `xml:"DisplayTime,attr"`
	EventType    string `xml:"EventType,attr"`
	Description  string `xml:"EventType,attr"`
}

// <Event InternalTime="2013-04-02 03:56:19" DisplayTime="2013-04-01 20:55:46" EventTime="2013-04-01 20:55:00"
//               EventType="Carbs" Decription="Carbs 8 grams"/>

func Parse(filepath string) (reads []models.MeterRead, carbIntakes []models.CarbIntake, exercises []models.Exercise, injections []models.Injection) {
	// open input file
	fi, err := os.Open(filepath)
	if err != nil { panic(err) }
	// close fi on exit and check for its returned error
	defer func() {
		if fi.Close() != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReader(fi)

	return ParseContent(r)
}

func ParseContent(reader io.Reader) (reads []models.MeterRead, carbIntakes []models.CarbIntake, exercises []models.Exercise, injections []models.Injection) {
	decoder := xml.NewDecoder(reader)
	injections = make([]models.Injection,0, 100)
	carbIntakes = make([]models.CarbIntake,0, 100)
	exercises = make([]models.Exercise,0, 100)

	readsList := list.New()
	readsList.Init()

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			// ...and its name is "Glucose"
			switch se.Name.Local {
			case "Glucose":
				var read Meter
				// decode a whole chunk of following XML into the
				decoder.DecodeElement(&read, &se)
				if (read.Value > 0) {
					readsList.PushBack(read)
				}
			case "Event":
				var event Event
				decoder.DecodeElement(&event, &se)
				if (event.EventType == "Carbs") {
					var carbQuantityInGrams int
					fmt.Sscanf(event.Description, "Carbs %d grams", &carbQuantityInGrams)
					carbIntake := models.CarbIntake{event.DisplayTime, utils.GetTimeInSeconds(event.DisplayTime, utils.TIMEZONE), carbQuantityInGrams}
					carbIntakes = append(carbIntakes, carbIntake)
				} else if (event.EventType == "Insulin") {
					var insulinUnits int
					fmt.Sscanf(event.Description, "Insulin %d units", &insulinUnits)
					injection := models.Injection{event.DisplayTime, utils.GetTimeInSeconds(event.DisplayTime, utils.TIMEZONE), insulinUnits}
					injections = append(injections, injection)
				} else if (strings.HasPrefix(event.EventType, "Exercise")) {
					var duration int
					var intensity string
					fmt.Sscanf(event.Description, "Exercise %s (%d minutes)", &intensity, &duration)
					exercise := models.Exercise{event.DisplayTime, utils.GetTimeInSeconds(event.DisplayTime, utils.TIMEZONE), duration, intensity}
					exercises = append(exercises, exercise)
				}

			case "Meter":
				// TODO: Read the meter calibrations?
			}
		}
	}

	return ConvertAsReadsArray(readsList), carbIntakes, exercises, injections
}

func ConvertAsReadsArray(meterReads *list.List) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, meterReads.Len())
	for e, i := meterReads.Front(), 0; e != nil; e, i = e.Next(), i + 1 {
		meter := e.Value.(Meter)
		reads[i] = models.MeterRead{meter.DisplayTime, utils.GetTimeInSeconds(meter.DisplayTime, utils.TIMEZONE), meter.Value}
	}

	return reads
}
