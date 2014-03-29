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

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
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

// processNewGlucoseReadData Handles a Post to the glucosereads endpoint and
// handles all data to be stored for a given user
func processNewGlucoseReadData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
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

	glucoseReadStreamer.Close()
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

// processNewInjectionData Handles a Post to the injections endpoint and
// handles all data to be stored for a given user
func processNewInjectionData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process injection data, user email is [%s]", user.Email, err)
		http.Error(writer, "Error getting user to process injection data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreInjectionBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewInjectionWriterSize(dataStoreWriter, 200)
	injectionStreamer := bufio.NewInjectionStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var p []apimodel.Injection

		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error getting user to process injection data, user email is [%s]", user.Email)
			http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
			break
		}

		injections := convertToInjection(p)
		context.Debugf("Writing [%d] new injections", len(injections))
		injectionStreamer.WriteInjections(injections)
	}

	if err != io.EOF {
		context.Warningf("Error getting user to process injection data, user email is [%s]", user.Email)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	injectionStreamer.Close()
	context.Infof("Wrote injections to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

func convertToInjection(p []apimodel.Injection) []model.Injection {
	r := make([]model.Injection, len(p))
	for i, e := range p {
		r[i] = model.Injection{model.Timestamp{e.EventTime, util.GetTimeInSeconds(e.InternalTime)}, e.Units, model.UNDEFINED_READ}
	}
	return r
}

// processNewMealData Handles a Post to the Meals endpoint and
// handles all data to be stored for a given user
func processNewMealData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process meal data, user email is [%s]", user.Email, err)
		http.Error(writer, "Error getting user to process meal data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreCarbBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewCarbWriterSize(dataStoreWriter, 200)
	carbStreamer := bufio.NewCarbStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var p []apimodel.Meal

		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error getting user to process meal data, user email is [%s]", user.Email)
			http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
			break
		}

		meals := convertToMeal(p)
		context.Debugf("Writing [%d] new meals", len(meals))
		carbStreamer.WriteCarbs(meals)
	}

	if err != io.EOF {
		context.Warningf("Error getting user to process meal data, user email is [%s]", user.Email)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	carbStreamer.Close()
	context.Infof("Wrote meals to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

func convertToMeal(p []apimodel.Meal) []model.Carb {
	r := make([]model.Carb, len(p))
	for i, e := range p {
		r[i] = model.Carb{model.Timestamp{e.EventTime, util.GetTimeInSeconds(e.InternalTime)}, e.Carbohydrates, model.UNDEFINED_READ}
	}
	return r
}
