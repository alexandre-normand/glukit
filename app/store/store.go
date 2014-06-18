// This package provides the interface and functions to deal with the storage of data.
package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/util"
	"math"
	"time"
)

const (
	// Number of GlukitScores to batch in a single PutMulti
	GLUKIT_SCORE_PUT_MULTI_SIZE = 25
)

// Error interface to distinguish between temporary errors from permanent ones
type StoreError struct {
	msg       string
	Temporary bool // Is the error temporary?
}

func (e StoreError) Error() string {
	return e.msg
}

type GlukitScoreScanQuery struct {
	Limit *int
	From  *time.Time
	To    *time.Time
}

var (
	// ErrNoImportedDataFound is returned when the user doesn't have data imported yet.
	ErrNoImportedDataFound = StoreError{"store: no imported data found", true}

	// ErrNoSteadySailorMatchFound is returned when the user doesn't have any steady sailor matching his profile
	ErrNoSteadySailorMatchFound = StoreError{"store: no match for a steady sailor found", true}
)

// GetUserKey returns the GlukitUser datastore key given its email address.
func GetUserKey(context appengine.Context, email string) (key *datastore.Key) {
	return datastore.NewKey(context, "GlukitUser", email, 0, nil)
}

// StoreUserProfile stores a GlukitUser profile to the datastore. If the entry already exists, it is overriden and it is created
// otherwise
func StoreUserProfile(context appengine.Context, updatedAt time.Time, userProfile model.GlukitUser) (key *datastore.Key, err error) {
	key, error := datastore.Put(context, GetUserKey(context, userProfile.Email), &userProfile)
	if error != nil {
		util.Propagate(error)
	}

	return key, nil
}

// GetUserProfile returns the GlukitUser entry associated with the given datastore key. This can be obtained
// by calling GetUserKey.
func GetUserProfile(context appengine.Context, key *datastore.Key) (userProfile *model.GlukitUser, err error) {
	userProfile = new(model.GlukitUser)
	context.Infof("Fetching user profile for key: %s", key.String())
	error := datastore.Get(context, key, userProfile)
	if error != nil {
		return nil, error
	}

	return userProfile, nil
}

