package store

import (
	"app/model"
	"sort"
	"time"
)

// filterReads filters out any GlucoseRead that is outside of the lower and upper bounds. The two bounds are inclusive.
func filterReads(data []model.GlucoseRead, lowerBound, upperBound time.Time) (filteredData []model.GlucoseRead) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := 0
	endOfDayIndex := arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.GlucoseReadSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.After(upperBound) {
			endOfDayIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.Before(lowerBound) {
			startOfDayIndex = i
			break
		}
	}

	return data[startOfDayIndex:endOfDayIndex]
}

// filterInjections filters out any Injection that is outside of the lower and upper bounds. The two bounds are inclusive.
func filterInjections(data []model.Injection, lowerBound, upperBound time.Time) (filteredData []model.Injection) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := 0
	endOfDayIndex := arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.InjectionSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.After(upperBound) {
			endOfDayIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.Before(lowerBound) {
			startOfDayIndex = i
			break
		}
	}

	return data[startOfDayIndex:endOfDayIndex]
}

// filterCarbs filters out any Carb element that is outside of the lower and upper bounds. The two bounds are inclusive.
func filterCarbs(data []model.Carb, lowerBound, upperBound time.Time) (filteredData []model.Carb) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := 0
	endOfDayIndex := arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.CarbSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.After(upperBound) {
			endOfDayIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.Before(lowerBound) {
			startOfDayIndex = i
			break
		}
	}

	return data[startOfDayIndex:endOfDayIndex]
}

// filterExercises filters out any Exercise element that is outside of the lower and upper bounds. The two bounds are inclusive.
func filterExercises(data []model.Exercise, lowerBound, upperBound time.Time) (filteredData []model.Exercise) {
	// Nothing to sort, return immediately
	if len(data) == 0 {
		return data
	}

	arraySize := len(data)
	startOfDayIndex := 0
	endOfDayIndex := arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(model.ExerciseSlice(data))

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.After(upperBound) {
			endOfDayIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(int64(data[i].Timestamp), 0)
		if !elementTime.Before(lowerBound) {
			startOfDayIndex = i
			break
		}
	}

	return data[startOfDayIndex:endOfDayIndex]
}

// mergeGlucoseReadArrays merges two arrays of GlucoseRead elements.
func mergeGlucoseReadArrays(first, second []model.GlucoseRead) []model.GlucoseRead {
	newslice := make([]model.GlucoseRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeCarbArrays merges two arrays of Carb elements.
func mergeCarbArrays(first, second []model.Carb) []model.Carb {
	newslice := make([]model.Carb, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeInjectionArrays merges two arrays of Injection elements.
func mergeInjectionArrays(first, second []model.Injection) []model.Injection {
	newslice := make([]model.Injection, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
