package engine

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/grd/stat"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"sort"
	"time"
)

const (
	A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS = 90

	// A1C estimation scoring period requirement
	A1C_ESTIMATION_SCORE_PERIOD = 95

	A1C_SCORING_VERSION = 3
)

// CalculateA1CEstimate calculates an estimate of a a1c given the last 3 months of data. The current algo is naively assuming that the average of the last
// 3 months will be an approximation of the a1c.
func CalculateA1CEstimate(context context.Context, reads []apimodel.GlucoseRead) (a1c *model.A1CEstimate, err error) {
	if len(reads) == 0 {
		return nil, errors.New(fmt.Sprintf("Insufficient read coverage to estimate a1c, got no reads"))
	}
	lowerBound := reads[0].GetTime()
	upperBound := reads[len(reads)-1].GetTime()

	log.Debugf(context, "Estimating a1c from [%d] reads from [%s] to [%s]", len(reads), lowerBound, upperBound)

	coverage := upperBound.Sub(lowerBound)
	days := coverage / (time.Hour * 24)
	log.Debugf(context, "Coverage is [%d]", days)

	if days < A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS {
		return nil, errors.New(fmt.Sprintf("Insufficient read coverage to estimate a1c, got [%d] days but requires [%d]", days, A1C_READ_COVERAGE_REQUIREMENT_IN_DAYS))
	} else {
		sortedReads := model.ReadStatsSlice(reads)
		sort.Sort(sortedReads)
		median := stat.MedianFromSortedData(sortedReads)
		//a1c := (average + 77.3) / 35.6
		a1c := (median + 77.3) / 35.6
		log.Debugf(context, "Estimated a1c is [%f]", a1c)
		return &model.A1CEstimate{
			Value:          a1c,
			LowerBound:     lowerBound,
			UpperBound:     upperBound,
			CalculatedOn:   time.Now(),
			ScoringVersion: A1C_SCORING_VERSION}, nil
	}
}

func EstimateA1C(context context.Context, glukitUser *model.GlukitUser, endOfPeriod time.Time) (a1c *model.A1CEstimate, err error) {
	// Get the last period's worth of reads
	upperBound := util.GetMidnightUTCBefore(endOfPeriod)
	lowerBound := upperBound.AddDate(0, 0, -1*A1C_ESTIMATION_SCORE_PERIOD)

	log.Debugf(context, "Getting reads for a1c estimate calculation from [%s] to [%s]", lowerBound, upperBound)
	if reads, err := store.GetGlucoseReads(context, glukitUser.Email, lowerBound, upperBound); err != nil {
		return &model.UNDEFINED_A1C_ESTIMATE, err
	} else {
		return CalculateA1CEstimate(context, reads)
	}
}
