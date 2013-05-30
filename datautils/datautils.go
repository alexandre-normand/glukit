package datautils

import (
	"models"
	"time"
	"sort"
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

func FilterReads(data []models.MeterRead, lowerBound, upperBound time.Time) (filteredData []models.MeterRead) {
	// Nothing to sort, return immediately
	if (len(data) == 0) {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(models.MeterReadSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].TimeValue), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex:endOfDayIndex + 1]

}
