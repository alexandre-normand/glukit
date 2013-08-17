package model

import (
	"app/util"
	"time"
)

// GlucoseRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type GlucoseRead struct {
	LocalTime string    `json:"label" datastore:"localtime,noindex"`
	Timestamp Timestamp `json:"x" datastore:"timestamp"`
	Value     int       `json:"y" datastore:"value,noindex"`
}

// This holds an array of reads for a whole day
type DayOfGlucoseReads struct {
	Reads []GlucoseRead
}

// GetTime gets the time of a GlucoseRead value
func (read GlucoseRead) GetTime() (timeValue time.Time) {
	value := time.Unix(int64(read.Timestamp), 0)
	return value
}

type GlucoseReadSlice []GlucoseRead

func (slice GlucoseReadSlice) Len() int {
	return len(slice)
}

func (slice GlucoseReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice GlucoseReadSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice GlucoseReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts a GlucoseReadSlice into a generic DataPoint array
func (slice GlucoseReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].Timestamp, slice[i].Value, float32(slice[i].Value)}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_GLUCOSE_READ = GlucoseRead{"2004-01-01 00:00:00 UTC", Timestamp(util.GLUKIT_EPOCH_TIME.Unix()), UNDEFINED_READ}
