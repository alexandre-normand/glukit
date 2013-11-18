package model

const (
	INSULIN_TAG = "Insulin"
)

// Injection represents an insulin injection
type Injection struct {
	Timestamp
	Units              float32 `json:"units" datastore:"units,noindex"`
	ReferenceReadValue int     `json:"y" datastore:"referenceReadValue,noindex"`
}

// This holds an array of injections for a whole day
type DayOfInjections struct {
	Injections []Injection
}

type InjectionSlice []Injection

// ToTimestampSlice converts an InjectionSlice into a generic TimestampSlice
func (slice InjectionSlice) ToTimestampSlice() (timestamps TimestampSlice) {
	timestamps = make([]Timestamp, len(slice))

	for i := range slice {
		timestamp := slice[i].Timestamp
		timestamps[i] = timestamp
	}

	return timestamps
}

// ToDataPointSlice converts an InjectionSlice into a generic DataPoint array
func (slice InjectionSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))

	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime,
			linearInterpolateY(matchingReads, slice[i].Timestamp), slice[i].Units, INSULIN_TAG}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
