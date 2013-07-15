package store

import (
	"app/model"
	"sort"
	"time"
)

// FilterReads filters out any GlucoseRead that is outside of the lower and upper bounds. The two bounds are inclusive.
func FilterReads(data []model.GlucoseRead, lowerBound, upperBound time.Time) (filteredData []model.GlucoseRead) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.GlucoseReadSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	// We don't have data in those boundaries so return an empty array
	if startOfDayIndex == -1 || endOfDayIndex == -1 {
		return make([]model.GlucoseRead, 0)
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// FilterReads filters out any Injection that is outside of the lower and upper bounds. The two bounds are inclusive.
func FilterInjections(data []model.Injection, lowerBound, upperBound time.Time) (filteredData []model.Injection) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.InjectionSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	// We don't have data in those boundaries so return an empty array
	if startOfDayIndex == -1 || endOfDayIndex == -1 {
		return make([]model.Injection, 0)
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// FilterReads filters out any Carb element that is outside of the lower and upper bounds. The two bounds are inclusive.
func FilterCarbs(data []model.Carb, lowerBound, upperBound time.Time) (filteredData []model.Carb) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.CarbSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	// We don't have data in those boundaries so return an empty array
	if startOfDayIndex == -1 || endOfDayIndex == -1 {
		return make([]model.Carb, 0)
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// FilterReads filters out any Exercise element that is outside of the lower and upper bounds. The two bounds are inclusive.
func FilterExercises(data []model.Exercise, lowerBound, upperBound time.Time) (filteredData []model.Exercise) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.ExerciseSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	// We don't have data in those boundaries so return an empty array
	if startOfDayIndex == -1 || endOfDayIndex == -1 {
		return make([]model.Exercise, 0)
	}

	return data[startOfDayIndex : endOfDayIndex+1]
}

// MergeGlucoseReadArrays merges two arrays of GlucoseRead elements.
func MergeGlucoseReadArrays(first, second []model.GlucoseRead) []model.GlucoseRead {
	newslice := make([]model.GlucoseRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// MergeCarbArrays merges two arrays of Carb elements.
func MergeCarbArrays(first, second []model.Carb) []model.Carb {
	newslice := make([]model.Carb, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// MergeInjectionArrays merges two arrays of Injection elements.
func MergeInjectionArrays(first, second []model.Injection) []model.Injection {
	newslice := make([]model.Injection, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
