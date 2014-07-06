package importer

import (
	"appengine"
	"appengine/datastore"
	"encoding/xml"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/dexcomimporter"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"github.com/alexandre-normand/glukit/app/util"
	"io"
	"strings"
	"time"
)

// ParseContent is the big function that parses the Dexcom xml file. It is given a reader to the file and it parses batches of days of GlucoseReads/Events. It streams the content but
// keeps some in memory until it reaches a full batch of a type. A batch is an array of DayOf[GlucoseReads,Injection,Meals,Exercises]. A batch is flushed to the datastore once it reaches
// the given batchSize or we reach the end of the file.
func ParseContent(context appengine.Context, reader io.Reader, parentKey *datastore.Key, startTime time.Time, readsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, meals []apimodel.DayOfGlucoseReads) ([]*datastore.Key, error), mealsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfMeals []apimodel.DayOfMeals) ([]*datastore.Key, error), injectionBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []apimodel.DayOfInjections) ([]*datastore.Key, error), exerciseBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []apimodel.DayOfExercises) ([]*datastore.Key, error)) (lastReadTime time.Time, err error) {
	decoder := xml.NewDecoder(reader)

	calibrationDataStoreWriter := store.NewDataStoreCalibrationBatchWriter(context, parentKey)
	calibrationBatchingWriter := bufio.NewCalibrationWriterSize(calibrationDataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	calibrationStreamer := streaming.NewCalibrationReadStreamerDuration(calibrationBatchingWriter, apimodel.DAY_OF_DATA_DURATION)

	glucoseDataStoreWriter := store.NewDataStoreGlucoseReadBatchWriter(context, parentKey)
	glucoseBatchingWriter := bufio.NewGlucoseReadWriterSize(glucoseDataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	glucoseStreamer := streaming.NewGlucoseStreamerDuration(glucoseBatchingWriter, apimodel.DAY_OF_DATA_DURATION)

	injectionDataStoreWriter := store.NewDataStoreInjectionBatchWriter(context, parentKey)
	injectionBatchingWriter := bufio.NewInjectionWriterSize(injectionDataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	injectionStreamer := streaming.NewInjectionStreamerDuration(injectionBatchingWriter, apimodel.DAY_OF_DATA_DURATION)

	mealDataStoreWriter := store.NewDataStoreMealBatchWriter(context, parentKey)
	mealBatchingWriter := bufio.NewMealWriterSize(mealDataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	mealStreamer := streaming.NewMealStreamerDuration(mealBatchingWriter, apimodel.DAY_OF_DATA_DURATION)

	exerciseDataStoreWriter := store.NewDataStoreExerciseBatchWriter(context, parentKey)
	exerciseBatchingWriter := bufio.NewExerciseWriterSize(exerciseDataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	exerciseStreamer := streaming.NewExerciseStreamerDuration(exerciseBatchingWriter, apimodel.DAY_OF_DATA_DURATION)

	var lastRead *apimodel.GlucoseRead
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
				var read dexcomimporter.Glucose
				// decode a whole chunk of following XML into the
				decoder.DecodeElement(&read, &se)
				glucoseRead, err := dexcomimporter.ConvertXmlGlucoseRead(read)
				if err != nil {
					return lastRead.GetTime(), err
				}

				if glucoseRead != nil && glucoseRead.Value > 0 {
					glucoseStreamer, err = glucoseStreamer.WriteGlucoseRead(*glucoseRead)

					if err != nil {
						return lastRead.GetTime(), err
					}

					lastRead = glucoseRead
				}
			case "Event":
				var event dexcomimporter.Event
				decoder.DecodeElement(&event, &se)
				internalEventTime, err := util.GetTimeUTC(event.InternalTime)
				if err != nil {
					context.Warningf("Skipping [%s] event [%v], bad internal time [%s]: %v", event.EventType, event, event.InternalTime, err)
					continue
				}

				// Skip everything that's before the last import's read time
				if internalEventTime.Unix() > startTime.Unix() {
					location := util.GetLocaltimeOffset(event.EventTime, internalEventTime)

					eventTime, err := util.GetTimeWithImpliedLocation(event.EventTime, location)
					if err != nil {
						context.Warningf("Skipping [%s] event [%v], bad event time [%s]: %v", event.EventType, event, event.EventTime, err)
						continue
					}

					if event.EventType == "Carbs" {
						var mealQuantityInGrams int
						fmt.Sscanf(event.Description, "Carbs %d grams", &mealQuantityInGrams)

						meal := apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(eventTime), location.String()}, float32(mealQuantityInGrams), 0., 0., 0.}

						mealStreamer, err = mealStreamer.WriteMeal(meal)
						if err != nil {
							return lastRead.GetTime(), err
						}

					} else if event.EventType == "Insulin" {
						var insulinUnits float32
						_, err := fmt.Sscanf(event.Description, "Insulin %f units", &insulinUnits)
						if err != nil {
							context.Warningf("Failed to parse event as injection [%s]: %v", event.Description, err)
						} else {
							injection := apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(eventTime), location.String()}, float32(insulinUnits), "", ""}

							injectionStreamer, err = injectionStreamer.WriteInjection(injection)

							if err != nil {
								return lastRead.GetTime(), err
							}
						}
					} else if strings.HasPrefix(event.EventType, "Exercise") {
						var duration int
						var intensity string
						fmt.Sscanf(event.Description, "Exercise %s (%d minutes)", &intensity, &duration)

						exercise := apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(eventTime), location.String()}, duration, intensity, ""}
						exerciseStreamer, err = exerciseStreamer.WriteExercise(exercise)
						if err != nil {
							return lastRead.GetTime(), err
						}
					}
				}
			case "Meter":
				var c dexcomimporter.Calibration
				decoder.DecodeElement(&c, &se)

				if calibrationRead, err := dexcomimporter.ConvertXmlCalibrationRead(c); err != nil {
					return lastRead.GetTime(), err
				} else {
					calibrationStreamer, err = calibrationStreamer.WriteCalibration(*calibrationRead)

					if err != nil {
						return lastRead.GetTime(), err
					}
				}
			}
		}
	}

	// Close the streams and flush anything pending
	glucoseStreamer, err = glucoseStreamer.Close()
	if err != nil {
		return lastRead.GetTime(), err
	}
	calibrationStreamer, err = calibrationStreamer.Close()
	if err != nil {
		return lastRead.GetTime(), err
	}

	injectionStreamer, err = injectionStreamer.Close()
	if err != nil {
		return lastRead.GetTime(), err
	}

	mealStreamer, err = mealStreamer.Close()
	if err != nil {
		return lastRead.GetTime(), err
	}

	exerciseStreamer, err = exerciseStreamer.Close()
	if err != nil {
		return lastRead.GetTime(), err
	}

	context.Infof("Done parsing and storing all data")
	return lastRead.GetTime(), nil
}
