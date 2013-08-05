// The engine package is where the magic happens. The analysis and insightful bits of data we generate are computed by this
// package. It's the Glukit engine™.
package engine

import (
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"math"
	"time"
)

const (
	// The multiplier applied to any deviation from the target, on the low spectrum (i.e. anything less than 83)
	LOW_MULTIPLIER = 1
	// The multiplier applied to any deviation from the target, on the high spectrum (i.e. anything above 83)
	HIGH_MULTIPLIER = 2
	// Two full weeks of reads
	READS_REQUIREMENT = 288 * 14
)

// CalculateGlukitScore computes the GlukitScore for a given user. This is done in a few steps:
//   1. Get the latest 16 days of reads
//   2. For the most recent reads up to READS_REQUIREMENT, calculate the individual score
//      contribution and add it to the GlukitScore.
//   3. If we had enough reads to satisfy the requirements, we return the sum of
//      all individual score contributions.
func CalculateGlukitScore(context appengine.Context, glukitUser *model.GlukitUser) (glukitScore *model.GlukitScore, err error) {
	// Get the last 2 weeks of data plus 2 days
	upperBound := util.GetEndOfDayBoundaryBefore(glukitUser.MostRecentRead)
	lowerBound := upperBound.AddDate(0, 0, -16)
	score := model.UNDEFINED_SCORE_VALUE

	context.Debugf("Getting reads for glukit score calculation from [%s] to [%s]", lowerBound, upperBound)
	if reads, err := store.GetGlucoseReads(context, glukitUser.Email, lowerBound, upperBound); err != nil {
		return &model.UNDEFINED_SCORE, err
	} else {
		// We might want to do some interpolation of missing reads at some point but for now, we'll only use
		// actual values. Since we know we'll have gaps in a 2 weeks window because of sensor warm-ups, let's
		// just normalize by stopping after the equivalent of full 14 days of reads (assuming most people won't have
		// more than 2 days worth of missing data)
		readCount := 0
		score = 0

		for i := 0; i < len(reads) && i < READS_REQUIREMENT; i++ {
			score = score + int64(CalculateIndividualReadScoreWeight(context, reads[i]))
			readCount = readCount + 1
		}

		context.Infof("Readcount of [%d] used for glukit score calculation of [%d]", readCount, score)
		if readCount < READS_REQUIREMENT {
			context.Warningf("Received only [%d] but required [%d] to calculate valid GlukitScore", readCount, READS_REQUIREMENT)
			return &model.UNDEFINED_SCORE, nil
		}
	}

	if score == model.UNDEFINED_SCORE_VALUE {
		glukitScore = &model.UNDEFINED_SCORE
	} else {
		glukitScore = &model.GlukitScore{
			Value:      score,
			UpperBound: upperBound,
			UpdatedAt:  time.Now()}
	}

	return glukitScore, nil
}

// An individual score is either 0 if it's straight on perfection (83) or it's the deviation from 83 weighted
// by whether it's high (multiplier of 2) or lower (multiplier of 1)
func CalculateIndividualReadScoreWeight(context appengine.Context, read model.GlucoseRead) (weightedScoreContribution int) {
	weightedScoreContribution = 0

	if read.Value > model.TARGET_GLUCOSE_VALUE {
		weightedScoreContribution = (read.Value - model.TARGET_GLUCOSE_VALUE) * HIGH_MULTIPLIER
	} else if read.Value < model.TARGET_GLUCOSE_VALUE {
		weightedScoreContribution = -(read.Value - model.TARGET_GLUCOSE_VALUE) * LOW_MULTIPLIER
	}

	return weightedScoreContribution
}

// CalculateUserFacingScore maps an internal GlukitScore to a user facing value (should be between 0 and 100)
func CalculateUserFacingScore(internal model.GlukitScore) (external int64) {
	internalAsFloat := float64(internal.Value)
	externalAsFloat := 100.689 + 3.8502e-6*internalAsFloat - 0.0631063*math.Sqrt(internalAsFloat-3014.77)
	return int64(externalAsFloat)
}
