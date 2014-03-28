package glukit

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"io"
	"net/http"
	"time"
)

// processNewCalibrationData Handles a Post to the calibration endpoint and
// handles all data to be stored for a given user
func processNewCalibrationData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, userProfileKey, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process calibration data, user email is [%s]: %v", user.Email, err)
		http.Error(writer, "Error getting user to process calibration data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreCalibrationBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewCalibrationWriterSize(dataStoreWriter, 200)
	calibrationStreamer := bufio.NewCalibrationReadStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var c []apimodel.Calibration

		if err = decoder.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error getting user to process calibration data, user email is [%s]", user.Email)
			http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
			break
		}

		reads := convertToCalibrationRead(c)
		context.Debugf("Writing new calibration reads [%v]", reads)
		calibrationStreamer.WriteCalibrations(reads)
	}

	if err != io.EOF {
		context.Warningf("Error getting user to process calibration data, user email is [%s]", user.Email)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	calibrationStreamer.Flush()
	context.Infof("Wrote calibrations to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

func convertToCalibrationRead(p []apimodel.Calibration) []model.CalibrationRead {
	r := make([]model.CalibrationRead, len(p))
	for i, c := range p {
		r[i] = model.CalibrationRead{model.Timestamp{c.DisplayTime, util.GetTimeInSeconds(c.InternalTime)}, c.Value}
	}
	return r
}

// processNewGlucoseReadData Handles a Post to the glucoseread endpoint and
// handles all data to be stored for a given user
func processNewGlucoseReadData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, userProfileKey, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process glucose read data, user email is [%s]", user.Email, err)
		http.Error(writer, "Error getting user to process glucose read data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreGlucoseReadBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewGlucoseReadWriterSize(dataStoreWriter, 200)
	glucoseReadStreamer := bufio.NewGlucoseStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var c []apimodel.Glucose

		if err = decoder.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error getting user to process glucose read data, user email is [%s]", user.Email)
			http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
			break
		}

		reads := convertToGlucoseRead(c)
		context.Debugf("Writing [%d] new glucose reads", len(reads))
		glucoseReadStreamer.WriteGlucoseReads(reads)
	}

	if err != io.EOF {
		context.Warningf("Error getting user to process glucose read data, user email is [%s]", user.Email)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	glucoseReadStreamer.Flush()
	context.Infof("Wrote glucose reads to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

func convertToGlucoseRead(p []apimodel.Glucose) []model.GlucoseRead {
	r := make([]model.GlucoseRead, len(p))
	for i, c := range p {
		r[i] = model.GlucoseRead{model.Timestamp{c.DisplayTime, util.GetTimeInSeconds(c.InternalTime)}, c.Value}
	}
	return r
}
