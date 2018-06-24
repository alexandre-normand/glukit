package model

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"time"
)

// Represents a GlukitUser profile
type GlukitUser struct {
	Email           string               `datastore:"email"`
	FirstName       string               `datastore:"firstName,noindex"`
	LastName        string               `datastore:"lastName,noindex"`
	DateOfBirth     time.Time            `datastore:"birthdate"`
	DiabetesType    string               `datastore:"diabetesType"`
	Timezone        string               `datastore:"timezoneId,noindex"`
	LastUpdated     time.Time            `datastore:"lastUpdated"`
	MostRecentRead  apimodel.GlucoseRead `datastore:"mostRecentRead"`
	BestScore       GlukitScore          `datastore:"bestScore"`
	MostRecentScore GlukitScore          `datastore:"mostRecentScore"`
	Internal        bool                 `datastore:"internal"`
	PictureUrl      string               `datastore:"pictureUrl,noindex"`
	AccountCreated  time.Time            `datastore:"joinedOn"`
	MostRecentA1C   A1CEstimate          `datastore:"mostRecentA1C"`
}

// Represents a GlukitScore value, the lower and upper bounds
// should match the date of the first and last read of the period
// used to calculate the score. The GlukitScore scoring version represents
// the version of the calculation algorithm used to calculate a given
// score. It is used to discard/recalculate older versions of glukit
// scores in the eventuality where we change how we calculate the internal
// score.
type GlukitScore struct {
	Value          int64     `datastore:"value"`
	LowerBound     time.Time `datastore:"lowerBound"`
	UpperBound     time.Time `datastore:"upperBound"`
	CalculatedOn   time.Time `datastore:"calculatedOn"`
	ScoringVersion int       `datastore:"scoringVersion`
}

// Type of diabetes
const (
	DIABETES_TYPE_1 = "T1"
	DIABETES_TYPE_2 = "T2"
)

// Compares one GlukitScore with another and returns true if
// the score is better than the reference score passed in argument
// i.e. newValue.IsBetterThan(oldValue) would return true if the
// newValue's score is better than the oldValue's score.
func (score GlukitScore) IsBetterThan(reference GlukitScore) bool {
	return score.Value < reference.Value
}
