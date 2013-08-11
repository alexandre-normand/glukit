package glukit

import (
	"app/engine"
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/user"
	"encoding/json"
	"lib/github.com/grd/stat"
	"net/http"
	"sort"
)

// Represents a DataResponse with an array of DataSeries and some metadata
type DataResponse struct {
	Name      string             `json:"name"`
	Score     *model.GlukitScore `json:"score"`
	Data      []DataSeries       `json:"data"`
	AvatarUrl string             `json:"avatar"`
}

// Represents a generic DataSeries structure with a series of DataPoints
type DataSeries struct {
	Name string            `json:"name"`
	Data []model.DataPoint `json:"data"`
	Type string            `json:"type"`
}

// content renders the most recent day's worth of data as json for the active user
func personalData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	mostRecentDayAsJson(writer, request, user.Email)
}

// demoContent renders the most recent day's worth of data as json for the demo user
func demoContent(writer http.ResponseWriter, request *http.Request) {
	mostRecentDayAsJson(writer, request, DEMO_EMAIL)
}

// mostRecentDayAsJson retrieves the most recent day's worth of data for the user identified by
// the given email address and writes to the response writer as json
func mostRecentDayAsJson(writer http.ResponseWriter, request *http.Request, email string) {
	context := appengine.NewContext(request)
	glukitUser, _, lowerBound, upperBound, err := store.GetUserData(context, email)
	if err != nil && err == store.ErrNoImportedDataFound {
		context.Debugf("No imported data found for user [%s]", email)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}
		injections, err := store.GetInjections(context, email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}
		carbs, err := store.GetCarbs(context, email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}
		exercises, err := store.GetExercises(context, email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}

		value := writer.Header()
		value.Add("Content-type", "application/json")

		response := DataResponse{Name: glukitUser.FirstName, Score: &glukitUser.Score, Data: generateDataSeriesFromData(reads, injections, carbs, exercises), AvatarUrl: ""}
		writeAsJson(writer, response)
	}
}

// find the steady sailor and retrieve his most recent day's worth of data.
func steadySailorData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	steadySailorDataForEmail(writer, request, user.Email)
}

// find the steady sailor for the demo user and retrieve his most recent day's worth of data.
func demoSteadySailorData(writer http.ResponseWriter, request *http.Request) {
	steadySailorDataForEmail(writer, request, DEMO_EMAIL)
}

// find the steady sailor and retrieve his most recent day's worth of data.
func steadySailorDataForEmail(writer http.ResponseWriter, request *http.Request, recipientEmail string) {
	context := appengine.NewContext(request)
	steadySailor, _, lowerBound, upperBound, err := store.FindSteadySailor(context, recipientEmail)

	if err != nil && err == store.ErrNoSteadySailorMatchFound {
		context.Debugf("No steady sailor match found for user [%s]", recipientEmail)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, steadySailor.Email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}

		value := writer.Header()
		value.Add("Content-type", "application/json")

		response := DataResponse{Name: steadySailor.FirstName, Score: &steadySailor.Score, Data: generateDataSeriesFromData(reads, nil, nil, nil), AvatarUrl: ""}
		writeAsJson(writer, response)
	}
}

// writeAsJson writes a DataResponse with its set of GlucoseReads, Injections, Carbs and Exercises as json. This is what is called from the javascript
// front-end to get the data.
func writeAsJson(writer http.ResponseWriter, response DataResponse) {
	enc := json.NewEncoder(writer)
	enc.Encode(response)
}

func generateDataSeriesFromData(reads []model.GlucoseRead, injections []model.Injection, carbs []model.Carb, exercises []model.Exercise) (dataSeries []DataSeries) {
	data := make([]DataSeries, 1)

	data[0] = DataSeries{"GlucoseReads", model.GlucoseReadSlice(reads).ToDataPointSlice(), "GlucoseReads"}
	if injections != nil {
		data = append(data, DataSeries{"Injections", model.InjectionSlice(injections).ToDataPointSlice(reads), "Injections"})
	}

	if carbs != nil {
		data = append(data, DataSeries{"Carbohydrates", model.CarbSlice(carbs).ToDataPointSlice(reads), "Carbs"})
	}

	if exercises != nil {
		data = append(data, DataSeries{"Exercises", model.ExerciseSlice(exercises).ToDataPointSlice(reads), "Exercises"})
	}

	return data
}

// dashboard renders the dashboard statistics as json
func dashboard(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	dashboardDataForUser(writer, request, user.Email)
}

// demodashboard renders the dashboard statistics as json for the demo user
func demoDashboard(writer http.ResponseWriter, request *http.Request) {
	dashboardDataForUser(writer, request, DEMO_EMAIL)
}

// dashboardDataForUser retrieves reads and generates dashboard statistics from them
func dashboardDataForUser(writer http.ResponseWriter, request *http.Request, email string) {
	context := appengine.NewContext(request)

	userProfile, _, lowerBound, upperBound, err := store.GetUserData(context, email)
	if err != nil && err == store.ErrNoImportedDataFound {
		context.Debugf("No imported data found for user [%s]", email)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, email, lowerBound, upperBound)
		if err != nil {
			util.Propagate(err)
		}

		writeDashboardDataAsJson(writer, request, reads, userProfile)
	}
}

// writedashboardDataAsJson calculates dashboard statistics from an array of GlucoseReads and writes it
// as json
func writeDashboardDataAsJson(writer http.ResponseWriter, request *http.Request, reads []model.GlucoseRead, userProfile *model.GlukitUser) {
	value := writer.Header()
	value.Add("Content-type", "application/json")

	var dashboardData model.DashboardData
	if len(reads) > 0 {
		sort.Sort(model.ReadStatsSlice(reads))
		dashboardData.Average = stat.Mean(model.ReadStatsSlice(reads))
		dashboardData.High, _ = stat.Max(model.ReadStatsSlice(reads))
		dashboardData.Low, _ = stat.Min(model.ReadStatsSlice(reads))
		dashboardData.Median = stat.MedianFromSortedData(model.ReadStatsSlice(reads))
		dashboardData.Score = engine.CalculateUserFacingScore(userProfile.Score)
	}

	enc := json.NewEncoder(writer)
	enc.Encode(dashboardData)
}
