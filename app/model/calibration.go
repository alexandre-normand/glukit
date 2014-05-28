package model

import (
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	CALIBRATION_READ_TAG = "CalibrationRead"
)

// CalibrationRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type CalibrationRead struct {
	Time  Time    `json:"time"`
	Unit  string  `json:"unit"`
	Value float32 `json:"value"`
}

// This holds an array of reads for a whole day
type DayOfCalibrationReads struct {
	Reads []CalibrationRead
}

// GetTime gets the time of a Timestamp value
func (element CalibrationRead) GetTime() time.Time {
	return element.Time.GetTime()
}

type CalibrationReadSlice []CalibrationRead

func (slice CalibrationReadSlice) Len() int {
	return len(slice)
}

func (slice CalibrationReadSlice) Less(i, j int) bool {
	return slice[i].Time.Timestamp < slice[j].Time.Timestamp
}

func (slice CalibrationReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CalibrationReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice CalibrationReadSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Time.Timestamp / 1000
}

// ToDataPointSlice converts a CalibrationReadSlice into a generic DataPoint array
func (slice CalibrationReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i), slice[i].Value, float32(slice[i].Value), CALIBRATION_READ_TAG}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_CALIBRATION_READ = CalibrationRead{Time{0, "UTC"}, "NONE", -1.}
