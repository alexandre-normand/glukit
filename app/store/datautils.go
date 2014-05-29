package store

import (
	"github.com/alexandre-normand/glukit/app/model"
)

// mergeGlucoseReadArrays merges two arrays of GlucoseRead elements.
func mergeGlucoseReadArrays(first, second []model.GlucoseRead) []model.GlucoseRead {
	newslice := make([]model.GlucoseRead, len(first)+len(second))
	copy(newslice, first)
	copy(newslice[len(first):], second)
	return newslice
}

// mergeMealArrays merges two arrays of Meal elements.
func mergeMealArrays(first, second []model.Meal) []model.Meal {
	newslice := make([]model.Meal, len(first)+len(second))
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
