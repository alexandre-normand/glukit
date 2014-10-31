package model

import (
	"github.com/alexandre-normand/glukit/app/util"
	"math"
	"time"
)

// A1CEstimate is a calculated estimate of an a1c. The lower and upper bounds
// should match the date of the first and last read of the period
// used to calculate the score. The scoring version represents
// the version of the calculation algorithm used to calculate a given estimate
// It is used to discard/recalculate older versions of glukit
// scores in the eventuality where we change how we calculate the internal
// estimation.
type A1CEstimate struct {
	Value          float64   `datastore:"value"`
	LowerBound     time.Time `datastore:"lowerBound"`
	UpperBound     time.Time `datastore:"upperBound"`
	CalculatedOn   time.Time `datastore:"calculatedOn"`
	ScoringVersion int       `datastore:"scoringVersion`
}

const (
	UNDEFINED_A1C_VALUE = math.SmallestNonzeroFloat64
)

// "Dynamic" constants, those should never be updated
var UNDEFINED_A1C_ESTIMATE = A1CEstimate{Value: UNDEFINED_A1C_VALUE, LowerBound: util.GLUKIT_EPOCH_TIME, UpperBound: util.GLUKIT_EPOCH_TIME, CalculatedOn: util.GLUKIT_EPOCH_TIME, ScoringVersion: -1}