// StoreDaysOfReads stores a batch of DayOfReads elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all reads for a single day.
//    2. We have multiple DayOfReads elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfReads is physically stored, see the implementation of apimodel.Store and apimodel.Load.
// Also important to note, this store operation also handles updating the GlukitUser entry with the most recent read, if applicable.
func StoreDaysOfReads(context appengine.Context, userProfileKey *datastore.Key, daysOfReads []apimodel.DayOfGlucoseReads) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfReads))
	for i := range daysOfReads {
		elementKeys[i] = datastore.NewKey(context, "DayOfReads", "", daysOfReads[i].Reads[0].Time.Timestamp, userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of reads", len(elementKeys), len(daysOfReads))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfReads)
	if error != nil {
		context.Warningf("Error writing %d days of reads with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	// Get the time of the batch's last read and update the most recent read timestamp if necessary
	userProfile, err := GetGlukitUserWithKey(context, userProfileKey)
	if err != nil {
		context.Criticalf("Error reading user profile [%s] for its most recent read value: %v", userProfileKey, err)
		return nil, err
	}

	lastDayOfRead := daysOfReads[len(daysOfReads)-1]
	lastRead := lastDayOfRead.Reads[len(lastDayOfRead.Reads)-1]
	if userProfile.MostRecentRead.GetTime().Before(lastRead.GetTime()) {
		context.Infof("Updating most recent read date to %s", lastRead.GetTime())
		userProfile.MostRecentRead = lastRead
		_, err := StoreUserProfile(context, time.Now(), *userProfile)
		if err != nil {
			context.Criticalf("Error storing updated user profile [%s] with most recent read value of %s: %v", userProfileKey, userProfile.MostRecentRead, err)
			return nil, err
		}
	}

	return elementKeys, nil
}

// StoreCalibrationReads stores a batch of DayOfCalibrations elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all calibration reads for a single day.
//    2. We have multiple DayOfReads elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfReads is physically stored, see the implementation of apimodel.Store and apimodel.Load.
func StoreCalibrationReads(context appengine.Context, userProfileKey *datastore.Key, daysOfCalibrationReads []apimodel.DayOfCalibrationReads) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfCalibrationReads))
	for i := range daysOfCalibrationReads {
		elementKeys[i] = datastore.NewKey(context, "DayOfCalibrationReads", "", daysOfCalibrationReads[i].Reads[0].Time.Timestamp, userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of calibration reads", len(elementKeys), len(daysOfCalibrationReads))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfCalibrationReads)
	if error != nil {
		context.Criticalf("Error writing %d days of calibration reads with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetGlucoseReads returns all GlucoseReads given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetGlucoseReads(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (reads []apimodel.GlucoseRead, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for reads between %s and %s to get reads between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfReads").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfReads apimodel.DayOfGlucoseReads
	readsForPeriod := make([]apimodel.GlucoseRead, 0)

	iterator := query.Run(context)
	for _, err := iterator.Next(&daysOfReads); err == nil; _, err = iterator.Next(&daysOfReads) {
		context.Debugf("Loaded batch of %d reads...", len(daysOfReads.Reads))
		readsForPeriod = mergeGlucoseReadArrays(readsForPeriod, daysOfReads.Reads)
	}

	readSlice := apimodel.GlucoseReadSlice(readsForPeriod)
	startIndex, endIndex := apimodel.GetBoundariesOfElementsInRange(readSlice, lowerBound, upperBound)
	filteredReads := readsForPeriod[startIndex : endIndex+1]

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredReads, nil
}

// StoreDaysOfInjections stores a batch of DayOfInjections elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all injections for a single day.
//    2. We have multiple DayOfInjections elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfInjections is physically stored, see the implementation of apimodel.Store and apimodel.Load.
func StoreDaysOfInjections(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []apimodel.DayOfInjections) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfInjections))
	for i := range daysOfInjections {
		elementKeys[i] = datastore.NewKey(context, "DayOfInjections", "", daysOfInjections[i].Injections[0].Time.Timestamp, userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of injections", len(elementKeys), len(daysOfInjections))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfInjections)
	if error != nil {
		context.Criticalf("Error writing %d days of injections with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetInjections returns all Injection entries given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetInjections(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (injections []apimodel.Injection, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for injections between %s and %s to get injections between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfInjections").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfInjections apimodel.DayOfInjections
	injectionsForPeriod := make([]apimodel.Injection, 0)

	iterator := query.Run(context)
	for _, err := iterator.Next(&daysOfInjections); err == nil; _, err = iterator.Next(&daysOfInjections) {
		context.Debugf("Loaded batch of %d injections...", len(daysOfInjections.Injections))
		injectionsForPeriod = mergeInjectionArrays(injectionsForPeriod, daysOfInjections.Injections)
	}

	injectionSlice := apimodel.InjectionSlice(injectionsForPeriod)
	startIndex, endIndex := apimodel.GetBoundariesOfElementsInRange(injectionSlice, lowerBound, upperBound)
	filteredInjections := injectionsForPeriod[startIndex : endIndex+1]

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredInjections, nil
}

// StoreDaysOfMeals stores a batch of DayOfMeals elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all Meals for a single day.
//    2. We have multiple DayOfMeals elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfMeals is physically stored, see the implementation of apimodel.Store and apimodel.Load.
func StoreDaysOfMeals(context appengine.Context, userProfileKey *datastore.Key, daysOfMeals []apimodel.DayOfMeals) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfMeals))
	for i := range daysOfMeals {
		elementKeys[i] = datastore.NewKey(context, "DayOfMeals", "", daysOfMeals[i].Meals[0].Time.Timestamp, userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of carbs", len(elementKeys), len(daysOfMeals))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfMeals)
	if error != nil {
		context.Criticalf("Error writing %d days of carbs with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetMeals returns all Meal entries given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetMeals(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (carbs []apimodel.Meal, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for carbs between %s and %s to get carbs between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfMeals").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfMeals apimodel.DayOfMeals
	mealsForPeriod := make([]apimodel.Meal, 0)

	iterator := query.Run(context)
	for _, err := iterator.Next(&daysOfMeals); err == nil; _, err = iterator.Next(&daysOfMeals) {
		context.Debugf("Loaded batch of %d carbs...", len(daysOfMeals.Meals))
		mealsForPeriod = mergeMealArrays(mealsForPeriod, daysOfMeals.Meals)
	}

	context.Debugf("Filtering between [%s] and [%s], %d carbs: %v", lowerBound, upperBound, len(mealsForPeriod), mealsForPeriod)
	carbSlice := apimodel.MealSlice(mealsForPeriod)
	startIndex, endIndex := apimodel.GetBoundariesOfElementsInRange(carbSlice, lowerBound, upperBound)
	filteredMeals := mealsForPeriod[startIndex : endIndex+1]
	context.Debugf("Finished filtering with %d carbs", len(filteredMeals))

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredMeals, nil
}

// StoreDaysOfExercises stores a batch of DayOfExercises elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all Exercises for a single day.
//    2. We have multiple DayOfExercises elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfExercises is physically stored, see the implementation of apimodel.Store and apimodel.Load.
func StoreDaysOfExercises(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []apimodel.DayOfExercises) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfExercises))
	for i := range daysOfExercises {
		elementKeys[i] = datastore.NewKey(context, "DayOfExercises", "", daysOfExercises[i].Exercises[0].Time.Timestamp, userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of exercises", len(elementKeys), len(daysOfExercises))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfExercises)
	if error != nil {
		context.Criticalf("Error writing %d days of exercises with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetExercises returns all Exercise entries given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetExercises(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (exercises []apimodel.Exercise, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for exercises between %s and %s to get exercises between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfExercises").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfExercises apimodel.DayOfExercises
	exercisesForPeriod := make([]apimodel.Exercise, 0)

	iterator := query.Run(context)
	for _, err := iterator.Next(&daysOfExercises); err == nil; _, err = iterator.Next(&daysOfExercises) {
		context.Debugf("Loaded batch of %d exercises...", len(daysOfExercises.Exercises))
		exercisesForPeriod = mergeExerciseArrays(exercisesForPeriod, daysOfExercises.Exercises)
	}

	exerciseSlice := apimodel.ExerciseSlice(exercisesForPeriod)
	startIndex, endIndex := apimodel.GetBoundariesOfElementsInRange(exerciseSlice, lowerBound, upperBound)
	filteredExercises := exercisesForPeriod[startIndex : endIndex+1]

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredExercises, nil
}

// LogFileImport persist a log of a file import operation. A log entry is actually kept for each distinct file and NOT for every log import
// operation. That is, if we re-import and updated file, we should update the FileImportLog for that file but not create a new one.
// This is used to optimize and not reimport a file that hasn't been updated.
func LogFileImport(context appengine.Context, userProfileKey *datastore.Key, fileImport model.FileImportLog) (key *datastore.Key, err error) {
	key = datastore.NewKey(context, "FileImportLog", fileImport.Id, 0, userProfileKey)

	context.Infof("Emitting a Put for file import log with key [%s] for file id [%s]", key, fileImport.Id)
	key, err = datastore.Put(context, key, &fileImport)
	if err != nil {
		context.Criticalf("Error storing file import log with key [%s] for file id [%s]: %v", key, fileImport.Id, err)
		return nil, err
	}

	return key, nil
}

// GetFileImportLog retrieves a FileImportLog entry for a given file id. If it's the first time we import this file id, a zeroed FileImportLog
// element is returned
func GetFileImportLog(context appengine.Context, userProfileKey *datastore.Key, fileId string) (fileImport *model.FileImportLog, err error) {
	key := datastore.NewKey(context, "FileImportLog", fileId, 0, userProfileKey)

	context.Infof("Reading file import log for file id [%s]", fileId)
	fileImport = new(model.FileImportLog)
	error := datastore.Get(context, key, fileImport)
	if error != nil {
		return nil, error
	}

	return fileImport, nil
}

func GetGlukitUser(context appengine.Context, email string) (key *datastore.Key, userProfile *model.GlukitUser, err error) {
	key = GetUserKey(context, email)
	userProfile, err = GetGlukitUserWithKey(context, key)
	return key, userProfile, err
}

func GetGlukitUserWithKey(context appengine.Context, key *datastore.Key) (userProfile *model.GlukitUser, err error) {
	userProfile, err = GetUserProfile(context, key)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

// GetUserData returns a GlukitUser entry and the boundaries of its most recent complete reads. This is aligned to complete days, meaning that
// it "snaps" to the last day ending at 06:00 am. If the user doesn't have any imported data yet, GetUserData returns ErrNoImportedDataFound
func GetUserData(context appengine.Context, email string) (userProfile *model.GlukitUser, key *datastore.Key, upperBound time.Time, err error) {
	key = GetUserKey(context, email)
	userProfile, err = GetUserProfile(context, key)
	if err != nil {
		return nil, nil, util.GLUKIT_EPOCH_TIME, err
	}

	// If the most recent read is still at the beginning on time, we know no data has been imported yet
	if util.GLUKIT_EPOCH_TIME.Equal(userProfile.MostRecentRead.GetTime()) {
		return userProfile, key, util.GLUKIT_EPOCH_TIME, ErrNoImportedDataFound
	} else {
		upperBound = util.GetEndOfDayBoundaryBefore(userProfile.MostRecentRead.GetTime())
		return userProfile, key, upperBound, nil
	}
}

// FindSteadySailor queries the datastore for others users of the same type of diabetes. It will then select the match that
// has a top glukit score and return that user profile along with the upper boundary for its most recent day of reads.
// The steps involved are:
//    - Find the user profile of the recipient
//    - Query the data store for profile data that matches (using the type of diabetes) in ascending order of score value
//       * A first time for users that are NOT internal
//       * A second time including internal users (if the first one returns no match)
//    - Filter out the recipient profile that could be returned in the search
//    - If match found, get the profile of the steady sailor
func FindSteadySailor(context appengine.Context, recipientEmail string) (sailorProfile *model.GlukitUser, key *datastore.Key, upperBound time.Time, err error) {
	key = GetUserKey(context, recipientEmail)

	recipientProfile, err := GetUserProfile(context, key)
	if err != nil {
		return nil, nil, util.GLUKIT_EPOCH_TIME, err
	}

	context.Debugf("Looking for other diabetes of type [%s]", recipientProfile.DiabetesType)

	// Only get the top-sailor of the same type of diabetes. We might want to throw some randomization in there and pick one of the top 10
	// using cursors or offsets. We need to check at least for two because the recipient user will always be returned by the query.
	// It's more efficient to filter the recipient after the fact than before.
	query := datastore.NewQuery("GlukitUser").
		Filter("diabetesType =", recipientProfile.DiabetesType).
		Order("mostRecentScore.value").Limit(5)

	var steadySailors []model.GlukitUser
	_, err = query.GetAll(context, &steadySailors)
	if err != nil {
		return nil, nil, util.GLUKIT_EPOCH_TIME, err
	}

	context.Debugf("Found a few unfiltered matches [%v]", steadySailors)

	// Stop when we find the first match.

	// First, try for real users
	for i := 0; sailorProfile == nil && i < len(steadySailors); i++ {
		context.Debugf("Loaded steady sailor with email [%s]...", steadySailors[i].Email)
		if steadySailors[i].Email != recipientProfile.Email && !steadySailors[i].Internal {
			sailorProfile = &steadySailors[i]
		}
	}

	// Failing that, get an internal user
	if sailorProfile == nil {
		context.Debugf("Could not find a real user match for recipient [%s], falling back to internal users...", recipientEmail)
		for i := 0; sailorProfile == nil && i < len(steadySailors); i++ {
			context.Debugf("Loaded steady sailor with email [%s]...", steadySailors[i].Email)
			if steadySailors[i].Email != recipientProfile.Email {
				sailorProfile = &steadySailors[i]
			}
		}
	}

	if sailorProfile == nil {
		context.Warningf("No steady sailor match found for user [%s] with type of diabetes [%s]", recipientEmail, recipientProfile.DiabetesType)
		return nil, nil, util.GLUKIT_EPOCH_TIME, ErrNoSteadySailorMatchFound
	} else {
		context.Infof("Found a steady sailor match for user [%s]: healthy [%s]", recipientEmail, sailorProfile.Email)
		upperBound = util.GetEndOfDayBoundaryBefore(sailorProfile.MostRecentRead.GetTime())
		return sailorProfile, GetUserKey(context, sailorProfile.Email), upperBound, nil
	}
}

// StoreGlukitScoreBatch stores a batch of GlukitScores. The array could be of any size. A large batch of GlukitScores
// will be internally split into multiple PutMultis.
func StoreGlukitScoreBatch(context appengine.Context, userEmail string, glukitScores []model.GlukitScore) error {
	parentKey := GetUserKey(context, userEmail)

	totalBatchSize := float64(len(glukitScores))
	for chunkStartIndex := 0; chunkStartIndex < len(glukitScores); chunkStartIndex = chunkStartIndex + GLUKIT_SCORE_PUT_MULTI_SIZE {
		chunkEndIndex := int(math.Min(float64(chunkStartIndex+GLUKIT_SCORE_PUT_MULTI_SIZE), totalBatchSize))
		glukitScoreChunk := glukitScores[chunkStartIndex:chunkEndIndex]
		storeGlukitScoreChunk(context, parentKey, glukitScoreChunk)
	}

	return nil
}

func storeGlukitScoreChunk(context appengine.Context, parentKey *datastore.Key, glukitScoreChunk []model.GlukitScore) (keys []*datastore.Key, err error) {
	context.Debugf("Storing chunk of [%d] glukit scores", len(glukitScoreChunk))

	elementKeys := make([]*datastore.Key, len(glukitScoreChunk))
	for i := range glukitScoreChunk {
		elementKeys[i] = datastore.NewKey(context, "GlukitScore", "", glukitScoreChunk[i].UpperBound.Unix(), parentKey)
	}

	context.Infof("Emitting a PutMulti with [%d] keys for all [%d] glukit scores of chunk", len(elementKeys), len(glukitScoreChunk))
	keys, error := datastore.PutMulti(context, elementKeys, glukitScoreChunk)
	if error != nil {
		context.Criticalf("Error writing [%d] glukit scores with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetGlukitScores returns all GlukitScores for the given email address and matching the query parameters
func GetGlukitScores(context appengine.Context, email string, scanQuery GlukitScoreScanQuery) (scores []model.GlukitScore, err error) {
	key := GetUserKey(context, email)

	context.Infof("Scanning for glukit scores with limit [%d], from [%s], to [%s]", scanQuery.Limit, scanQuery.From, scanQuery.To)

	query := datastore.NewQuery("GlukitScore").Ancestor(key)
	if scanQuery.From != nil {
		query = query.Filter("upperBound >=", scanQuery.From)
	}
	if scanQuery.To != nil {
		query = query.Filter("upperBound <=", scanQuery.To)
	}
	if scanQuery.Limit != nil {
		query = query.Limit(*scanQuery.Limit)
	}
	query = query.Order("-upperBound")

	_, err = query.GetAll(context, &scores)

	if err != datastore.Done {
		util.Propagate(err)
	}

	context.Infof("Found [%d] glukit scores.", len(scores))
	return scores, nil
}
