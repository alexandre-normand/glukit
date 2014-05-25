package model

import (
	"github.com/alexandre-normand/glukit/app/util"
)

const (
	GLUCOSE_READ_TAG = "GlucoseRead"

	// Units
	MMOL_PER_L = "mmolPerL"
	MG_PER_DL  = "mgPerDL"
)

// GlucoseRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type GlucoseRead struct {
	Time  Time    `json:"time"`
	Unit  string  `json:"unit"`
	Value float32 `json:"value"`
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
	return slice[i].Time.Timestamp < slice[j].Time.EpochTime
}

func (slice GlucoseReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice GlucoseReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice GlucoseReadSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Time.Timestamp / 1000
}

// ToDataPointSlice converts a GlucoseReadSlice into a generic DataPoint array
func (slice GlucoseReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Time.LocalTime, slice.GetEpochTime(i), slice[i].Value, slice[i].Value, GLUCOSE_READ_TAG}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_GLUCOSE_READ = GlucoseRead{Timestamp{"2004-01-01 00:00:00 UTC", util.GLUKIT_EPOCH_TIME.Unix()}, UNDEFINED_READ}
