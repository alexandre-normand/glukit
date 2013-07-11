package models

import (
	"timeutils"
	"time"
)

// Represents an exercise event
type Exercise struct {
	LocalTime         string      `json:"label" datastore:"localtime,noindex"`
	TimeValue         TimeValue   `json:"unixtime" datastore:"timestamp"`
	DurationInMinutes int         `json:"duration" datastore:"duration,noindex"`
	// One of: light, medium, heavy
	Intensity         string      `json:"intensity" datastore:"intensity,noindex"`
}

// Gets the time of the exercise
func (exercise Exercise) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(exercise.LocalTime, timeutils.TIMEZONE)
	return value
}

// This holds an array of exercise events for a whole day
type DayOfExercises struct {
	Exercises      []Exercise
}

type ExerciseSlice []Exercise

func (slice ExerciseSlice) Len() int {
	return len(slice)
}

func (slice ExerciseSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice ExerciseSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an ExerciseSlice into a generic DataPoint array
func (slice ExerciseSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue,
			LinearInterpolateY(matchingReads, slice[i].TimeValue), float32(slice[i].DurationInMinutes)}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
