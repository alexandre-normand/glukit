package importer

import (
	"app/model"
	"app/util"
	"appengine"
	"appengine/datastore"
	"container/list"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
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

func ParseContent(context appengine.Context, reader io.Reader, batchSize int, parentKey *datastore.Key, startTime time.Time, readsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, carbs []model.DayOfReads) ([]*datastore.Key, error), carbsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfCarbs []model.DayOfCarbs) ([]*datastore.Key, error), injectionBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []model.DayOfInjections) ([]*datastore.Key, error), exerciseBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []model.DayOfExercises) ([]*datastore.Key, error)) (lastReadTime time.Time) {
	decoder := xml.NewDecoder(reader)

	reads := make([]model.GlucoseRead, 0)
	daysOfReads := make([]model.DayOfReads, 0, batchSize)

	injections := make([]model.Injection, 0)
	daysOfInjections := make([]model.DayOfInjections, 0, batchSize)
	var lastInjection model.Injection

	carbs := make([]model.Carb, 0)
	daysOfCarbs := make([]model.DayOfCarbs, 0, batchSize)
	var lastCarb model.Carb

	exercises := make([]model.Exercise, 0)
	daysOfExercises := make([]model.DayOfExercises, 0, batchSize)
	var lastExercise model.Exercise

	var lastRead model.GlucoseRead
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			context.Debugf("finished reading file")
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
				if read.Value > 0 {
					glucoseRead := model.GlucoseRead{read.DisplayTime, model.Timestamp(util.GetTimeInSeconds(read.InternalTime)), read.Value}

					// Skip all reads that are not after the last import's last read time
					if glucoseRead.GetTime().After(startTime) {
						// This should only happen once as we start parsing, we initialize the previous day to the current
						// and the rest of the logic should gracefully handle this case
						if len(reads) == 0 {
							lastRead = glucoseRead
						}

						// We're crossing a day boundery, we cut a batch store it and start a new one with the most recently
						// read read. This assumes that we will never get a gap big enough that two consecutive reads could
						// have the same day value while being months apart.
						if glucoseRead.GetTime().Day() != lastRead.GetTime().Day() {
							// Create a day of reads and append it to the batch
							daysOfReads = append(daysOfReads, model.DayOfReads{reads})

							if len(daysOfReads) == batchSize {
								// Send the batch to be handled and restart another one
								readsBatchHandler(context, parentKey, daysOfReads)
								daysOfReads = make([]model.DayOfReads, 0, batchSize)
							}

							reads = make([]model.GlucoseRead, 0, batchSize)
						}

						reads = append(reads, glucoseRead)
					}

					lastRead = glucoseRead
				}
			case "Event":
				var event Event
				decoder.DecodeElement(&event, &se)
				internalEventTime := util.GetTimeInSeconds(event.InternalTime)

				// Skip everything that's before the last import's read time
				if internalEventTime > startTime.Unix() {
					if event.EventType == "Carbs" {
						var carbQuantityInGrams int
						fmt.Sscanf(event.Description, "Carbs %d grams", &carbQuantityInGrams)
						carb := model.Carb{event.EventTime, model.Timestamp(internalEventTime), float32(carbQuantityInGrams), model.UNDEFINED_READ}

						if !carb.GetTime().After(startTime) {
							context.Debugf("Skipping already imported carb dated [%s]", carb.GetTime().Format(util.TIMEFORMAT))
						} else {
							// This should only happen once as we start parsing, we initialize the previous day to the current
							// and the rest of the logic should gracefully handle this case
							if len(carbs) == 0 {
								lastCarb = carb
							}

							// We're crossing a day boundery, we cut a batch store it and start a new one with the most recently
							// read carb. This assumes that we will never get a gap big enough that two consecutive carbs could
							// have the same day value while being months apart.
							if carb.GetTime().Day() != lastCarb.GetTime().Day() {
								// Create a day of reads and append it to the batch
								daysOfCarbs = append(daysOfCarbs, model.DayOfCarbs{carbs})

								if len(daysOfCarbs) == batchSize {
									// Send the batch to be handled and restart another one
									carbsBatchHandler(context, parentKey, daysOfCarbs)
									daysOfCarbs = make([]model.DayOfCarbs, 0, batchSize)
								}

								carbs = make([]model.Carb, 0)
							}

							carbs = append(carbs, carb)
						}

						lastCarb = carb
					} else if event.EventType == "Insulin" {
						var insulinUnits float32
						_, err := fmt.Sscanf(event.Description, "Insulin %f units", &insulinUnits)
						if err != nil {
							util.Propagate(err)
						}
						injection := model.Injection{event.EventTime, model.Timestamp(internalEventTime), float32(insulinUnits), model.UNDEFINED_READ}
						if !injection.GetTime().After(startTime) {
							context.Debugf("Skipping already imported injection dated [%s]", injection.GetTime().Format(util.TIMEFORMAT))
						} else {
							// This should only happen once as we start parsing, we initialize the previous day to the current
							// and the rest of the logic should gracefully handle this case
							if len(injections) == 0 {
								lastInjection = injection
							}

							// We're crossing a day boundery, we cut a batch store it and start a new one with the most recently
							// read injection. This assumes that we will never get a gap big enough that two consecutive injections could
							// have the same day value while being months apart.
							if injection.GetTime().Day() != lastInjection.GetTime().Day() {
								// Create a day of reads and append it to the batch
								daysOfInjections = append(daysOfInjections, model.DayOfInjections{injections})

								if len(daysOfInjections) == batchSize {
									// Send the batch to be handled and restart another one
									injectionBatchHandler(context, parentKey, daysOfInjections)
									daysOfInjections = make([]model.DayOfInjections, 0, batchSize)
								}

								injections = make([]model.Injection, 0)
							}

							injections = append(injections, injection)
						}

						lastInjection = injection
					} else if strings.HasPrefix(event.EventType, "Exercise") {
						var duration int
						var intensity string
						fmt.Sscanf(event.Description, "Exercise %s (%d minutes)", &intensity, &duration)
						exercise := model.Exercise{event.EventTime, model.Timestamp(internalEventTime), duration, intensity}

						if !exercise.GetTime().After(startTime) {
							context.Debugf("Skipping already imported exercise dated [%s]", exercise.GetTime().Format(util.TIMEFORMAT))
						} else {
							// This should only happen once as we start parsing, we initialize the previous day to the current
							// and the rest of the logic should gracefully handle this case
							if len(exercises) == 0 {
								lastExercise = exercise
							}

							// We're crossing a day boundery, we cut a batch store it and start a new one with the most recently
							// read exercise. This assumes that we will never get a gap big enough that two consecutive exercises could
							// have the same day value while being months apart.
							if exercise.GetTime().Day() != lastExercise.GetTime().Day() {
								// Create a day of reads and append it to the batch
								daysOfExercises = append(daysOfExercises, model.DayOfExercises{exercises})

								if len(daysOfExercises) == batchSize {
									// Send the batch to be handled and restart another one
									exerciseBatchHandler(context, parentKey, daysOfExercises)
									daysOfExercises = make([]model.DayOfExercises, 0, batchSize)
								}

								exercises = make([]model.Exercise, 0)
							}

							exercises = append(exercises, exercise)
						}

						lastExercise = exercise
					}
				}

			case "Meter":
				// TODO: Read the meter calibrations? No need for it yet but we could
			}
		}
	}

	// Run the final batch for each
	if len(reads) > 0 {
		daysOfReads = append(daysOfReads, model.DayOfReads{reads})
		context.Infof("Flushing %d days of reads", len(daysOfReads))
		readsBatchHandler(context, parentKey, daysOfReads)
	}

	if len(injections) > 0 {
		// Store the last batch
		daysOfInjections = append(daysOfInjections, model.DayOfInjections{injections})
		context.Infof("Flushing %d days of injections", len(daysOfInjections))
		injectionBatchHandler(context, parentKey, daysOfInjections)
	}

	if len(carbs) > 0 {
		// Store the last batch
		daysOfCarbs = append(daysOfCarbs, model.DayOfCarbs{carbs})
		context.Infof("Flushing %d days of carbs", len(daysOfCarbs))
		carbsBatchHandler(context, parentKey, daysOfCarbs)
	}

	if len(exercises) > 0 {
		// Store the last batch of exercises
		daysOfExercises = append(daysOfExercises, model.DayOfExercises{exercises})
		context.Infof("Flushing %d days of exercises", len(daysOfExercises))
		exerciseBatchHandler(context, parentKey, daysOfExercises)
	}

	context.Infof("Done parsing and storing all data")
	return lastRead.GetTime()
}

func ConvertAsReadsArray(glucoseReads *list.List) (reads []model.GlucoseRead) {
	reads = make([]model.GlucoseRead, glucoseReads.Len())
	for e, i := glucoseReads.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		meter := e.Value.(Meter)
		reads[i] = model.GlucoseRead{meter.DisplayTime, model.Timestamp(util.GetTimeInSeconds(meter.InternalTime)), meter.Value}
	}

	return reads
}
