package engine

import (
	"appengine"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/lib/github.com/grd/stat"
	"time"
)

const (
	A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS = 90
)

// CalculateA1CEstimate calculates an estimate of a a1c given the last 3 months of data. The current algo is naively assuming that the average of the last
// 3 months will be an approximation of the a1c.
func CalculateA1CEstimate(context appengine.Context, reads []apimodel.GlucoseRead) (glukitScore *model.A1CEstimate, err error) {
	lowerBound := reads[0].GetTime()
	upperBound := reads[len(reads)-1].GetTime()

	context.Debugf("Estimating a1c from [%d] reads from [%s] to [%s]", len(reads), lowerBound, upperBound)

	coverage := upperBound.Sub(lowerBound)
	days := coverage / time.Hour

	if days < A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS {
		return &model.UNDEFINED_A1C_ESTIMATE, nil
	} else {
		average := stat.Mean(model.ReadStatsSlice(reads))
		a1c := (average + 46.7) / 28.7
		return &model.A1CEstimate{
			Value:          a1c,
			LowerBound:     lowerBound,
			UpperBound:     upperBound,
			CalculatedOn:   time.Now(),
			ScoringVersion: SCORING_VERSION}, nil
	}
}
