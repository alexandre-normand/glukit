package store

import (
	"appengine"
	"appengine/datastore"
	"datautils"
	"math"
	"models"
	"sysutils"
	"time"
	"timeutils"
)

func StoreUserProfile(context appengine.Context, updatedAt time.Time, userProfile models.GlukitUser) (key *datastore.Key, err error) {
	key, error := datastore.Put(context, GetUserKey(context, userProfile.Email), &userProfile)
	if error != nil {
		sysutils.Propagate(error)
	}

	return key, nil
}

func StoreDaysOfReads(context appengine.Context, userProfileKey *datastore.Key, daysOfReads []models.DayOfReads) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfReads))
	for i := range daysOfReads {
		elementKeys[i] = datastore.NewKey(context, "DayOfReads", "", int64(daysOfReads[i].Reads[0].TimeValue), userProfileKey)
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
	if userProfile.MostRecentRead.Before(lastRead.GetTime()) {
		context.Infof("Updating most recent read date to %s", lastRead.GetTime())
		userProfile.MostRecentRead = lastRead.GetTime()
		_, err := StoreUserProfile(context, time.Now(), *userProfile)
		if err != nil {
			context.Criticalf("Error storing updated user profile [%s] with most recent read value of %s: %v", userProfileKey, userProfile.MostRecentRead, err)
			return nil, err
		}
	}

	return elementKeys, nil
}

