package models

import (
	"time"
	"timeutils"
)

// Injection represents an insulin injection
type Injection struct {
	LocalTime          string    `json:"label" datastore:"localtime,noindex"`
	TimeValue          TimeValue `json:"x" datastore:"timestamp"`
	Units              float32   `json:"units" datastore:"units,noindex"`
	ReferenceReadValue int       `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfInjections struct {
	Injections []Injection
}

// GetTime gets the time of an Injection value
func (injection Injection) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(injection.LocalTime, timeutils.TIMEZONE)
	return value
}

type InjectionSlice []Injection

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an InjectionSlice into a generic DataPoint array
func (slice InjectionSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue,
			LinearInterpolateY(matchingReads, slice[i].TimeValue), slice[i].Units}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
