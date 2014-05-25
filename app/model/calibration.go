package model

import (
	"github.com/alexandre-normand/glukit/app/util"
)

const (
	CALIBRATION_READ_TAG = "CalibrationRead"
)

// CalibrationRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type CalibrationRead struct {
	Time
	Value int `json:"y" datastore:"value,noindex"`
}

// This holds an array of reads for a whole day
type DayOfCalibrationReads struct {
	Reads []CalibrationRead
}

// func (slice CalibrationReadSlice) Len() int {
// 	return len(slice)
// }
type CalibrationReadSlice []CalibrationRead

func (slice CalibrationReadSlice) Len() int {
	return len(slice)
}

func (slice CalibrationReadSlice) Less(i, j int) bool {
	return slice[i].Timestamp.EpochTime < slice[j].Timestamp.EpochTime
}

func (slice CalibrationReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CalibrationReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice CalibrationReadSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Timestamp.EpochTime
}

// ToDataPointSlice converts a CalibrationReadSlice into a generic DataPoint array
func (slice CalibrationReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime, slice[i].Value, float32(slice[i].Value), CALIBRATION_READ_TAG}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_Calibration_READ = CalibrationRead{Timestamp{"2004-01-01 00:00:00 UTC", util.GLUKIT_EPOCH_TIME.Unix()}, UNDEFINED_READ}
