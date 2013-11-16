package model

import (
	"time"
)

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

// Represents the structure of the dashboard data for a user
type DashboardData struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Picture      string      `json:"picture"`
	LastSync     time.Time   `json:"lastSync"`
	Average      float64     `json:"average"`
	Median       float64     `json:"median"`
	High         float64     `json:"high"`
	Low          float64     `json:"low"`
	Score        *int64      `json:"score"`
	ScoreDetails GlukitScore `json:"scoreDetails"`
	JoinedOn     time.Time   `json:"joinedOn"`
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
