// This package provides the interface and functions to deal with the storage of data.
package store

import (
	"app/model"
	"app/util"
	"appengine"
	"appengine/datastore"
	"math"
	"time"
)

const (
	// Number of GlukitScores to batch in a single PutMulti
	GLUKIT_SCORE_PUT_MULTI_SIZE = 200
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
// For details of how a single element of DayOfReads is physically stored, see the implementation of model.Store and model.Load.
// Also important to note, this store operation also handles updating the GlukitUser entry with the most recent read, if applicable.
func StoreDaysOfReads(context appengine.Context, userProfileKey *datastore.Key, daysOfReads []model.DayOfGlucoseReads) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfReads))
	for i := range daysOfReads {
		elementKeys[i] = datastore.NewKey(context, "DayOfReads", "", int64(daysOfReads[i].Reads[0].Timestamp), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of reads", len(elementKeys), len(daysOfReads))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfReads)
	if error != nil {
		context.Criticalf("Error writing %d days of reads with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	// Get the time of the batch's last read and update the most recent read timestamp if necessary
	userProfile, err := GetUserProfile(context, userProfileKey)
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

// GetGlucoseReads returns all GlucoseReads given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetGlucoseReads(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (reads []model.GlucoseRead, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for reads between %s and %s to get reads between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfReads").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfReads model.DayOfGlucoseReads

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfReads); err == nil; _, err = iterator.Next(&daysOfReads) {
		batchSize := len(daysOfReads.Reads) - count
		context.Debugf("Loaded batch of %d reads...", batchSize)
		count = len(daysOfReads.Reads)
	}

	filteredReads := filterReads(daysOfReads.Reads, lowerBound, upperBound)

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredReads, nil
}

// StoreDaysOfInjections stores a batch of DayOfInjections elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all injections for a single day.
//    2. We have multiple DayOfInjections elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfInjections is physically stored, see the implementation of model.Store and model.Load.
func StoreDaysOfInjections(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []model.DayOfInjections) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfInjections))
	for i := range daysOfInjections {
		elementKeys[i] = datastore.NewKey(context, "DayOfInjections", "", int64(daysOfInjections[i].Injections[0].Timestamp), userProfileKey)
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
func GetInjections(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (injections []model.Injection, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for injections between %s and %s to get injections between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfInjections").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfInjections model.DayOfInjections

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfInjections); err == nil; _, err = iterator.Next(&daysOfInjections) {
		batchSize := len(daysOfInjections.Injections) - count
		context.Debugf("Loaded batch of %d injections...", batchSize)
		count = len(daysOfInjections.Injections)
	}

	filteredInjections := filterInjections(daysOfInjections.Injections, lowerBound, upperBound)

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredInjections, nil
}

// StoreDaysOfCarbs stores a batch of DayOfCarbs elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all Carbs for a single day.
//    2. We have multiple DayOfCarbs elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfCarbs is physically stored, see the implementation of model.Store and model.Load.
func StoreDaysOfCarbs(context appengine.Context, userProfileKey *datastore.Key, daysOfCarbs []model.DayOfCarbs) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfCarbs))
	for i := range daysOfCarbs {
		elementKeys[i] = datastore.NewKey(context, "DayOfCarbs", "", int64(daysOfCarbs[i].Carbs[0].Timestamp), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of carbs", len(elementKeys), len(daysOfCarbs))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfCarbs)
	if error != nil {
		context.Criticalf("Error writing %d days of carbs with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

// GetCarbs returns all Carb entries given a user's email address and the time boundaries. Not that the boundaries are both inclusive.
func GetCarbs(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (carbs []model.Carb, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for carbs between %s and %s to get carbs between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfCarbs").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfCarbs model.DayOfCarbs

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfCarbs); err == nil; _, err = iterator.Next(&daysOfCarbs) {
		batchSize := len(daysOfCarbs.Carbs) - count
		context.Debugf("Loaded batch of %d carbs...", batchSize)
		count = len(daysOfCarbs.Carbs)
	}

	context.Debugf("Filtering between [%s] and [%s], %d carbs: %v", lowerBound, upperBound, len(daysOfCarbs.Carbs), daysOfCarbs.Carbs)
	filteredCarbs := filterCarbs(daysOfCarbs.Carbs, lowerBound, upperBound)
	context.Debugf("Finished filterting with %d carbs", len(filteredCarbs))

	if err != datastore.Done {
		util.Propagate(err)
	}

	return filteredCarbs, nil
}

// StoreDaysOfExercises stores a batch of DayOfExercises elements. It is a optimized operation in that:
//    1. One element represents a relatively short-and-wide entry of all Exercises for a single day.
//    2. We have multiple DayOfExercises elements and we use a PutMulti to make this faster.
// For details of how a single element of DayOfExercises is physically stored, see the implementation of model.Store and model.Load.
func StoreDaysOfExercises(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []model.DayOfExercises) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfExercises))
	for i := range daysOfExercises {
		elementKeys[i] = datastore.NewKey(context, "DayOfExercises", "", int64(daysOfExercises[i].Exercises[0].Timestamp), userProfileKey)
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
func GetExercises(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (exercises []model.Exercise, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for exercises between %s and %s to get exercises between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfExercises").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfExercises model.DayOfExercises

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfExercises); err == nil; _, err = iterator.Next(&daysOfExercises) {
		batchSize := len(daysOfExercises.Exercises) - count
		context.Debugf("Loaded batch of %d exercises...", batchSize)
		count = len(daysOfExercises.Exercises)
	}

	filteredExercises := filterExercises(daysOfExercises.Exercises, lowerBound, upperBound)

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

// GetUserData returns a GlukitUser entry and the boundaries of its most recent complete day of reads. If the user doesn't have any imported data yet,
// GetUserData returns ErrNoImportedDataFound
func GetUserData(context appengine.Context, email string) (userProfile *model.GlukitUser, key *datastore.Key, lowerBound time.Time, upperBound time.Time, err error) {
	key = GetUserKey(context, email)
	userProfile, err = GetUserProfile(context, key)
	if err != nil {
		return nil, nil, time.Unix(0, 0), time.Unix(math.MaxInt64, math.MaxInt64), err
	}

	// If the most recent read is still at the beginning on time, we know no data has been imported yet
	if util.GLUKIT_EPOCH_TIME.Equal(userProfile.MostRecentRead.GetTime()) {
		return userProfile, key, util.GLUKIT_EPOCH_TIME, util.GLUKIT_EPOCH_TIME, ErrNoImportedDataFound
	} else {
		upperBound = util.GetEndOfDayBoundaryBefore(util.GetLocalTimeInProperLocation(userProfile.MostRecentRead.LocalTime, userProfile.MostRecentRead.GetTime()))
		lowerBound = upperBound.Add(time.Duration(-24 * time.Hour))
		return userProfile, key, lowerBound, upperBound, nil
	}
}

// FindSteadySailor queries the datastore for others users of the same type of diabetes. It will then select the match that
// has a top glukit score and return that user profile along with the boundaries for its most recent day of reads.
// The steps involved are:
//    - Find the user profile of the recipient
//    - Query the data store for profile data that matches (using the type of diabetes) in ascending order of score value
//       * A first time for users that are NOT internal
//       * A second time including internal users (if the first one returns no match)
//    - Filter out the recipient profile that could be returned in the search
//    - If match found, get the profile of the steady sailor
func FindSteadySailor(context appengine.Context, recipientEmail string) (sailorProfile *model.GlukitUser, key *datastore.Key, lowerBound time.Time, upperBound time.Time, err error) {
	key = GetUserKey(context, recipientEmail)

	recipientProfile, err := GetUserProfile(context, key)
	if err != nil {
		return nil, nil, time.Unix(0, 0), time.Unix(math.MaxInt64, math.MaxInt64), err
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
		return nil, nil, time.Unix(0, 0), time.Unix(math.MaxInt64, math.MaxInt64), err
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
		return nil, nil, util.GLUKIT_EPOCH_TIME, util.GLUKIT_EPOCH_TIME, ErrNoSteadySailorMatchFound
	} else {
		context.Warningf("Found a steady sailor match for user [%s]: healthy [%s]", recipientEmail, sailorProfile.Email)
		upperBound = util.GetEndOfDayBoundaryBefore(util.GetLocalTimeInProperLocation(sailorProfile.MostRecentRead.LocalTime, sailorProfile.MostRecentRead.GetTime()))
		lowerBound = upperBound.Add(time.Duration(-24 * time.Hour))
		return sailorProfile, GetUserKey(context, sailorProfile.Email), lowerBound, upperBound, nil
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
