package glukit

import (
	"appengine"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
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

	writer.WriteHeader(200)

	context := appengine.NewContext(request)
	logWriter := &loggingCalibrationWriter{context}

	calibrationBuffer := bufio.NewWriterSize(logWriter, 200)
	calibrationBuffer.WriteCalibration(model.CalibrationRead{model.Timestamp{"", 0}, 75})
	// TODO: Parse the calibration reads from the json

	// TODO: Write values to the buffer

	// TODO: Flush to trigger the persistence

	// context := appengine.NewContext(request)
	// user := user.Current(context)
}
