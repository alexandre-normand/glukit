package model

const (
	CARB_TAG = "Carbs"
)

// Represents a carbohydrate intake
type Carb struct {
	Timestamp
	Grams              float32 `json:"carbs" datastore:"grams,noindex"`
	ReferenceReadValue int     `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfCarbs struct {
	Carbs []Carb
}

type CarbSlice []Carb

func (slice CarbSlice) Len() int {
	return len(slice)
}

func (slice CarbSlice) Less(i, j int) bool {
	return slice[i].Timestamp.EpochTime < slice[j].Timestamp.EpochTime
}

func (slice CarbSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CarbSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Timestamp.EpochTime
}

// ToDataPointSlice converts an CarbSlice into a generic DataPoint array
func (slice CarbSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime,
			linearInterpolateY(matchingReads, slice[i].Timestamp), slice[i].Grams, CARB_TAG}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
