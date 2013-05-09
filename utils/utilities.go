package utils

import (
	"models"
)

func MergeReadArrays(first, second []models.MeterRead) []models.MeterRead {
   newslice := make([]models.MeterRead, len(first) + len(second))
   copy(newslice, first)
   copy(newslice[len(first):], second)
   return newslice
}

func MergeCarbIntakeArrays(first, second []models.CarbIntake) []models.CarbIntake {
   newslice := make([]models.CarbIntake, len(first) + len(second))
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
