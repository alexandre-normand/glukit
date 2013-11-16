// The engine package is where the magic happens. The analysis and insightful bits of data we generate are computed by this
// package. It's the Glukit engine™.
package engine

import (
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/taskqueue"
	"math"
	"time"
)

const (
	// The multiplier applied to any deviation from the target, on the low spectrum (i.e. anything less than 83)
	LOW_MULTIPLIER = 1
	// The multiplier applied to any deviation from the target, on the high spectrum (i.e. anything above 83)
	HIGH_MULTIPLIER = 2
	// Glukit score calculation period
	GLUKIT_SCORE_PERIOD = 7
	// One period of reads minus on day for potential data gaps
	READS_REQUIREMENT = 288 * (GLUKIT_SCORE_PERIOD - 1)
	// The current Glukit scoring version
	SCORING_VERSION = 1
)

// CalculateGlukitScore computes the GlukitScore for a given user. This is done in a few steps:
//   1. Get the latest GLUKIT_SCORE_PERIOD days of reads
//   2. For the most recent reads up to READS_REQUIREMENT, calculate the individual score
//      contribution and add it to the GlukitScore.
//   3. If we had enough reads to satisfy the requirements, we return the sum of
//      all individual score contributions.
func CalculateGlukitScore(context appengine.Context, glukitUser *model.GlukitUser, endOfPeriod time.Time) (glukitScore *model.GlukitScore, err error) {
	// Get the last period's worth of reads
	upperBound := util.GetMidnightUTCBefore(endOfPeriod)
	lowerBound := upperBound.AddDate(0, 0, -1*GLUKIT_SCORE_PERIOD)
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
			context.Infof("Received only [%d] but required [%d] to calculate valid GlukitScore", readCount, READS_REQUIREMENT)
			return &model.UNDEFINED_SCORE, nil
		}
	}

	if score == model.UNDEFINED_SCORE_VALUE {
		glukitScore = &model.UNDEFINED_SCORE
	} else {
		glukitScore = &model.GlukitScore{
			Value:          score,
			LowerBound:     lowerBound,
			UpperBound:     upperBound,
			CalculatedOn:   time.Now(),
			ScoringVersion: SCORING_VERSION}
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
func CalculateUserFacingScore(internal model.GlukitScore) (external *int64) {
	if internal.Value == model.UNDEFINED_SCORE_VALUE {
		return nil
	} else {
		internalAsFloat := float64(internal.Value)
		externalAsFloat := 100 + 1.043e-9*math.Pow(internalAsFloat, 2) + 6.517e-22*math.Pow(internalAsFloat, 4) - 0.0003676*internalAsFloat - 1.434e-15*math.Pow(internalAsFloat, 3)
		externalAsInt := int64(externalAsFloat)
		return &externalAsInt
	}
}

// CalculateGlukitScoreBatch tries to calculate glukit scores for any week following the most recent calculated score
func CalculateGlukitScoreBatch(context appengine.Context, glukitUser *model.GlukitUser) (err error) {
	lastScoredRead := glukitUser.MostRecentScore.UpperBound

	// Kick off the first chunk of glukit score calculation
	task, err := RunGlukitScoreCalculationChunk.Task(glukitUser.Email, lastScoredRead)
	if err != nil {
		context.Criticalf("Couldn't schedule the next execution of runGlukitScoreCalculationChunk for user [%s]. "+
			"This breaks batch calculation of glukit scores for that user!: %v", glukitUser.Email, err)
	}
	taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)
	context.Infof("Queued up first chunk of glukit score calculation for user [%s] and lowerBound [%s]", glukitUser.Email, lastScoredRead.Format(util.TIMEFORMAT))

	return nil
}
