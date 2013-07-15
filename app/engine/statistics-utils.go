package engine

import (
	"app/model"
)

// BuildHistogram generates a histogram from an array of GlucoseRead. The resulting histogram is
// an array of Coordinates where the X value is the value of the read and the Y value is the number
// of instances of that value
func BuildHistogram(reads []model.GlucoseRead) (data []model.Coordinate) {
	distribution := make(map[int]int)

	for i := range reads {
		currentReadValue := reads[i].Value
		currentValue, found := distribution[currentReadValue]
		if found {
			distribution[currentReadValue] = currentValue + 1
		} else {
			distribution[currentReadValue] = 1
		}
	}

	data = make([]model.Coordinate, len(distribution))
	j := 0
	for key, value := range distribution {
		data[j] = model.Coordinate{key, value}
		j = j + 1
	}

	return data
}
