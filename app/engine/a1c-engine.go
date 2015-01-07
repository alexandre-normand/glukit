package engine

import (
	"appengine"
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/github.com/grd/stat"
	"time"
)

const (
	A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS = 90

	// A1C estimation scoring period requirement
	A1C_ESTIMATION_SCORE_PERIOD = 95
)

// CalculateA1CEstimate calculates an estimate of a a1c given the last 3 months of data. The current algo is naively assuming that the average of the last
// 3 months will be an approximation of the a1c.
func CalculateA1CEstimate(context appengine.Context, reads []apimodel.GlucoseRead) (a1c *model.A1CEstimate, err error) {
	if len(reads) == 0 {
		return nil, errors.New(fmt.Sprintf("Insufficient read coverage to estimate a1c, got no reads"))
	}
	lowerBound := reads[0].GetTime()
	upperBound := reads[len(reads)-1].GetTime()

	context.Debugf("Estimating a1c from [%d] reads from [%s] to [%s]", len(reads), lowerBound, upperBound)

	coverage := upperBound.Sub(lowerBound)
	days := coverage / (time.Hour * 24)
	context.Debugf("Coverage is [%d]", days)

	if days < A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS {
		context.Debugf(fmt.Sprintf("Insufficient read coverage to estimate a1c, got [%d] days but requires [%d]", days, A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS))
		return nil, errors.New(fmt.Sprintf("Insufficient read coverage to estimate a1c, got [%d] days but requires [%d]", days, A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS))
	} else {
		average := stat.Mean(model.ReadStatsSlice(reads))
		a1c := (average + 77.3) / 35.6
		context.Debugf("Estimated a1c is [%f]", a1c)
		return &model.A1CEstimate{
			Value:          a1c,
			LowerBound:     lowerBound,
			UpperBound:     upperBound,
			CalculatedOn:   time.Now(),
			ScoringVersion: SCORING_VERSION}, nil
	}
}

func EstimateA1C(context appengine.Context, glukitUser *model.GlukitUser, endOfPeriod time.Time) (a1c *model.A1CEstimate, err error) {
	// Get the last period's worth of reads
	upperBound := util.GetMidnightUTCBefore(endOfPeriod)
	lowerBound := upperBound.AddDate(0, 0, -1*A1C_ESTIMATION_SCORE_PERIOD)
	context.Debugf("fuck this")
	context.Debugf("Getting reads for a1c estimate calculation from [%s] to [%s]", lowerBound, upperBound)
	if reads, err := store.GetGlucoseReads(context, glukitUser.Email, lowerBound, upperBound); err != nil {
		return &model.UNDEFINED_A1C_ESTIMATE, err
	} else {
		return CalculateA1CEstimate(context, reads)
	}
}
