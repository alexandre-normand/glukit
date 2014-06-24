package apimodel

import (
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	CARB_TAG = "Carbs"
)

// Meal is the data structure that represents a meal of food intake. Only carbohydrates
// are fully supported at the moment.
type Meal struct {
	Time          Time    `json:"time" datastore:"time,noindex"`
	Carbohydrates float32 `json:"carbohydrates" datastore:"carbohydrates,noindex"`
	Proteins      float32 `json:"proteins" datastore:"proteins,noindex"`
	Fat           float32 `json:"fat" datastore:"fat,noindex"`
	SaturatedFat  float32 `json:"saturatedFat" datastore:"saturatedFat,noindex"`
}

// This holds an array of injections for a whole day
type DayOfMeals struct {
	Meals     []Meal    `datastore:"meals,noindex"`
	StartTime time.Time `datastore:"startTime"`
	EndTime   time.Time `datastore:"endTime"`
}

func NewDayOfMeals(meals []Meal) DayOfMeals {
	return DayOfMeals{meals, meals[0].GetTime(), meals[len(meals)-1].GetTime()}
}

// GetTime gets the time of a Timestamp value
func (element Meal) GetTime() time.Time {
	return element.Time.GetTime()
}

type MealSlice []Meal

func (slice MealSlice) Len() int {
	return len(slice)
}

func (slice MealSlice) Less(i, j int) bool {
	return slice[i].Time.Timestamp < slice[j].Time.Timestamp
}

func (slice MealSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice MealSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Time.Timestamp / 1000
}

// ToDataPointSlice converts an MealSlice into a generic DataPoint array
func (slice MealSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i),
			linearInterpolateY(matchingReads, slice[i].Time), slice[i].Carbohydrates, CARB_TAG, "grams"}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
