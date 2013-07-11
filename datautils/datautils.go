package datautils

import (
	"models"
	"sort"
	"time"
)

func BuildHistogram(reads []models.GlucoseRead) (data []models.Coordinate) {
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

	data = make([]models.Coordinate, len(distribution))
	j := 0
	for key, value := range distribution {
		data[j] = models.Coordinate{key, value}
		j = j + 1
	}

	return data
}

func FilterReads(data []models.GlucoseRead, lowerBound, upperBound time.Time) (filteredData []models.GlucoseRead) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(models.GlucoseReadSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].TimeValue), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]

}

func FilterInjections(data []models.Injection, lowerBound, upperBound time.Time) (filteredData []models.Injection) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(models.InjectionSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].TimeValue), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

func FilterCarbs(data []models.Carb, lowerBound, upperBound time.Time) (filteredData []models.Carb) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(models.CarbSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].TimeValue), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

func FilterExercises(data []models.Exercise, lowerBound, upperBound time.Time) (filteredData []models.Exercise) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(models.ExerciseSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].TimeValue), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}
