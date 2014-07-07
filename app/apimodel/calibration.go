package apimodel

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	CALIBRATION_READ_TAG = "CalibrationRead"
)

// CalibrationRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type CalibrationRead struct {
	Time  Time        `json:"time" datastore:"time,noindex"`
	Unit  GlucoseUnit `json:"unit" datastore:"unit,noindex"`
	Value float32     `json:"value" datastore:"value,noindex"`
}

// This holds an array of reads for a whole day
type DayOfCalibrationReads struct {
	Reads     []CalibrationRead `datastore:"calibrations,noindex"`
	StartTime time.Time         `datastore:"startTime"`
	EndTime   time.Time         `datastore:"endTime"`
}

// GetTime gets the time of a Timestamp value
func (element CalibrationRead) GetTime() time.Time {
	return element.Time.GetTime()
}

func NewDayOfCalibrationReads(reads []CalibrationRead) DayOfCalibrationReads {
	return DayOfCalibrationReads{reads, reads[0].GetTime().Truncate(DAY_OF_DATA_DURATION), reads[len(reads)-1].GetTime()}
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

// GetNormalizedValue gets the normalized value to the requested unit
func (element CalibrationRead) GetNormalizedValue(unit GlucoseUnit) (float32, error) {
	if unit == element.Unit {
		return element.Value, nil
	}

	if element.Unit == UNKNOWN_GLUCOSE_MEASUREMENT_UNIT {
		return element.Value, nil
	}

	// This switch can focus on only conversion cases because the obvious
	// cases have been sorted out already
	switch unit {
	case MMOL_PER_L:
		return element.Value * 0.0555, nil
	case MG_PER_DL:
		return element.Value * 18.0182, nil
	default:
		return -1., errors.New(fmt.Sprintf("Bad unit requested, [%s] is not one of [%s, %s]", unit, MG_PER_DL, MMOL_PER_L))
	}
}

// ToDataPointSlice converts a CalibrationReadSlice into a generic DataPoint array
func (slice CalibrationReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		// It's pretty terrible if this happens and we crash the app but this is a coding error and I want to know early
		mgPerDlValue, err := slice[i].GetNormalizedValue(MG_PER_DL)
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i), mgPerDlValue, float32(slice[i].Value), CALIBRATION_READ_TAG, MG_PER_DL}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_CALIBRATION_READ = CalibrationRead{Time{0, "UTC"}, "NONE", -1.}
