package apimodel

// Represents a generic data point in time
type DataPoint struct {
	LocalTime string      `json:"label"`
	EpochTime int64       `json:"x"`
	Y         float32     `json:"y"`
	Value     float32     `json:"value"`
	Tag       string      `json:"tag"`
	Unit      GlucoseUnit `json:"unit"`
}

type DataPointSlice []DataPoint

func (slice DataPointSlice) Len() int {
	return len(slice)
}

func (slice DataPointSlice) Less(i, j int) bool {
	return slice[i].EpochTime < slice[j].EpochTime
}

func (slice DataPointSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
