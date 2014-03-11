/*
Package models provides model types
*/
package model

import (
	"github.com/alexandre-normand/glukit/app/util"
	"math"
	"time"
)

const (
	TARGET_GLUCOSE_VALUE    = 83
	UNDEFINED_READ          = -1
	EXERCISE_VALUE_FORMAT   = "%d,%s"
	UNDEFINED_SCORE_VALUE   = int64(math.MaxInt64)
	DEFAULT_LOOKBACK_PERIOD = time.Duration(-7*24) * time.Hour
)

// "Dynamic" constants, those should never be updated
var UNDEFINED_SCORE = GlukitScore{Value: UNDEFINED_SCORE_VALUE, LowerBound: util.GLUKIT_EPOCH_TIME, UpperBound: util.GLUKIT_EPOCH_TIME, CalculatedOn: util.GLUKIT_EPOCH_TIME, ScoringVersion: -1}

// Represents a cartesian coordinate
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Represents the logging of a file import
type FileImportLog struct {
	Id                string
	Md5Checksum       string
	LastDataProcessed time.Time
}
