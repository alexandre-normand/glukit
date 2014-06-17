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
