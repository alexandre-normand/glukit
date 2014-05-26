package glukit

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/github.com/grd/stat"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Represents a DataResponse with an array of DataSeries and some metadata
type DataResponse struct {
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Data      []DataSeries `json:"data"`
	Picture   string       `json:"picture"`
	LastSync  time.Time    `json:"lastSync"`
}

// Represents a generic DataSeries structure with a series of DataPoints
type DataSeries struct {
	Name string            `json:"name"`
	Data []model.DataPoint `json:"data"`
	Type string            `json:"type"`
}

const (
	QUERY_PARAM_LIMIT = "limit"
	QUERY_PARAM_FROM  = "from"
	QUERY_PARAM_TO    = "to"
)

// content renders the most recent day's worth of data as json for the active user
func personalData(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	mostRecentWeekAsJson(writer, request, user.Email)
}

// demoContent renders the most recent day's worth of data as json for the demo user
func demoContent(writer http.ResponseWriter, request *http.Request) {
	mostRecentWeekAsJson(writer, request, DEMO_EMAIL)
}

// mostRecentWeekAsJson retrieves the most recent day's week worth of data for the user identified by
// the given email address and writes to the response writer as json
func mostRecentWeekAsJson(writer http.ResponseWriter, request *http.Request, email string) {
	context := appengine.NewContext(request)
	glukitUser, _, upperBound, err := store.GetUserData(context, email)
	lowerBound := upperBound.Add(model.DEFAULT_LOOKBACK_PERIOD)

	if err != nil && err == store.ErrNoImportedDataFound {
		context.Debugf("No imported data found for user [%s]", email)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}
		injections, err := store.GetInjections(context, email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}
		carbs, err := store.GetMeals(context, email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}
		exercises, err := store.GetExercises(context, email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}

		value := writer.Header()
		value.Add("Content-type", "application/json")

		response := DataResponse{FirstName: glukitUser.FirstName, LastName: glukitUser.LastName, Data: generateDataSeriesFromData(reads, injections, carbs, exercises), Picture: glukitUser.PictureUrl, LastSync: glukitUser.MostRecentRead.GetTime()}
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
	steadySailor, _, upperBound, err := store.FindSteadySailor(context, recipientEmail)

	lowerBound := upperBound.Add(model.DEFAULT_LOOKBACK_PERIOD)
	if err != nil && err == store.ErrNoSteadySailorMatchFound {
		context.Debugf("No steady sailor match found for user [%s]", recipientEmail)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, steadySailor.Email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}

		value := writer.Header()
		value.Add("Content-type", "application/json")

		response := DataResponse{FirstName: steadySailor.FirstName, LastName: steadySailor.LastName, Data: generateDataSeriesFromData(reads, nil, nil, nil), Picture: steadySailor.PictureUrl, LastSync: steadySailor.MostRecentRead.GetTime()}
		writeAsJson(writer, response)
	}
}

// writeAsJson writes a DataResponse with its set of GlucoseReads, Injections, Meals and Exercises as json. This is what is called from the javascript
// front-end to get the data.
func writeAsJson(writer http.ResponseWriter, response DataResponse) {
	enc := json.NewEncoder(writer)
	enc.Encode(response)
}

func generateDataSeriesFromData(reads []model.GlucoseRead, injections []model.Injection, carbs []model.Meal, exercises []model.Exercise) (dataSeries []DataSeries) {
	data := make([]DataSeries, 1)

	data[0] = DataSeries{"GlucoseReads", model.GlucoseReadSlice(reads).ToDataPointSlice(), "GlucoseReads"}
	var userEvents []model.DataPoint
	if injections != nil {
		userEvents = model.MergeDataPointArrays(userEvents, model.InjectionSlice(injections).ToDataPointSlice(reads))
	}

	if carbs != nil {
		userEvents = model.MergeDataPointArrays(userEvents, model.MealSlice(carbs).ToDataPointSlice(reads))
	}

	// TODO: clean up exercise from all the app or restore it. We won't be using it at the moment as we don't think the exercise data
	// from the dexcom is good enough
	// if exercises != nil {
	// 	userEvents = model.MergeDataPointArrays(userEvents, model.ExerciseSlice(exercises).ToDataPointSlice(reads))
	// }

	sort.Sort(model.DataPointSlice(userEvents))

	data = append(data, DataSeries{"UserEvents", userEvents, "UserEvents"})

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

	userProfile, _, upperBound, err := store.GetUserData(context, email)
	lowerBound := upperBound.Add(time.Duration(-1*24) * time.Hour)

	if err != nil && err == store.ErrNoImportedDataFound {
		context.Debugf("No imported data found for user [%s]", email)
		http.Error(writer, err.Error(), 204)
	} else if err != nil {
		util.Propagate(err)
	} else {
		reads, err := store.GetGlucoseReads(context, email, lowerBound, time.Now())
		if err != nil {
			util.Propagate(err)
		}

		writeDashboardDataAsJson(writer, request, reads, userProfile)
	}
}

