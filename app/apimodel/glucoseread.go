package apimodel

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	GLUCOSE_READ_TAG = "GlucoseRead"

	// Units
	MMOL_PER_L                       = "mmolPerL"
	MG_PER_DL                        = "mgPerDL"
	UNKNOWN_GLUCOSE_MEASUREMENT_UNIT = "Unknown"
)

type GlucoseUnit string

// GlucoseRead represents a CGM read (not to be confused with a MeterRead which is a calibration value from an external
// meter
type GlucoseRead struct {
	Time  Time        `json:"time" datastore:"time,noindex"`
	Unit  GlucoseUnit `json:"unit" datastore:"unit,noindex"`
	Value float32     `json:"value" datastore:"value,noindex"`
}

// This holds an array of reads for a whole day
type DayOfGlucoseReads struct {
	Reads     []GlucoseRead `datastore:"reads,noindex"`
	StartTime time.Time     `datastore:"startTime"`
	EndTime   time.Time     `datastore:"endTime"`
}

func NewDayOfGlucoseReads(reads []GlucoseRead) DayOfGlucoseReads {
	return DayOfGlucoseReads{reads, reads[0].GetTime().Truncate(DAY_OF_DATA_DURATION), reads[len(reads)-1].GetTime()}
}

// GetTime gets the time of a Timestamp value
func (element GlucoseRead) GetTime() time.Time {
	return element.Time.GetTime()
}

// GetNormalizedValue gets the normalized value to the requested unit
func (element GlucoseRead) GetNormalizedValue(unit GlucoseUnit) (float32, error) {
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

type GlucoseReadSlice []GlucoseRead

func (slice GlucoseReadSlice) Len() int {
	return len(slice)
}

func (slice GlucoseReadSlice) Less(i, j int) bool {
	return slice[i].Time.Timestamp < slice[j].Time.Timestamp
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
func (slice GlucoseReadSlice) ToDataPointSlice(glucoseUnit GlucoseUnit) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		// It's pretty terrible if this happens and we crash the app but this is a coding error and I want to know early
		convertedValue, err := slice[i].GetNormalizedValue(glucoseUnit)
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i), convertedValue, convertedValue, GLUCOSE_READ_TAG, glucoseUnit}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

var UNDEFINED_GLUCOSE_READ = GlucoseRead{Time{GetTimeMillis(util.GLUKIT_EPOCH_TIME), "UTC"}, "NONE", UNDEFINED_READ}
