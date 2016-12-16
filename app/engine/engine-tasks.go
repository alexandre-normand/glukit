package engine

import (
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"golang.org/x/net/context"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"time"
)

var RunGlukitScoreCalculationChunk = delay.Func(GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, func(context context.Context, userEmail string,
	lowerBound time.Time) {
	log.Criticalf(context, "This function purely exists as a workaround to the \"initialization loop\" error that "+
		"shows up because the function calls itself. This implementation defines the same signature as the "+
		"real one which we define in init() to override this implementation!")
})

var RunA1CCalculationChunk = delay.Func(A1C_BATCH_CALCULATION_FUNCTION_NAME, func(context context.Context, userEmail string,
	lowerBound time.Time) {
	log.Criticalf(context, "This function purely exists as a workaround to the \"initialization loop\" error that "+
		"shows up because the function calls itself. This implementation defines the same signature as the "+
		"real one which we define in init() to override this implementation!")
})

const (
	PERIODS_PER_BATCH                            = 6
	BATCH_CALCULATION_QUEUE_NAME                 = "batch-calculation"
	GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME = "runGlukitScoreCalculationChunk"
	A1C_BATCH_CALCULATION_FUNCTION_NAME          = "runA1CCalculationChunk"
)

func RunGlukitScoreBatchCalculation(context context.Context, userEmail string, lowerBound time.Time) {
	glukitUser, _, _, err := store.GetUserData(context, userEmail)
	if _, ok := err.(store.StoreError); err != nil && !ok {
		log.Errorf(context, "We're trying to run a batch glukit score calculation for user [%s] that doesn't exist. "+
			"Got error: %v", userEmail, err)
		return
	}

	bestScore := glukitUser.BestScore
	mostRecentScore := glukitUser.MostRecentScore
	glukitScoreBatch := make([]model.GlukitScore, 0)
	var periodUpperBound time.Time

	log.Debugf(context, "Calculating batch of GlukitScores for user [%s] with current best of [%v] and most recent score of [%v]",
		userEmail, glukitUser.BestScore, mostRecentScore)
	upperBound := lowerBound.AddDate(0, 0, PERIODS_PER_BATCH*GLUKIT_SCORE_PERIOD)

	// Calculate the GlukitScore for every period until now by increment of 1 day. This is a moving score over the last GLUKIT_SCORE_PERIOD that gets a new value every day.
	// This will likely go through a few calculations for which we don't have data yet but this seems like the fair
	// price to pay for making sure we don't stop processing glukit scores because someone might have stopped using their CGM for a week or so.
	for periodUpperBound = lowerBound.AddDate(0, 0, 1); periodUpperBound.Before(time.Now()) && periodUpperBound.Before(upperBound); periodUpperBound = periodUpperBound.AddDate(0, 0, 1) {
		glukitScore, err := CalculateGlukitScore(context, glukitUser, periodUpperBound)
		if err != nil {
			util.Propagate(err)
		}

		if glukitScore.IsBetterThan(glukitUser.BestScore) {
			bestScore = *glukitScore
		}

		if glukitScore.Value != model.UNDEFINED_SCORE_VALUE {
			if periodUpperBound.After(mostRecentScore.UpperBound) {
				mostRecentScore = *glukitScore
			}

			glukitScoreBatch = append(glukitScoreBatch, *glukitScore)
		}
	}

	// Store the batch
	store.StoreGlukitScoreBatch(context, userEmail, glukitScoreBatch)

	// Update the bestScore/LastScoredRead if one of them is different than what was already there
	if bestScore != glukitUser.BestScore || mostRecentScore != glukitUser.MostRecentScore {
		glukitUser.BestScore = bestScore
		glukitUser.MostRecentScore = mostRecentScore
		if _, err := store.StoreUserProfile(context, time.Now(), *glukitUser); err != nil {
			util.Propagate(err)
		} else {
			log.Debugf(context, "Updated glukit user [%s] with an improved GlukitScore of [%v] and most recent score of [%v]",
				glukitUser.Email, bestScore, mostRecentScore)
		}
	}

	// Kick off the next chunk of glukit score calculation
	if !periodUpperBound.Before(upperBound) {
		task, err := RunGlukitScoreCalculationChunk.Task(userEmail, periodUpperBound)
		if err != nil {
			log.Criticalf(context, "Couldn't schedule the next execution of [%s] for user [%s]. "+
				"This breaks batch calculation of glukit scores for that user!: %v", GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, userEmail, err)
		}
		taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)

		log.Infof(context, "Queued up next chunk of glukit score calculation for user [%s] and lowerBound [%s]", userEmail, periodUpperBound.Format(util.TIMEFORMAT))
	} else {
		log.Infof(context, "Done with glukit score calculation for user [%s]", userEmail)
	}
}

