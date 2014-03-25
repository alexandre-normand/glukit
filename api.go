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

	_, userProfileKey, _, err := store.GetUserData(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process calibration data, user email is [%s]", user.Email)
		http.Error(writer, "Error getting user to process calibration data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreCalibrationBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewWriterSize(dataStoreWriter, 200)
	calibrationStreamer := bufio.NewWriterDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var c []apimodel.Calibration
		// TODO: this is broken, fix it
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
