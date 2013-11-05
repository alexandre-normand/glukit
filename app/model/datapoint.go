package model

// Represents a generic data point in time
type DataPoint struct {
	LocalTime string    `json:"label"`
	Timestamp Timestamp `json:"x"`
	Y         int       `json:"y"`
	Value     float32   `json:"value"`
	Tag       string    `json:"tag"`
}

type DataPointSlice []DataPoint

func (slice DataPointSlice) Len() int {
	return len(slice)
}

func (slice DataPointSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice DataPointSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
