package apimodel

import (
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	INSULIN_TAG = "Insulin"
)

// Injection represents an insulin injection
type Injection struct {
	Time        Time    `json:"time" datastore:"time,noindex"`
	Units       float32 `json:"units" datastore:"units,noindex"`
	InsulinName string  `json:"insulinName" datastore:"insulinName,noindex"`
	InsulinType string  `json:"insulinType" datastore:"insulinType,noindex"`
}

// This holds an array of injections for a whole day
type DayOfInjections struct {
	Injections []Injection `datastore:"injections,noindex"`
	StartTime  time.Time   `datastore:"startTime"`
	EndTime    time.Time   `datastore:"endTime"`
}

func NewDayOfInjections(injections []Injection) DayOfInjections {
	return DayOfInjections{injections, injections[0].GetTime().Truncate(DAY_OF_DATA_DURATION), injections[len(injections)-1].GetTime()}
}

// GetTime gets the time of a Timestamp value
func (element Injection) GetTime() time.Time {
	return element.Time.GetTime()
}

type InjectionSlice []Injection

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].Time.Timestamp < slice[j].Time.Timestamp
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice InjectionSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Time.Timestamp / 1000
}

// ToDataPointSlice converts an InjectionSlice into a generic DataPoint array
func (slice InjectionSlice) ToDataPointSlice(matchingReads []GlucoseRead, glucoseUnit GlucoseUnit) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))

	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i),
			linearInterpolateY(matchingReads, slice[i].Time, glucoseUnit), slice[i].Units, INSULIN_TAG, "units"}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
