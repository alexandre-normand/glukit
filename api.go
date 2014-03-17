package glukit

import (
	"appengine"
	"appengine/user"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/io"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"net/http"
)

type loggingCalibrationWriter struct {
	c appengine.Context
}

func (w loggingCalibrationWriter) WriteCalibrations(p []model.CalibrationRead) (n int, err error) {
	w.c.Infof("Received write [%v]", p)
	return len(p), nil
}

// processNewCalibrationData Handles a Post to the calibration endpoint and
// handles all data to be stored for a given user
func processNewCalibrationData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	_, userProfileKey, _, err := store.GetUserData(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process calibration data, user email is [%s]", user.Email)
		writer.WriteHeader(500)
		return
	}

	dataStoreWriter := io.NewDataStoreCalibrationWriter(context, userProfileKey)

	// TODO: Parse the calibration reads from the json
	calibrationBuffer := bufio.NewWriterSize(dataStoreWriter, 200)
	calibrationBuffer.WriteCalibration(model.CalibrationRead{model.Timestamp{"", 0}, 75})
	calibrationBuffer.Flush()

	// context := appengine.NewContext(request)
	// user := user.Current(context)
	writer.WriteHeader(200)
}
