package apimodel

import (
	"github.com/alexandre-normand/glukit/app/util"
)

// linearInterpolateY does a linear interpolation of the Y value of a given GlucoseRead for a given
// time value
func linearInterpolateY(reads []GlucoseRead, timeValue Time) (yValue float32) {
	lowerIndex := 0
	upperIndex := len(reads) - 1

	for i := range reads {
		if reads[i].Time.Timestamp > timeValue.Timestamp {
			upperIndex = i
			break
		}
	}
	// Handle the case where the timestamp is before the first read we have.
	// In such as case, we don't interpolate and return the Y value of that read
	if upperIndex == 0 {
		value, err := reads[upperIndex].GetNormalizedValue(MG_PER_DL)
		if err != nil {
			util.Propagate(err)
		}
		return value
	} else {
		lowerIndex = upperIndex - 1
	}

	lowerTimeValue := reads[lowerIndex].Time
	upperTimeValue := reads[upperIndex].Time
	lowerYValue, err := reads[lowerIndex].GetNormalizedValue(MG_PER_DL)
	if err != nil {
		util.Propagate(err)
	}
	upperYValue, err := reads[upperIndex].GetNormalizedValue(MG_PER_DL)
	if err != nil {
		util.Propagate(err)
	}

	relativeTimePosition := float32((timeValue.Timestamp - lowerTimeValue.Timestamp)) / float32((upperTimeValue.Timestamp - lowerTimeValue.Timestamp))
	yValue = relativeTimePosition*float32(upperYValue-lowerYValue) + float32(lowerYValue)

	return yValue
}

func MergeDataPointArrays(first, second []DataPoint) []DataPoint {
	newslice := make([]DataPoint, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
