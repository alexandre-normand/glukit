// This package provides some useful fonctions to filter/merge and do various operations on the data.
package datautils

import (
	"app/models"
	"sort"
	"time"
)

// FilterReads filters out any GlucoseRead that is outside of the lower and upper bounds. The two bounds are inclusive.
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
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]

}

// FilterReads filters out any Injection that is outside of the lower and upper bounds. The two bounds are inclusive.
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
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// FilterReads filters out any Carb element that is outside of the lower and upper bounds. The two bounds are inclusive.
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
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// FilterReads filters out any Exercise element that is outside of the lower and upper bounds. The two bounds are inclusive.
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
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// MergeGlucoseReadArrays merges two arrays of GlucoseRead elements.
func MergeGlucoseReadArrays(first, second []models.GlucoseRead) []models.GlucoseRead {
	newslice := make([]models.GlucoseRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// MergeCarbArrays merges two arrays of Carb elements.
func MergeCarbArrays(first, second []models.Carb) []models.Carb {
	newslice := make([]models.Carb, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// MergeInjectionArrays merges two arrays of Injection elements.
func MergeInjectionArrays(first, second []models.Injection) []models.Injection {
	newslice := make([]models.Injection, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
