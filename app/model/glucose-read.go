package model

import (
	"app/util"
)

const (
	GLUCOSE_READ_TAG = "GlucoseRead"
)

// GlucoseRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type GlucoseRead struct {
	Timestamp
	Value int `json:"y" datastore:"value,noindex"`
}

// This holds an array of reads for a whole day
type DayOfGlucoseReads struct {
	Reads []GlucoseRead
}

// func (slice GlucoseReadSlice) Len() int {
// 	return len(slice)
// }
type GlucoseReadSlice []GlucoseRead

func (slice GlucoseReadSlice) Len() int {
	return len(slice)
}

func (slice GlucoseReadSlice) Less(i, j int) bool {
	return slice[i].Timestamp.EpochTime < slice[j].Timestamp.EpochTime
}

func (slice GlucoseReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice GlucoseReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice GlucoseReadSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Timestamp.EpochTime
}

// ToDataPointSlice converts a GlucoseReadSlice into a generic DataPoint array
func (slice GlucoseReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime, slice[i].Value, float32(slice[i].Value), GLUCOSE_READ_TAG}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_GLUCOSE_READ = GlucoseRead{Timestamp{"2004-01-01 00:00:00 UTC", util.GLUKIT_EPOCH_TIME.Unix()}, UNDEFINED_READ}
