package parser

import (
	"encoding/xml"
	"io"
	"log"
	"container/list"
	"appengine/datastore"
	"appengine"
	"models"
	"sysutils"
	"timeutils"
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
	EventTime    string `xml:"EventTime,attr"`
	EventType    string `xml:"EventType,attr"`
	Description  string `xml:"Decription,attr"`
}

// <Event InternalTime="2013-04-02 03:56:19" DisplayTime="2013-04-01 20:55:46" EventTime="2013-04-01 20:55:00"
//               EventType="Carbs" Decription="Carbs 8 grams"/>

//func Parse(filepath string) (reads []models.MeterRead, carbIntakes []models.CarbIntake, exercises []models.Exercise, injections []models.Injection) {
//	// open input file
//	fi, err := os.Open(filepath)
//	if err != nil { panic(err) }
//	// close fi on exit and check for its returned error
//	defer func() {
//		if fi.Close() != nil {
//			panic(err)
//		}
//	}()
//	// make a read buffer
//	r := bufio.NewReader(fi)
//
//	return ParseContent(r)
//}

func ParseContent(reader io.Reader, batchSize int, context appengine.Context, parentKey *datastore.Key, readsBatchHandler func (context appengine.Context, userProfileKey *datastore.Key, reads []models.MeterRead) (*datastore.Key, error), carbsBatchHandler func (context appengine.Context, userProfileKey *datastore.Key, carbs []models.CarbIntake) ([] *datastore.Key, error), injectionBatchHandler func (context appengine.Context, userProfileKey *datastore.Key, injections []models.Injection) ([] *datastore.Key, error), exerciseBatchHandler func (context appengine.Context, userProfileKey *datastore.Key, exercises []models.Exercise) ([] *datastore.Key, error)) {
	decoder := xml.NewDecoder(reader)
	reads := make([]models.MeterRead,0, batchSize)
	injections := make([]models.Injection,0, batchSize)
	carbIntakes := make([]models.CarbIntake,0, batchSize)
	exercises := make([]models.Exercise,0, batchSize)

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
					meterRead := models.MeterRead{read.DisplayTime, models.TimeValue(timeutils.GetTimeInSeconds(read.DisplayTime, timeutils.TIMEZONE)), read.Value}
					reads = append(reads, meterRead)
					if (len(reads) == batchSize) {
						// Send the batch to be handled and restart another one
						readsBatchHandler(context, parentKey, reads)
						reads = make([]models.MeterRead,0, batchSize)
					}
				}
			case "Event":
				var event Event
				decoder.DecodeElement(&event, &se)
				if (event.EventType == "Carbs") {
					var carbQuantityInGrams int
					fmt.Sscanf(event.Description, "Carbs %d grams", &carbQuantityInGrams)
					carbIntake := models.CarbIntake{event.EventTime, models.TimeValue(timeutils.GetTimeInSeconds(event.EventTime, timeutils.TIMEZONE)), float32(carbQuantityInGrams), models.UNDEFINED_READ}
					carbIntakes = append(carbIntakes, carbIntake)
					if (len(carbIntakes) == batchSize) {
						// Send the batch to be handled and restart another one
						carbsBatchHandler(context, parentKey, carbIntakes)
						carbIntakes = make([]models.CarbIntake,0, batchSize)
					}
				} else if (event.EventType == "Insulin") {
					var insulinUnits float32
					_, err := fmt.Sscanf(event.Description, "Insulin %f units", &insulinUnits)
					if err != nil {
						sysutils.Propagate(err)
					}
					injection := models.Injection{event.EventTime, models.TimeValue(timeutils.GetTimeInSeconds(event.EventTime, timeutils.TIMEZONE)), float32(insulinUnits), models.UNDEFINED_READ}
					injections = append(injections, injection)
					if (len(injections) == batchSize) {
						// Send the batch to be handled and restart another one
						injectionBatchHandler(context, parentKey, injections)
						injections = make([]models.Injection,0, batchSize)
					}
				} else if (strings.HasPrefix(event.EventType, "Exercise")) {
					var duration int
					var intensity string
					fmt.Sscanf(event.Description, "Exercise %s (%d minutes)", &intensity, &duration)
					exercise := models.Exercise{event.EventTime, models.TimeValue(timeutils.GetTimeInSeconds(event.EventTime, timeutils.TIMEZONE)), duration, intensity}
					exercises = append(exercises, exercise)
					if (len(exercises) == batchSize) {
						// Send the batch to be handled and restart another one
						exerciseBatchHandler(context, parentKey, exercises)
						exercises = make([]models.Exercise,0, batchSize)
					}
				}

			case "Meter":
				// TODO: Read the meter calibrations?
			}
		}
	}

	// Run the final batch for each
	if (len(reads) > 0) {
		// Send the batch to be handled and restart another one
		readsBatchHandler(context, parentKey, reads)
	}

	if (len(injections) > 0) {
		// Send the batch to be handled and restart another one
		injectionBatchHandler(context, parentKey, injections)
	}

	if (len(carbIntakes) > 0) {
		// Send the batch to be handled and restart another one
		carbsBatchHandler(context, parentKey, carbIntakes)
	}

	if (len(exercises) > 0) {
		// Send the batch to be handled and restart another one
		exerciseBatchHandler(context, parentKey, exercises)
	}
	log.Printf("Done parsing and storing all data")
}

func ConvertAsReadsArray(meterReads *list.List) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, meterReads.Len())
	for e, i := meterReads.Front(), 0; e != nil; e, i = e.Next(), i + 1 {
		meter := e.Value.(Meter)
		reads[i] = models.MeterRead{meter.DisplayTime, models.TimeValue(timeutils.GetTimeInSeconds(meter.DisplayTime, timeutils.TIMEZONE)), meter.Value}
	}

	return reads
}