func RunA1CBatchCalculation(context context.Context, userEmail string, lowerBound time.Time) {
	glukitUser, _, _, err := store.GetUserData(context, userEmail)
	if _, ok := err.(store.StoreError); err != nil && !ok {
		log.Errorf(context, "We're trying to run a batch of a1c estimates for user [%s] that doesn't exist. "+
			"Got error: %v", userEmail, err)
		return
	}

	mostRecentA1C := glukitUser.MostRecentA1C
	a1cBatch := make([]model.A1CEstimate, 0)
	var periodUpperBound time.Time

	log.Debugf(context, "Calculating batch of A1C estimates for user [%s]", userEmail)
	upperBound := lowerBound.AddDate(0, 0, PERIODS_PER_BATCH*A1C_ESTIMATION_SCORE_PERIOD)

	// Calculate the GlukitScore for every period until now by increment of 1 day. This is a moving estimate over the last A1C_ESTIMATION_SCORE_PERIOD that gets a new value every day.
	// This will likely go through a few calculations for which we don't have data yet but this seems like the fair
	// price to pay for making sure we don't stop processing estimates because someone might have stopped using their CGM for a week or so.
	for periodUpperBound = lowerBound.AddDate(0, 0, 1); periodUpperBound.Before(time.Now()) && periodUpperBound.Before(upperBound); periodUpperBound = periodUpperBound.AddDate(0, 0, 1) {
		a1cEstimate, err := EstimateA1C(context, glukitUser, periodUpperBound)
		if err != nil {
			log.Warningf(context, "Error trying to calculate a1c for user [%s] with upper bound [%s]: %v", userEmail, periodUpperBound, err)
		} else {
			a1cBatch = append(a1cBatch, *a1cEstimate)
			if periodUpperBound.After(mostRecentA1C.UpperBound) {
				mostRecentA1C = *a1cEstimate
			}
		}
	}

	// Store the batch
	store.StoreA1CBatch(context, userEmail, a1cBatch)

	if mostRecentA1C != glukitUser.MostRecentA1C {
		glukitUser.MostRecentA1C = mostRecentA1C

		if _, err := store.StoreUserProfile(context, time.Now(), *glukitUser); err != nil {
			util.Propagate(err)
		} else {
			log.Debugf(context, "Updated glukit user [%s] with a most recent a1c [%v]",
				glukitUser.Email, mostRecentA1C)
		}
	}

	// Kick off the next chunk of glukit score calculation
	if !periodUpperBound.Before(upperBound) {
		task, err := RunA1CCalculationChunk.Task(userEmail, periodUpperBound)
		if err != nil {
			log.Criticalf(context, "Couldn't schedule the next execution of [%s] for user [%s]. "+
				"This breaks batch calculation of glukit scores for that user!: %v", A1C_BATCH_CALCULATION_FUNCTION_NAME, userEmail, err)
		}
		taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)

		log.Infof(context, "Queued up next chunk of a1c calculation for user [%s] and lowerBound [%s]", userEmail, periodUpperBound.Format(util.TIMEFORMAT))
	} else {
		log.Infof(context, "Done with a1c estimation for user [%s]", userEmail)
	}
}