// writedashboardDataAsJson calculates dashboard statistics from an array of GlucoseReads and writes it
// as json
func writeDashboardDataAsJson(writer http.ResponseWriter, request *http.Request, reads []model.GlucoseRead, userProfile *model.GlukitUser) {
	context := appengine.NewContext(request)
	value := writer.Header()
	value.Add("Content-type", "application/json")

	var dashboardData model.DashboardData
	if len(reads) > 0 {
		sort.Sort(model.ReadStatsSlice(reads))
		dashboardData.FirstName = userProfile.FirstName
		dashboardData.LastName = userProfile.LastName
		dashboardData.Picture = userProfile.PictureUrl
		dashboardData.LastSync = userProfile.MostRecentRead.GetTime()
		dashboardData.Average = stat.Mean(model.ReadStatsSlice(reads))
		dashboardData.High, _ = stat.Max(model.ReadStatsSlice(reads))
		dashboardData.Low, _ = stat.Min(model.ReadStatsSlice(reads))
		dashboardData.Median = stat.MedianFromSortedData(model.ReadStatsSlice(reads))
		dashboardData.Score = engine.CalculateUserFacingScore(userProfile.MostRecentScore)
		dashboardData.ScoreDetails = userProfile.MostRecentScore
		dashboardData.JoinedOn = userProfile.AccountCreated
		context.Debugf("Calculated user score of [%d] from internal score of [%d]", dashboardData.Score, userProfile.MostRecentScore.Value)
	}

	enc := json.NewEncoder(writer)
	enc.Encode(dashboardData)
}

func glukitScores(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	user := user.Current(context)

	glukitScoresForEmail(writer, request, user.Email)
}

func glukitScoresForDemo(writer http.ResponseWriter, request *http.Request) {
	glukitScoresForEmail(writer, request, DEMO_EMAIL)
}

// glukitScoresForEmail is the endpoint to retrieve a list of glukitscores.
func glukitScoresForEmail(writer http.ResponseWriter, request *http.Request, email string) {
	context := appengine.NewContext(request)

	limit := request.FormValue(QUERY_PARAM_LIMIT)
	fromTimestamp := request.FormValue(QUERY_PARAM_FROM)
	toTimestamp := request.FormValue(QUERY_PARAM_TO)

	// The request must include at least one of from/to/limit to be considered valid
	// as leaving it too open would could open the door for costly queries
	if len(limit) == 0 && len(fromTimestamp) == 0 && len(toTimestamp) == 0 {
		http.Error(writer, fmt.Sprintf("Query must specify at least one of: %s, %s or %s.",
			QUERY_PARAM_LIMIT, QUERY_PARAM_FROM, QUERY_PARAM_TO), 400)
		return
	}

	var scanQuery store.GlukitScoreScanQuery
	if len(limit) > 0 {
		if limitValue, err := strconv.ParseInt(limit, 10, 32); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid value for %s: [%v].", QUERY_PARAM_LIMIT, err), 400)
			return
		} else {
			limitVal := int(limitValue)
			scanQuery.Limit = &limitVal
		}
	}

	if len(fromTimestamp) > 0 {
		if fromValue, err := strconv.ParseInt(fromTimestamp, 10, 64); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid value for %s: [%v].", QUERY_PARAM_FROM, err), 400)
			return
		} else {
			fromTime := time.Unix(fromValue, 0)
			scanQuery.From = &fromTime
		}
	}

	if len(toTimestamp) > 0 {
		if toValue, err := strconv.ParseInt(toTimestamp, 10, 64); err != nil {
			http.Error(writer, fmt.Sprintf("Invalid value for %s: [%v].", QUERY_PARAM_TO, err), 400)
			return
		} else {
			toTime := time.Unix(toValue, 0)
			scanQuery.To = &toTime
		}
	}

	glukitScores, err := store.GetGlukitScores(context, email, scanQuery)
	if err != nil {
		util.Propagate(err)
	}

	if len(glukitScores) < 1 {
		http.Error(writer, "No glukit scores calculated yet.", 204)
		return
	}

	value := writer.Header()
	value.Add("Content-type", "application/json")

	enc := json.NewEncoder(writer)
	enc.Encode(glukitScores)

}
