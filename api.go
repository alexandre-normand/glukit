package glukit

import (
	"net/http"
)

// processNewCalibrationData Handles a Post to the calibration endpoint and
// handles all data to be stored for a given user
func processNewCalibrationData(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(200)

	// TODO: Parse the calibration reads from the json

	// TODO: Write values to the buffer

	// TODO: Flush to trigger the persistence

	// context := appengine.NewContext(request)
	// user := user.Current(context)
}
