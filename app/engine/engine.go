// The engine package is where the magic happens. The analysis and insightful bits of data we generate are computed by this
// package. It's the Glukit engine™.
package engine

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
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
	// The max number of days to look back when starting a new batch of calculation
	MAX_CALCULATION_DAYS_TO_LOOK_BACK = 30
)

// January 1st, 2014
var A1C_CALCULATION_START = time.Unix(1388534400, 0)

// CalculateGlukitScore computes the GlukitScore for a given user. This is done in a few steps:
//   1. Get the latest GLUKIT_SCORE_PERIOD days of reads
//   2. For the most recent reads up to READS_REQUIREMENT, calculate the individual score
//      contribution and add it to the GlukitScore.
//   3. If we had enough reads to satisfy the requirements, we return the sum of
//      all individual score contributions.
func CalculateGlukitScore(context context.Context, glukitUser *model.GlukitUser, endOfPeriod time.Time) (glukitScore *model.GlukitScore, err error) {
	// Get the last period's worth of reads
	upperBound := util.GetMidnightUTCBefore(endOfPeriod)
	lowerBound := upperBound.AddDate(0, 0, -1*GLUKIT_SCORE_PERIOD)
	score := model.UNDEFINED_SCORE_VALUE

	log.Debugf(context, "Getting reads for glukit score calculation from [%s] to [%s]", lowerBound, upperBound)
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

		log.Infof(context, "Readcount of [%d] used for glukit score calculation of [%d]", readCount, score)
		if readCount < READS_REQUIREMENT {
			log.Infof(context, "Received only [%d] but required [%d] to calculate valid GlukitScore", readCount, READS_REQUIREMENT)
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
func CalculateIndividualReadScoreWeight(context context.Context, read apimodel.GlucoseRead) (weightedScoreContribution float64) {
	weightedScoreContribution = 0.
	convertedValue, err := read.GetNormalizedValue(apimodel.MG_PER_DL)
	if err != nil {
		util.Propagate(err)
	}
	value := float64(convertedValue)

	if value > model.TARGET_GLUCOSE_VALUE {
		weightedScoreContribution = (value - model.TARGET_GLUCOSE_VALUE) * HIGH_MULTIPLIER
	} else if value < model.TARGET_GLUCOSE_VALUE {
		weightedScoreContribution = -(value - model.TARGET_GLUCOSE_VALUE) * LOW_MULTIPLIER
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

// StartGlukitScoreBatch tries to calculate glukit scores for any week following the most recent calculated score
func StartGlukitScoreBatch(context context.Context, glukitUser *model.GlukitUser) (err error) {
	lowerBoundOfLastScore := glukitUser.MostRecentScore.LowerBound

	// Calculate our minimum allowed lower bound since we don't want to incur the cost of too many reads when
	// a user is coming back from a long absence
	minLowerBound := time.Now().AddDate(0, 0, -1*MAX_CALCULATION_DAYS_TO_LOOK_BACK)

	// Set the lower bound to one day after the last lower bound
	lowerBound := lowerBoundOfLastScore.AddDate(0, 0, -1*GLUKIT_SCORE_PERIOD+1)

	if lowerBound.Before(minLowerBound) {
		lowerBound = minLowerBound
	}

	// Kick off the first chunk of glukit score calculation
	task, err := RunGlukitScoreCalculationChunk.Task(glukitUser.Email, lowerBound)
	if err != nil {
		log.Criticalf(context, "Couldn't schedule the next execution of [%s] for user [%s]. "+
			"This breaks batch calculation of glukit scores for that user!: %v", GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, glukitUser.Email, err)
	}
	taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)
	log.Infof(context, "Queued up first chunk of glukit score calculation for user [%s] and lowerBound [%s]", glukitUser.Email, lowerBound.Format(util.TIMEFORMAT))

	return nil
}

// StartA1CCalculationBatch tries to calculate a1c estimates for any week following the most recent calculated glukit score (a hack, we should have the most recent
// a1c calculation date)
func StartA1CCalculationBatch(context context.Context, glukitUser *model.GlukitUser) (err error) {
	lowerBoundOfLastA1C := glukitUser.MostRecentA1C.LowerBound

	// Uninitialized, default to January 1st, 2014
	if lowerBoundOfLastA1C.Before(A1C_CALCULATION_START) {
		lowerBoundOfLastA1C = A1C_CALCULATION_START
	}

	// Calculate our minimum allowed lower bound since we don't want to incur the cost of too many reads when
	// a user is coming back from a long absence
	minLowerBound := time.Now().AddDate(0, 0, -1*MAX_CALCULATION_DAYS_TO_LOOK_BACK)

	// Set the lower bound to one day after the last lower bound
	lowerBound := lowerBoundOfLastA1C.AddDate(0, 0, -1*A1C_ESTIMATION_SCORE_PERIOD+1)

	if lowerBound.Before(minLowerBound) {
		lowerBound = minLowerBound
	}

	// Kick off the first chunk of glukit score calculation
	task, err := RunA1CCalculationChunk.Task(glukitUser.Email, lowerBound)
	if err != nil {
		log.Criticalf(context, "Couldn't schedule the next execution of [%s] for user [%s]. "+
			"This breaks batch calculation of a1c estimates scores for that user!: %v", A1C_BATCH_CALCULATION_FUNCTION_NAME, glukitUser.Email, err)
	}
	taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)
	log.Infof(context, "Queued up first chunk of a1c calculation for user [%s] and lowerBound [%s]", glukitUser.Email, lowerBound.Format(util.TIMEFORMAT))

	return nil
}
