package importer

import (
	"appengine"
	"appengine/datastore"
	"encoding/xml"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"github.com/alexandre-normand/glukit/app/util"
	"io"
	"strings"
	"time"
)

// ParseContent is the big function that parses the Dexcom xml file. It is given a reader to the file and it parses batches of days of GlucoseReads/Events. It streams the content but
// keeps some in memory until it reaches a full batch of a type. A batch is an array of DayOf[GlucoseReads,Injection,Carbs,Exercises]. A batch is flushed to the datastore once it reaches
// the given batchSize or we reach the end of the file.
func ParseContent(context appengine.Context, reader io.Reader, batchSize int, parentKey *datastore.Key, startTime time.Time, readsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, carbs []model.DayOfGlucoseReads) ([]*datastore.Key, error), carbsBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfCarbs []model.DayOfCarbs) ([]*datastore.Key, error), injectionBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []model.DayOfInjections) ([]*datastore.Key, error), exerciseBatchHandler func(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []model.DayOfExercises) ([]*datastore.Key, error)) (lastReadTime time.Time, err error) {
	decoder := xml.NewDecoder(reader)

	calibrationDataStoreWriter := store.NewDataStoreCalibrationBatchWriter(context, parentKey)
	calibrationBatchingWriter := bufio.NewCalibrationWriterSize(calibrationDataStoreWriter, IMPORT_BATCH_SIZE)
	calibrationStreamer := streaming.NewCalibrationReadStreamerDuration(calibrationBatchingWriter, time.Hour*24)

	glucoseDataStoreWriter := store.NewDataStoreGlucoseReadBatchWriter(context, parentKey)
	glucoseBatchingWriter := bufio.NewGlucoseReadWriterSize(glucoseDataStoreWriter, IMPORT_BATCH_SIZE)
	glucoseStreamer := streaming.NewGlucoseStreamerDuration(glucoseBatchingWriter, time.Hour*24)

	injectionDataStoreWriter := store.NewDataStoreInjectionBatchWriter(context, parentKey)
	injectionBatchingWriter := bufio.NewInjectionWriterSize(injectionDataStoreWriter, IMPORT_BATCH_SIZE)
	injectionStreamer := streaming.NewInjectionStreamerDuration(injectionBatchingWriter, time.Hour*24)

	carbDataStoreWriter := store.NewDataStoreCarbBatchWriter(context, parentKey)
	carbBatchingWriter := bufio.NewCarbWriterSize(carbDataStoreWriter, IMPORT_BATCH_SIZE)
	carbStreamer := streaming.NewCarbStreamerDuration(carbBatchingWriter, time.Hour*24)

	exerciseDataStoreWriter := store.NewDataStoreExerciseBatchWriter(context, parentKey)
	exerciseBatchingWriter := bufio.NewExerciseWriterSize(exerciseDataStoreWriter, IMPORT_BATCH_SIZE)
	exerciseStreamer := streaming.NewExerciseStreamerDuration(exerciseBatchingWriter, time.Hour*24)

	var lastRead *model.GlucoseRead
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
				var read apimodel.Glucose
				// decode a whole chunk of following XML into the
				decoder.DecodeElement(&read, &se)

				if read.Value > 0 {
					glucoseRead := model.GlucoseRead{model.Timestamp{read.DisplayTime, util.GetTimeInSeconds(read.InternalTime)}, read.Value}
					glucoseStreamer, err = glucoseStreamer.WriteGlucoseRead(glucoseRead)

					if err != nil {
						return lastRead.GetTime(), err
					}

					lastRead = &glucoseRead
				}
			case "Event":
				var event apimodel.Event
				decoder.DecodeElement(&event, &se)
				internalEventTime := util.GetTimeInSeconds(event.InternalTime)

				// Skip everything that's before the last import's read time
				if internalEventTime > startTime.Unix() {
					if event.EventType == "Carbs" {
						var carbQuantityInGrams int
						fmt.Sscanf(event.Description, "Carbs %d grams", &carbQuantityInGrams)
						carb := model.Carb{model.Timestamp{event.EventTime, internalEventTime}, float32(carbQuantityInGrams), model.UNDEFINED_READ}

						carbStreamer, err = carbStreamer.WriteCarb(carb)
						if err != nil {
							return lastRead.GetTime(), err
						}

					} else if event.EventType == "Insulin" {
						var insulinUnits float32
						_, err := fmt.Sscanf(event.Description, "Insulin %f units", &insulinUnits)
						if err != nil {
							context.Warningf("Failed to parse event as injection [%s]: %v", event.Description, err)
						} else {
							injection := model.Injection{model.Timestamp{event.EventTime, internalEventTime}, float32(insulinUnits), model.UNDEFINED_READ}
							injectionStreamer, err = injectionStreamer.WriteInjection(injection)

							if err != nil {
								return lastRead.GetTime(), err
							}
						}
					} else if strings.HasPrefix(event.EventType, "Exercise") {
						var duration int
						var intensity string
						fmt.Sscanf(event.Description, "Exercise %s (%d minutes)", &intensity, &duration)
						exercise := model.Exercise{model.Timestamp{event.EventTime, internalEventTime}, duration, intensity}

						exerciseStreamer, err = exerciseStreamer.WriteExercise(exercise)
						if err != nil {
							return lastRead.GetTime(), err
						}
					}
				}

			case "Meter":
				var c apimodel.Calibration
				decoder.DecodeElement(&c, &se)
				calibrationStreamer, err = calibrationStreamer.WriteCalibration(model.CalibrationRead{model.Timestamp{c.DisplayTime, util.GetTimeInSeconds(c.InternalTime)}, c.Value})

				if err != nil {
					return lastRead.GetTime(), err
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

	carbStreamer, err = carbStreamer.Close()
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
