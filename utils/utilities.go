package utils

import (
	"models"
)

func MergeReadArrays(first, second []models.GlucoseRead) []models.GlucoseRead {
   newslice := make([]models.GlucoseRead, len(first) + len(second))
   copy(newslice, first)
   copy(newslice[len(first):], second)
   return newslice
}

func MergeCarbArrays(first, second []models.Carb) []models.Carb {
   newslice := make([]models.Carb, len(first) + len(second))
   copy(newslice, first)
   copy(newslice[len(first):], second)
   return newslice
}

func MergeInjectionArrays(first, second []models.Injection) []models.Injection {
   newslice := make([]models.Injection, len(first) + len(second))
   copy(newslice, first)
   copy(newslice[len(first):], second)
   return newslice
}
