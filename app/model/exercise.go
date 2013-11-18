package model

const (
	EXERCISE_TAG = "Exercise"
)

// Represents an exercise event
type Exercise struct {
	Timestamp
	DurationInMinutes int `json:"duration" datastore:"duration,noindex"`
	// One of: light, medium, heavy
	Intensity string `json:"intensity" datastore:"intensity,noindex"`
}

// This holds an array of exercise events for a whole day
type DayOfExercises struct {
	Exercises []Exercise
}

type ExerciseSlice []Exercise

// ToTimestampSlice converts a ExerciseSlice into a generic TimestampSlice array
func (slice ExerciseSlice) ToTimestampSlice() (timestamps TimestampSlice) {
	timestamps = make([]Timestamp, len(slice))

	for i := range slice {
		timestamp := slice[i].Timestamp
		timestamps[i] = timestamp
	}

	return timestamps
}

// ToDataPointSlice converts an ExerciseSlice into a generic DataPoint array
func (slice ExerciseSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime,
			linearInterpolateY(matchingReads, slice[i].Timestamp), float32(slice[i].DurationInMinutes), EXERCISE_TAG}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
