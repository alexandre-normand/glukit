package model

import (
	"time"
)

// Represents a carbohydrate intake
type Carb struct {
	LocalTime          string    `json:"label" datastore:"localtime,noindex"`
	Timestamp          Timestamp `json:"x" datastore:"timestamp"`
	Grams              float32   `json:"carbs" datastore:"grams,noindex"`
	ReferenceReadValue int       `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfCarbs struct {
	Carbs []Carb
}

// GetTime gets the time of a GlucoseRead value
func (carb Carb) GetTime() (timeValue time.Time) {
	value := time.Unix(int64(carb.Timestamp), 0)
	return value
}

type CarbSlice []Carb

func (slice CarbSlice) Len() int {
	return len(slice)
}

func (slice CarbSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice CarbSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an CarbSlice into a generic DataPoint array
func (slice CarbSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].Timestamp,
			linearInterpolateY(matchingReads, slice[i].Timestamp), slice[i].Grams}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
