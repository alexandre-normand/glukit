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

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].Timestamp.EpochTime < slice[j].Timestamp.EpochTime
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice InjectionSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].EpochTime
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
