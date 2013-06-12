package engine

import (
	"models"
	"store"
	"timeutils"
	"appengine"
	"time"
)

const (
	LOW_MULTIPLIER  = 1
	HIGH_MULTIPLIER = 2
	// Two full weeks of reads
	READS_REQUIREMENT = 288*14
)

func CalculateGlukitScore(context appengine.Context, glukitUser *models.GlukitUser) (glukitScore *models.GlukitScore, err error) {
	// Get the last 2 weeks of data plus 2 days
	upperBound := timeutils.GetEndOfDayBoundaryBefore(glukitUser.MostRecentRead)
	lowerBound := upperBound.AddDate(0, 0, -16)
	score := models.UNDEFINED_SCORE_VALUE

	if reads, err := store.GetUserReads(context, glukitUser.Email, lowerBound, upperBound); err != nil {
		return &models.UNDEFINED_SCORE, err
	} else {
		// We might want to do some interpolation of missing reads at some point but for now, we'll only use
		// actual values. Since we know we'll have gaps in a 2 weeks window because of sensor warm-ups, let's
		// just normalize by stopping after the equivalent of full 14 days of reads (assuming most people won't have
		// more than 2 days worth of missing data
		readCount := 0

		for i := 0; i < len(reads) && i < READS_REQUIREMENT; i++ {
			score = score + int64(CalculateIndividualReadScoreWeight(context, reads[i]))
			readCount = readCount + 1
		}
	}

	if (score == models.UNDEFINED_SCORE_VALUE) {
		glukitScore = &models.UNDEFINED_SCORE
	} else {
		glukitScore = &models.GlukitScore{
			Value: score,
			UpperBound: upperBound,
			UpdatedAt: time.Now()}
	}

	return glukitScore, nil
}

// An individual score is either 0 if it's straight on perfection (83) or it's the deviation from 83 weighted
// by whether it's high (multiplier of 2) or lower (multiplier of 1)
func CalculateIndividualReadScoreWeight(context appengine.Context, read models.GlucoseRead) (weightedScoreContribution int) {
	weightedScoreContribution = 0

	if read.Value > 83 {
		weightedScoreContribution = (read.Value - 83)*HIGH_MULTIPLIER
	} else if read.Value < 83 {
		weightedScoreContribution = -(read.Value - 83)*LOW_MULTIPLIER
	}
	//context.Debugf("Calculated individual score of [%d] with read value [%d]", weightedScoreContribution, read.Value)
	return weightedScoreContribution
}
