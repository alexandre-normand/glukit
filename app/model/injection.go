package model

import (
	"app/util"
	"time"
)

// Injection represents an insulin injection
type Injection struct {
	LocalTime          string    `json:"label" datastore:"localtime,noindex"`
	Timestamp          Timestamp `json:"x" datastore:"timestamp"`
	Units              float32   `json:"units" datastore:"units,noindex"`
	ReferenceReadValue int       `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfInjections struct {
	Injections []Injection
}

// GetTime gets the time of an Injection value
func (injection Injection) GetTime() (timeValue time.Time) {
	value, _ := util.ParseTime(injection.LocalTime, util.TIMEZONE)
	return value
}

type InjectionSlice []Injection

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an InjectionSlice into a generic DataPoint array
func (slice InjectionSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].Timestamp,
			linearInterpolateY(matchingReads, slice[i].Timestamp), slice[i].Units}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
