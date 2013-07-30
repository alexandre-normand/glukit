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

// Represents a generic DataSeries structure with a series of DataPoints
type DataSeries struct {
	Name string            `json:"name"`
	Data []model.DataPoint `json:"data"`
	Type string            `json:"type"`
}

// content renders the most recent day's worth of data as json for the active user
func content(writer http.ResponseWriter, request *http.Request) {
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
	_, _, lowerBound, upperBound, err := store.GetUserData(context, email)
	if err != nil {
		util.Propagate(err)
	}

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

	writeAsJson(writer, reads, injections, carbs, exercises)
}

// writeAsJson writes the set of GlucoseReads, Injections, Carbs and Exercises as json. This is what is called from the javascript
// front-end to get the data.
func writeAsJson(writer http.ResponseWriter, reads []model.GlucoseRead, injections []model.Injection, carbs []model.Carb, exercises []model.Exercise) {
	enc := json.NewEncoder(writer)
	individuals := make([]DataSeries, 5)

	if len(reads) == 0 {
		individuals[0] = DataSeries{"You", emptyDataPointSlice, "GlucoseReads"}
		individuals[1] = DataSeries{"You.Injection", emptyDataPointSlice, "Injections"}
		individuals[2] = DataSeries{"You.Carbohydrates", emptyDataPointSlice, "Carbs"}
		individuals[3] = DataSeries{"You.Exercises", emptyDataPointSlice, "Exercises"}
		individuals[4] = DataSeries{"Perfection", emptyDataPointSlice, "ComparisonReads"}
	} else {
		individuals[0] = DataSeries{"You", model.GlucoseReadSlice(reads).ToDataPointSlice(), "GlucoseReads"}
		individuals[1] = DataSeries{"You.Injection", model.InjectionSlice(injections).ToDataPointSlice(reads), "Injections"}
		individuals[2] = DataSeries{"You.Carbohydrates", model.CarbSlice(carbs).ToDataPointSlice(reads), "Carbs"}
		individuals[3] = DataSeries{"You.Exercises", model.ExerciseSlice(exercises).ToDataPointSlice(reads), "Exercises"}
		individuals[4] = DataSeries{"Perfection", model.GlucoseReadSlice(buildPerfectBaseline(reads)).ToDataPointSlice(), "ComparisonReads"}
	}

	enc.Encode(individuals)
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
	if err != nil {
		util.Propagate(err)
	}

	reads, err := store.GetGlucoseReads(context, email, lowerBound, upperBound)
	if err != nil {
		util.Propagate(err)
	}

	writeDashboardDataAsJson(writer, request, reads, userProfile)
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
