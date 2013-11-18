package store

import (
	"app/model"
)

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
