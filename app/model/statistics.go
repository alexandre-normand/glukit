package model

// Type for a slice of GlucoseReads with comparison based on value rather than time. It is used as read statistics.
type ReadStatsSlice []GlucoseRead

func (slice ReadStatsSlice) Len() int {
	return len(slice)
}

func (slice ReadStatsSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice ReadStatsSlice) Less(i, j int) bool {
	return slice[i].Value < slice[j].Value
}

func (slice ReadStatsSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Represents the structure of the tracking data for a user
type TrackingData struct {
	Mean         float64      `json:"mean"`
	Median       float64      `json:"median"`
	Deviation    float64      `json:"deviation"`
	Min          float64      `json:"min"`
	Max          float64      `json:"max"`
	Distribution []Coordinate `json:"distribution"`
}

type CoordinateSlice []Coordinate

func (slice CoordinateSlice) Len() int {
	return len(slice)
}

func (slice CoordinateSlice) Get(i int) int {
	return slice[i].X
}

func (slice CoordinateSlice) Less(i, j int) bool {
	return slice[i].X < slice[j].X
}

func (slice CoordinateSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
