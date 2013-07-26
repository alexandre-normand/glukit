package model

import (
	"app/util"
	"time"
)

// Represents an exercise event
type Exercise struct {
	LocalTime         string    `json:"label" datastore:"localtime,noindex"`
	Timestamp         Timestamp `json:"unixtime" datastore:"timestamp"`
	DurationInMinutes int       `json:"duration" datastore:"duration,noindex"`
	// One of: light, medium, heavy
	Intensity string `json:"intensity" datastore:"intensity,noindex"`
}

// Gets the time of the exercise
func (exercise Exercise) GetTime() (timeValue time.Time) {
	value, _ := util.ParseTime(exercise.LocalTime, util.TIMEZONE)
	return value
}

// This holds an array of exercise events for a whole day
type DayOfExercises struct {
	Exercises []Exercise
}

type ExerciseSlice []Exercise

func (slice ExerciseSlice) Len() int {
	return len(slice)
}

func (slice ExerciseSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice ExerciseSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ToDataPointSlice converts an ExerciseSlice into a generic DataPoint array
func (slice ExerciseSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].Timestamp,
			linearInterpolateY(matchingReads, slice[i].Timestamp), float32(slice[i].DurationInMinutes)}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
