package engine

import (
	"app/model"
	"app/store"
	"app/util"
	"appengine"
	"appengine/delay"
	"appengine/taskqueue"
	"time"
)

var RunGlukitScoreCalculationChunk = delay.Func(GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME, func(context appengine.Context, userEmail string,
	lowerBound time.Time) {
	context.Criticalf("This function purely exists as a workaround to the \"initialization loop\" error that " +
		"shows up because the function calls itself. This implementation defines the same signature as the " +
		"real one which we define in init() to override this implementation!")
})

const (
	PERIODS_PER_BATCH                            = 6
	BATCH_CALCULATION_QUEUE_NAME                 = "batch-calculation"
	GLUKIT_SCORE_BATCH_CALCULATION_FUNCTION_NAME = "runGlukitScoreCalculationChunk"
)

func RunGlukitScoreBatchCalculation(context appengine.Context, userEmail string, lowerBound time.Time) {
	glukitUser, _, _, _, err := store.GetUserData(context, userEmail)
	if _, ok := err.(store.StoreError); err != nil && !ok {
		context.Errorf("We're trying to run an batch glukit score calculation for user [%s] that doesn't exist. "+
			"Got error: %v", userEmail, err)
		return
	}

	bestScore := glukitUser.BestScore
	mostRecentScore := glukitUser.MostRecentScore
	glukitScoreBatch := make([]model.GlukitScore, 0)
	var periodUpperBound time.Time

	context.Debugf("Calculating batch of GlukitScores for user [%s] with current best of [%v] and most recent score of [%v]",
		userEmail, glukitUser.BestScore, glukitUser.MostRecentScore)
	upperBound := lowerBound.AddDate(0, 0, PERIODS_PER_BATCH*GLUKIT_SCORE_PERIOD)

	// Calculate the GlukitScore for every period until now. This will likely go through a few calculations for which we don't have data yet but this seems like the fair
	// price to pay for making sure we don't stop processing glukit scores because someone might have stopped using their CGM for a week or so.
	for periodUpperBound = lowerBound.AddDate(0, 0, GLUKIT_SCORE_PERIOD); periodUpperBound.Before(time.Now()) && periodUpperBound.Before(upperBound); periodUpperBound = periodUpperBound.AddDate(0, 0, GLUKIT_SCORE_PERIOD) {
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
			context.Debugf("Updated glukit user [%s] with an improved GlukitScore of [%v] and most recent score of [%v]",
				glukitUser.Email, bestScore, mostRecentScore)
		}
	}

	// Kick off the next chunk of glukit score calculation
	if !periodUpperBound.Before(upperBound) {
		task, err := RunGlukitScoreCalculationChunk.Task(userEmail, periodUpperBound)
		if err != nil {
			context.Criticalf("Couldn't schedule the next execution of runGlukitScoreCalculationChunk for user [%s]. "+
				"This breaks batch calculation of glukit scores for that user!: %v", userEmail, err)
		}
		taskqueue.Add(context, task, BATCH_CALCULATION_QUEUE_NAME)

		context.Infof("Queued up next chunk of glukit score calculation for user [%s] and lowerBound [%s]", userEmail, periodUpperBound.Format(util.TIMEFORMAT))
	} else {
		context.Infof("Done with glukit score calculation for user [%s]", userEmail)
	}
}
