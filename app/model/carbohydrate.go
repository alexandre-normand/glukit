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

// ToTimestampSlice converts a CarbSlice into a generic TimestampSlice array
func (slice CarbSlice) ToTimestampSlice() (timestamps TimestampSlice) {
	timestamps = make([]Timestamp, len(slice))

	for i := range slice {
		timestamp := slice[i].Timestamp
		timestamps[i] = timestamp
	}

	return timestamps
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