func StoreDaysOfInjections(context appengine.Context, userProfileKey *datastore.Key, daysOfInjections []models.DayOfInjections) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfInjections))
	for i := range daysOfInjections {
		elementKeys[i] = datastore.NewKey(context, "DayOfInjections", "", int64(daysOfInjections[i].Injections[0].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of injections", len(elementKeys), len(daysOfInjections))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfInjections)
	if error != nil {
		context.Criticalf("Error writing %d days of injections with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

func StoreDaysOfCarbs(context appengine.Context, userProfileKey *datastore.Key, daysOfCarbs []models.DayOfCarbs) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfCarbs))
	for i := range daysOfCarbs {
		elementKeys[i] = datastore.NewKey(context, "DayOfCarbs", "", int64(daysOfCarbs[i].Carbs[0].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of carbs", len(elementKeys), len(daysOfCarbs))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfCarbs)
	if error != nil {
		context.Criticalf("Error writing %d days of carbs with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

func StoreDaysOfExercises(context appengine.Context, userProfileKey *datastore.Key, daysOfExercises []models.DayOfExercises) (keys []*datastore.Key, err error) {
	elementKeys := make([]*datastore.Key, len(daysOfExercises))
	for i := range daysOfExercises {
		elementKeys[i] = datastore.NewKey(context, "DayOfExercises", "", int64(daysOfExercises[i].Exercises[0].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of exercises", len(elementKeys), len(daysOfExercises))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfExercises)
	if error != nil {
		context.Criticalf("Error writing %d days of exercises with keys [%s]: %v", len(elementKeys), elementKeys, error)
		return nil, error
	}

	return elementKeys, nil
}

func LogFileImport(context appengine.Context, userProfileKey *datastore.Key, fileImport models.FileImportLog) (key *datastore.Key, err error) {
	key = datastore.NewKey(context, "FileImportLog", fileImport.Id, 0, userProfileKey)

	context.Infof("Emitting a Put for file import log with key [%s] for file id [%s]", key, fileImport.Id)
	key, err = datastore.Put(context, key, &fileImport)
	if err != nil {
		context.Criticalf("Error storing file import log with key [%s] for file id [%s]: %v", key, fileImport.Id, err)
		return nil, err
	}

	return key, nil
}

func GetFileImportLog(context appengine.Context, userProfileKey *datastore.Key, fileId string) (fileImport *models.FileImportLog, err error) {
	key := datastore.NewKey(context, "FileImportLog", fileId, 0, userProfileKey)

	context.Infof("Reading file import log for file id [%s]", fileId)
	fileImport = new(models.FileImportLog)
	error := datastore.Get(context, key, fileImport)
	if error != nil {
		return nil, error
	}

	return fileImport, nil
}

func GetUserKey(context appengine.Context, email string) (key *datastore.Key) {
	return datastore.NewKey(context, "GlukitUser", email, 0, nil)
}

func GetUserProfile(context appengine.Context, key *datastore.Key) (userProfile *models.GlukitUser, err error) {
	userProfile = new(models.GlukitUser)
	context.Infof("Fetching user profile for key: %s", key.String())
	error := datastore.Get(context, key, userProfile)
	if error != nil {
		return nil, error
	}

	return userProfile, nil
}
func GetUserData(context appengine.Context, email string) (userProfile *models.GlukitUser, key *datastore.Key, lowerBound time.Time, upperBound time.Time, err error) {
	key = GetUserKey(context, email)
	userProfile, err = GetUserProfile(context, key)
	if err != nil {
		return nil, nil, time.Unix(0, 0), time.Unix(math.MaxInt64, math.MaxInt64), err
	}

	upperBound = timeutils.GetEndOfDayBoundaryBefore(userProfile.MostRecentRead)
	lowerBound = upperBound.Add(time.Duration(-24 * time.Hour))

	return userProfile, key, lowerBound, upperBound, nil
}

func GetUserReads(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (reads []models.GlucoseRead, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for reads between %s and %s to get reads between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfReads").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfReads models.DayOfReads

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfReads); err == nil; _, err = iterator.Next(&daysOfReads) {
		batchSize := len(daysOfReads.Reads) - count
		context.Debugf("Loaded batch of %d reads...", batchSize)
		count = len(daysOfReads.Reads)
	}

	filteredReads := datautils.FilterReads(daysOfReads.Reads, lowerBound, upperBound)

	if err != datastore.Done {
		sysutils.Propagate(err)
	}

	return filteredReads, nil
}

func GetUserInjections(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (injections []models.Injection, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for injections between %s and %s to get injections between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfInjections").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfInjections models.DayOfInjections

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfInjections); err == nil; _, err = iterator.Next(&daysOfInjections) {
		batchSize := len(daysOfInjections.Injections) - count
		context.Debugf("Loaded batch of %d injections...", batchSize)
		count = len(daysOfInjections.Injections)
	}

	filteredInjections := datautils.FilterInjections(daysOfInjections.Injections, lowerBound, upperBound)

	if err != datastore.Done {
		sysutils.Propagate(err)
	}

	return filteredInjections, nil
}

func GetUserCarbs(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (carbs []models.Carb, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for carbs between %s and %s to get carbs between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfCarbs").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfCarbs models.DayOfCarbs

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfCarbs); err == nil; _, err = iterator.Next(&daysOfCarbs) {
		batchSize := len(daysOfCarbs.Carbs) - count
		context.Debugf("Loaded batch of %d carbs...", batchSize)
		count = len(daysOfCarbs.Carbs)
	}

	filteredCarbs := datautils.FilterCarbs(daysOfCarbs.Carbs, lowerBound, upperBound)

	if err != datastore.Done {
		sysutils.Propagate(err)
	}

	return filteredCarbs, nil
}

// Retrieves exercise values for the specified lower and upper bounds for the user identified by the email address
func GetUserExercises(context appengine.Context, email string, lowerBound time.Time, upperBound time.Time) (exercises []models.Exercise, err error) {
	key := GetUserKey(context, email)

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24 * time.Hour))
	scanEnd := upperBound.Add(time.Duration(24 * time.Hour))

	context.Infof("Scanning for exercises between %s and %s to get exercises between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	query := datastore.NewQuery("DayOfExercises").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var daysOfExercises models.DayOfExercises

	iterator := query.Run(context)
	count := 0
	for _, err := iterator.Next(&daysOfExercises); err == nil; _, err = iterator.Next(&daysOfExercises) {
		batchSize := len(daysOfExercises.Exercises) - count
		context.Debugf("Loaded batch of %d exercises...", batchSize)
		count = len(daysOfExercises.Exercises)
	}

	filteredExercises := datautils.FilterExercises(daysOfExercises.Exercises, lowerBound, upperBound)

	if err != datastore.Done {
		sysutils.Propagate(err)
	}

	return filteredExercises, nil
}
