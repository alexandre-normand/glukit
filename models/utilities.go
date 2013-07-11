package models


// LinearInterpolateY does a linear interpolation of the Y value of a given GlucoseRead for a given
// time value
func LinearInterpolateY(reads []GlucoseRead, timeValue TimeValue) (yValue int) {
	lowerIndex := -1
	upperIndex := -1
	for i := range reads {
		if reads[i].TimeValue > timeValue {
			lowerIndex = i - 1;
			upperIndex = i
			break;
		}
	}

	lowerTimeValue := reads[lowerIndex].TimeValue
	upperTimeValue := reads[upperIndex].TimeValue
	lowerYValue := reads[lowerIndex].Value
	upperYValue := reads[upperIndex].Value

	relativeTimePosition := float32((timeValue - lowerTimeValue))/float32((upperTimeValue - lowerTimeValue))
	yValue = int(relativeTimePosition*float32(upperYValue - lowerYValue) + float32(lowerYValue))

	return yValue
}
