package glukit

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"io"
	"net/http"
	"time"
)

const (
	GLUCOSEREADS_V1_ROUTE = "v1_glucosereads"
	CALIBRATIONS_V1_ROUTE = "v1_calibrations"
	EXERCISES_V1_ROUTE    = "v1_exercises"
	MEALS_V1_ROUTE        = "v1_meals"
	INJECTIONS_V1_ROUTE   = "v1_injections"
)

func initApiEndpoints(writer http.ResponseWriter, request *http.Request) {
	muxRouter.Get(CALIBRATIONS_V1_ROUTE).Handler(newOauthAuthenticationHandler(http.HandlerFunc(processNewCalibrationData)))
	muxRouter.Get(INJECTIONS_V1_ROUTE).Handler(newOauthAuthenticationHandler(http.HandlerFunc(processNewInjectionData)))
	muxRouter.Get(MEALS_V1_ROUTE).Handler(newOauthAuthenticationHandler(http.HandlerFunc(processNewMealData)))
	muxRouter.Get(GLUCOSEREADS_V1_ROUTE).Handler(newOauthAuthenticationHandler(http.HandlerFunc(processNewGlucoseReadData)))
	muxRouter.Get(EXERCISES_V1_ROUTE).Handler(newOauthAuthenticationHandler(http.HandlerFunc(processNewExerciseData)))
}

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
	calibrationStreamer := streaming.NewCalibrationReadStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var c []model.CalibrationRead

		if err = decoder.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error processing calibration data for user [%s]: %v", user.Email, err)
			break
		}

		context.Debugf("Writing new calibration reads [%v]", c)
		calibrationStreamer, err = calibrationStreamer.WriteCalibrations(c)
		if err != nil {
			context.Warningf("Error storing calibration data [%v]: %v", c, err)
			http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
			return
		}
	}

	if err != io.EOF {
		context.Warningf("Error processing calibration read data for user [%s]: %v", user.Email, err)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	calibrationStreamer, err = calibrationStreamer.Close()
	if err != nil {
		context.Warningf("Error closing calibration read streamer: %v", err)
		http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
		return
	}

	context.Infof("Wrote calibrations to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

// processNewGlucoseReadData Handles a Post to the glucosereads endpoint and
// handles all data to be stored for a given user
func processNewGlucoseReadData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process glucose read data, user email is [%s]: %v", user.Email, err)
		http.Error(writer, "Error getting user to process glucose read data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreGlucoseReadBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewGlucoseReadWriterSize(dataStoreWriter, 200)
	glucoseReadStreamer := streaming.NewGlucoseStreamerDuration(batchingWriter, time.Hour*24)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(request.Body)

	// context.Debugf("Body is [%s]", buf.String())
	decoder := json.NewDecoder(request.Body)

	for {
		var c []model.GlucoseRead

		if err = decoder.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error processing glucose read data for user [%s]: %v", user.Email, err)
			break
		}

		context.Debugf("Writing [%d] new glucose reads", len(c))
		glucoseReadStreamer, err = glucoseReadStreamer.WriteGlucoseReads(c)
		if err != nil {
			context.Warningf("Error storing user data [%v]: %v", c, err)
			http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
			return
		}
	}

	if err != io.EOF {
		context.Warningf("Error processing glucose read data for user [%s]: %v", user.Email, err)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	glucoseReadStreamer, err = glucoseReadStreamer.Close()
	if err != nil {
		context.Warningf("Error closing glucose read streamer: %v", err)
		http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
		return
	}

	context.Infof("Wrote glucose reads to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

// processNewInjectionData Handles a Post to the injections endpoint and
// handles all data to be stored for a given user
func processNewInjectionData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process injection data, user email is [%s]: %v", user.Email, err)
		http.Error(writer, "Error getting user to process injection data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreInjectionBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewInjectionWriterSize(dataStoreWriter, 200)
	injectionStreamer := streaming.NewInjectionStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var p []model.Injection

		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error processing injection data for user [%s]: %v", user.Email, err)
			break
		}

		context.Debugf("Writing [%d] new injections", len(p))
		injectionStreamer, err = injectionStreamer.WriteInjections(p)
		if err != nil {
			context.Warningf("Error storing injection data [%v]: %v", p, err)
			http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
			return
		}
	}

	if err != io.EOF {
		context.Warningf("Error processing injection data for user [%s]: %v", user.Email, err)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	injectionStreamer, err = injectionStreamer.Close()
	if err != nil {
		context.Warningf("Error closing injection streamer: %v", err)
		http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
		return
	}

	context.Infof("Wrote injections to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

// processNewMealData Handles a Post to the Meals endpoint and
// handles all data to be stored for a given user
func processNewMealData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process meal data, user email is [%s]: %v", user.Email, err)
		http.Error(writer, "Error getting user to process meal data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreMealBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewMealWriterSize(dataStoreWriter, 200)
	mealStreamer := streaming.NewMealStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var meals []model.Meal

		if err = decoder.Decode(&meals); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error processing meal data for user [%s]: %v", user.Email, err)
			break
		}

		context.Debugf("Writing [%d] new meals", len(meals))
		mealStreamer, err = mealStreamer.WriteMeals(meals)
		if err != nil {
			context.Warningf("Error storing meal data [%v]: %v", meals, err)
			http.Error(writer, fmt.Sprintf("Error storing meal data: %v", err), 502)
			return
		}
	}

	if err != io.EOF {
		context.Warningf("Error processing meal data for user [%s]: %v", user.Email, err)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	mealStreamer, err = mealStreamer.Close()
	if err != nil {
		context.Warningf("Error closing meal streamer: %v", err)
		http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
		return
	}

	context.Infof("Wrote meals to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}

// processNewExerciseData Handles a Post to the exercises endpoint and
// handles all data to be stored for a given user
func processNewExerciseData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	userProfileKey, _, err := store.GetGlukitUser(context, user.Email)
	if err != nil {
		context.Warningf("Error getting user to process exercise data, user email is [%s]: %v", user.Email, err)
		http.Error(writer, "Error getting user to process exercise data", 500)
		return
	}

	dataStoreWriter := store.NewDataStoreExerciseBatchWriter(context, userProfileKey)
	batchingWriter := bufio.NewExerciseWriterSize(dataStoreWriter, 200)
	exerciseStreamer := streaming.NewExerciseStreamerDuration(batchingWriter, time.Hour*24)

	decoder := json.NewDecoder(request.Body)

	for {
		var exercises []model.Exercise

		if err = decoder.Decode(&exercises); err == io.EOF {
			break
		} else if err != nil {
			context.Warningf("Error processing exercise data for user [%s]: %v", user.Email, err)
			break
		}

		context.Debugf("Writing [%d] new Exercises", len(exercises))
		exerciseStreamer, err = exerciseStreamer.WriteExercises(exercises)
		if err != nil {
			context.Warningf("Error storing exercise data [%v]: %v", exercises, err)
			http.Error(writer, fmt.Sprintf("Error storing exercise data: %v", err), 502)
			return
		}
	}

	if err != io.EOF {
		context.Warningf("Error processing exercise data for user [%s]: %v", user.Email, err)
		http.Error(writer, fmt.Sprintf("Error decoding data: %v", err), 400)
		return
	}

	exerciseStreamer, err = exerciseStreamer.Close()
	if err != nil {
		context.Warningf("Error closing exercise streamer: %v", err)
		http.Error(writer, fmt.Sprintf("Error storing data: %v", err), 502)
		return
	}

	context.Infof("Wrote exercises to the datastore for user [%s]", user.Email)
	writer.WriteHeader(200)
}
