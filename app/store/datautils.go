package store

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
)

// mergeGlucoseReadArrays merges two arrays of GlucoseRead elements.
func mergeGlucoseReadArrays(first, second []apimodel.GlucoseRead) []apimodel.GlucoseRead {
	newslice := make([]apimodel.GlucoseRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeMealArrays merges two arrays of Meal elements.
func mergeMealArrays(first, second []apimodel.Meal) []apimodel.Meal {
	newslice := make([]apimodel.Meal, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeInjectionArrays merges two arrays of Injection elements.
func mergeInjectionArrays(first, second []apimodel.Injection) []apimodel.Injection {
	newslice := make([]apimodel.Injection, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeExerciseArrays merges two arrays of Exercise elements.
func mergeExerciseArrays(first, second []apimodel.Exercise) []apimodel.Exercise {
	newslice := make([]apimodel.Exercise, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeCalibrationReadArrays merges two arrays of CalibrationRead elements.
func mergeCalibrationReadArrays(first, second []apimodel.CalibrationRead) []apimodel.CalibrationRead {
	newslice := make([]apimodel.CalibrationRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}
