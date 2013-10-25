package model

// linearInterpolateY does a linear interpolation of the Y value of a given GlucoseRead for a given
// time value
func linearInterpolateY(reads []GlucoseRead, timeValue Timestamp) (yValue int) {
	lowerIndex := 0
	upperIndex := len(reads) - 1

	for i := range reads {
		if reads[i].Timestamp > timeValue {
			upperIndex = i
			break
		}
	}
	// Handle the case where the timestamp is before the first read we have.
	// In such as case, we don't interpolate and return the Y value of that read
	if upperIndex == 0 {
		return reads[upperIndex].Value
	} else {
		lowerIndex = upperIndex - 1
	}

	lowerTimeValue := reads[lowerIndex].Timestamp
	upperTimeValue := reads[upperIndex].Timestamp
	lowerYValue := reads[lowerIndex].Value
	upperYValue := reads[upperIndex].Value

	relativeTimePosition := float32((timeValue - lowerTimeValue)) / float32((upperTimeValue - lowerTimeValue))
	yValue = int(relativeTimePosition*float32(upperYValue-lowerYValue) + float32(lowerYValue))

	return yValue
}

func MergeDataPointArrays(first, second []DataPoint) []DataPoint {
	newslice := make([]DataPoint, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
