package datautils

import (
	"models"
)

func BuildHistogram(reads []models.MeterRead) (data []models.Coordinate) {
	distribution := make(map[int] int)

	for i := range reads {
		currentReadValue := reads[i].Value
		currentValue, found := distribution[currentReadValue]
		if found {
			distribution[currentReadValue] = currentValue + 1
		} else {
			distribution[currentReadValue] = 1
		}
	}

	data = make([]models.Coordinate, len(distribution))
    j := 0
	for key, value := range distribution  {
		data[j] = models.Coordinate{key, value}
		j = j + 1
	}

	return data
}
