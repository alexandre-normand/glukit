package models

import (
	"time"
	"timeutils"
)

// Represents a carbohydrate intake
type Carb struct {
	LocalTime          string    `json:"label" datastore:"localtime,noindex"`
	TimeValue          TimeValue `json:"x" datastore:"timestamp"`
	Grams              float32   `json:"carbs" datastore:"grams,noindex"`
	ReferenceReadValue int       `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfCarbs struct {
	Carbs []Carb
}

// GetTime gets the time of a GlucoseRead value
func (carb Carb) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(carb.LocalTime, timeutils.TIMEZONE)
	return value
}

type CarbSlice []Carb

func (slice CarbSlice) Len() int {
	return len(slice)
}

func (slice CarbSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice CarbSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an CarbSlice into a generic DataPoint array
func (slice CarbSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue,
			LinearInterpolateY(matchingReads, slice[i].TimeValue), slice[i].Grams}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
