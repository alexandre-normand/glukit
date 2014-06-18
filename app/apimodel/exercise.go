package apimodel

import (
	"github.com/alexandre-normand/glukit/app/util"
	"time"
)

const (
	EXERCISE_TAG = "Exercise"
)

type Exercise struct {
	Time            Time   `json:"time"`
	DurationMinutes int    `json:"durationInMinutes"`
	Intensity       string `json:"intensity"`
	Description     string `json:"description"`
}

// This holds an array of exercise events for a whole day
type DayOfExercises struct {
	Exercises []Exercise `datastore:"exercises"`
	StartTime time.Time  `datastore:"startTime"`
	EndTime   time.Time  `datastore:"endTime"`
}

func NewDayOfExercises(exercises []Exercise) DayOfExercises {
	return DayOfExercises{exercises, exercises[0].GetTime(), exercises[len(exercises)-1].GetTime()}
}

// GetTime gets the time of a Timestamp value
func (element Exercise) GetTime() time.Time {
	return element.Time.GetTime()
}

type ExerciseSlice []Exercise

func (slice ExerciseSlice) Len() int {
	return len(slice)
}

func (slice ExerciseSlice) Less(i, j int) bool {
	return slice[i].Time.Timestamp < slice[j].Time.Timestamp
}

func (slice ExerciseSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice ExerciseSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Time.Timestamp
}

// ToDataPointSlice converts an ExerciseSlice into a generic DataPoint array
func (slice ExerciseSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		localTime, err := slice[i].Time.Format()
		if err != nil {
			util.Propagate(err)
		}

		dataPoint := DataPoint{localTime, slice.GetEpochTime(i),
			linearInterpolateY(matchingReads, slice[i].Time), float32(slice[i].DurationMinutes), EXERCISE_TAG, "minutes"}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
