package model

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
	Exercises []Exercise
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
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime,
			linearInterpolateY(matchingReads, slice[i].Timestamp), float32(slice[i].DurationInMinutes), EXERCISE_TAG}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
