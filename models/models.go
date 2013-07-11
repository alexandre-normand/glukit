/*
Package models provides model types
*/
package models

import (
	"time"
	"timeutils"
)

const (
	UNDEFINED_READ        = -1
	EXERCISE_VALUE_FORMAT = "%d,%s"
	UNDEFINED_SCORE_VALUE = int64(-1)
)

// "Dynamic" constants, those should never be updated
var UNDEFINED_SCORE = GlukitScore{UNDEFINED_SCORE_VALUE, timeutils.BEGINNING_OF_TIME, timeutils.BEGINNING_OF_TIME}

type TimeValue int64

type TrackingData struct {
	Mean         float64      `json:"mean"`
	Median       float64      `json:"median"`
	Deviation    float64      `json:"deviation"`
	Min          float64      `json:"min"`
	Max          float64      `json:"max"`
	Distribution []Coordinate `json:"distribution"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type DataPoint struct {
	LocalTime string    `json:"label"`
	TimeValue TimeValue `json:"x"`
	Y         int       `json:"y"`
	Value     float32   `json:"r"`
}

type FileImportLog struct {
	Id                string
	Md5Checksum       string
	LastDataProcessed time.Time
}

type PointData interface {
	GetTime() time.Time
}

type ReadStatsSlice []GlucoseRead
type CoordinateSlice []Coordinate

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
