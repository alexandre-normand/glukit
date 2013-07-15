/*
Package models provides model types
*/
package model

import (
	"app/util"
	"time"
)

const (
	UNDEFINED_READ        = -1
	EXERCISE_VALUE_FORMAT = "%d,%s"
	UNDEFINED_SCORE_VALUE = int64(-1)
)

// "Dynamic" constants, those should never be updated
var UNDEFINED_SCORE = GlukitScore{UNDEFINED_SCORE_VALUE, util.BEGINNING_OF_TIME, util.BEGINNING_OF_TIME}

type Timestamp int64

// Represents a cartesian coordinate
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Represents a generic data point in time
type DataPoint struct {
	LocalTime string    `json:"label"`
	Timestamp Timestamp `json:"x"`
	Y         int       `json:"y"`
	Value     float32   `json:"r"`
}

// Represents the logging of a file import
type FileImportLog struct {
	Id                string
	Md5Checksum       string
	LastDataProcessed time.Time
}
