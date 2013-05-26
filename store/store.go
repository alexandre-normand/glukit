package store

import (
	"encoding/json"
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"time"
	"sysutils"
	"models"
	"timeutils"
	"datautils"
)

func StoreUserProfile(context appengine.Context, updatedAt time.Time, userProfile models.GlukitUser) (key *datastore.Key, err error) {
	key, error := datastore.Put(context, GetUserKey(context, userProfile.Email), &userProfile)
	if error != nil {
		sysutils.Propagate(error)
	}

	return key, nil
}

func StoreDaysOfReads(context appengine.Context, userProfileKey *datastore.Key, daysOfReads []models.DayOfReads) (keys []*datastore.Key, err error) {
	context.Infof("Called StoreDaysOfReads")
	elementKeys := make([] *datastore.Key, len(daysOfReads))
	for i := range daysOfReads {
		elementKeys[i] = datastore.NewKey(context, "DayOfReads", "", int64(daysOfReads[i].Reads[0].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d days of reads", len(elementKeys), len(daysOfReads))
	keys, error := datastore.PutMulti(context, elementKeys, daysOfReads)
	if error != nil {
		context.Criticalf("Error writing %d days of reads with keys [%s]", len(elementKeys), elementKeys)
		return nil, error
	}

	// Get the time of the batch's last read and update the most recent read timestamp if necessary
	userProfile, err := GetUserProfile(context, userProfileKey)
	if err != nil {
		context.Criticalf("Error reading user profile [%s] for its most recent read value: %v", userProfileKey, err)
		return nil, err
	}

	lastDayOfRead := daysOfReads[len(daysOfReads) - 1]
	lastRead := lastDayOfRead.Reads[len(lastDayOfRead.Reads) - 1]
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

func StoreReads(context appengine.Context, userProfileKey *datastore.Key, reads []models.MeterRead) (key *datastore.Key, err error) {
	elementKey := datastore.NewKey(context, "DayOfReads", "", int64(reads[0].TimeValue), userProfileKey)
	context.Debugf("Emitting a Put with %s key with all %d reads", elementKey, len(reads))
	context.Debugf("Putting [%s]", reads)
	key, err = datastore.Put(context, elementKey, &models.DayOfReads{reads})
	if err != nil {
		context.Criticalf("Error emitting put with key [%s] for user profile [%s]: %v", elementKey, userProfileKey, err)
		return nil, err
	}

	// Get the time of the batch's last read and update the most recent read timestamp if necessary
	userProfile, err := GetUserProfile(context, userProfileKey)
	if err != nil {
		context.Criticalf("Error reading user profile [%s] for its most recent read value: %v", userProfileKey, err)
		return nil, err
	}

	lastRead := reads[len(reads) - 1]
	if userProfile.MostRecentRead.Before(lastRead.GetTime()) {
		context.Infof("Updating most recent read date to %s", lastRead.GetTime())
		userProfile.MostRecentRead = lastRead.GetTime()
		_, err := StoreUserProfile(context, time.Now(), *userProfile)
		if err != nil {
			context.Criticalf("Error storing updated user profile [%s] with most recent read value of %s: %v", userProfileKey, userProfile.MostRecentRead, err)
			return nil, err
		}
	}

	return key, nil
}

func StoreInjections(context appengine.Context, userProfileKey *datastore.Key, injections []models.Injection) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(injections))
	for i := range injections {
		elementKeys[i] = datastore.NewKey(context, "Injection", "", int64(injections[i].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d injections", len(elementKeys), len(injections))
	key, error := datastore.PutMulti(context, elementKeys, injections)
	if error != nil {
		context.Criticalf("Error writing %d injections with keys [%s]", len(elementKeys), elementKeys)
		return nil, error
	}

	return key, nil
}

func StoreCarbs(context appengine.Context, userProfileKey *datastore.Key, carbs []models.CarbIntake) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(carbs))
	for i := range carbs {
		elementKeys[i] = datastore.NewKey(context, "CarbIntake", "", int64(carbs[i].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d carbs", len(elementKeys), len(carbs))
	key, error := datastore.PutMulti(context, elementKeys, carbs)
	if error != nil {
		context.Criticalf("Error writing %d carbs with keys [%s]", len(elementKeys), elementKeys)
		return nil, error
	}

	return key, nil
}

func StoreExerciseData(context appengine.Context, userProfileKey *datastore.Key, exercises []models.Exercise) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(exercises))
	for i := range exercises {
		elementKeys[i] = datastore.NewKey(context, "Exercise", "", int64(exercises[i].TimeValue), userProfileKey)
	}

	context.Infof("Emitting a PutMulti with %d keys for all %d exercises", len(elementKeys), len(exercises))
	key, error := datastore.PutMulti(context, elementKeys, exercises)
	if error != nil {
		context.Criticalf("Error writing %d exercises with keys [%s]", len(elementKeys), elementKeys)
		return nil, error
	}

	return key, nil
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


func FetchReadsBlob(context appengine.Context, blobKey appengine.BlobKey) (reads []models.MeterRead, err error) {
	reader := blobstore.NewReader(context, blobKey)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&reads); err != nil {
		return nil, err
	}

	return reads, nil
}

func FetchInjectionsBlob(context appengine.Context, blobKey appengine.BlobKey) (injections []models.Injection, err error) {
	reader := blobstore.NewReader(context, blobKey)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&injections); err != nil {
		return nil, err
	}

	return injections, nil
}

func FetchCarbIntakesBlob(context appengine.Context, blobKey appengine.BlobKey) (carbIntakes []models.CarbIntake, err error) {
	reader := blobstore.NewReader(context, blobKey)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&carbIntakes); err != nil {
		return nil, err
	}

	return carbIntakes, nil
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
func GetUserData(context appengine.Context, email string) (userProfile *models.GlukitUser, key *datastore.Key, reads []models.MeterRead, injections []models.Injection, carbIntakes []models.CarbIntake, err error) {
	key = GetUserKey(context, email)
	userProfile, err = GetUserProfile(context, key)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var upperBound time.Time;
	if (userProfile.MostRecentRead.Hour() < 6) {
		// Rewind by one more day
		previousDay := userProfile.MostRecentRead.Add(time.Duration(-24*time.Hour))
		upperBound = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, timeutils.TIMEZONE_LOCATION)
	} else {
		upperBound = time.Date(userProfile.MostRecentRead.Year(), userProfile.MostRecentRead.Month(), userProfile.MostRecentRead.Day(), 6, 0, 0, 0, timeutils.TIMEZONE_LOCATION)
	}
	lowerBound := upperBound.Add(time.Duration(-24*time.Hour))

	// Scan start should be one day prior and scan end should be one day later so that we can capture the day using
	// a single column inequality filter. The scan should actually capture at least one day and a maximum of 3
	scanStart := lowerBound.Add(time.Duration(-24*time.Hour))
	scanEnd := upperBound.Add(time.Duration(24*time.Hour))

	context.Infof("Scanning for reads between %s and %s to get reads between %s and %s", scanStart, scanEnd, lowerBound, upperBound)

	readQuery := datastore.NewQuery("DayOfReads").Ancestor(key).Filter("startTime >=", scanStart).Filter("startTime <=", scanEnd).Order("startTime")
	var dayOfReads models.DayOfReads

	iterator := readQuery.Run(context)
	count := 0
	for _, err := iterator.Next(&dayOfReads); err == nil; _, err = iterator.Next(&dayOfReads) {
		batchSize := len(dayOfReads.Reads) - count
		context.Infof("Loaded batch of %d reads...", batchSize)
		count = len(dayOfReads.Reads)
	}

	filteredReads := datautils.FilterReads(dayOfReads.Reads, lowerBound, upperBound)

	if err != datastore.Done {
		sysutils.Propagate(err)
	}

	carbQuery := datastore.NewQuery("CarbIntake").Ancestor(key).Filter("timestamp >=", lowerBound.Unix()).Filter("timestamp <=", upperBound.Unix()).Order("timestamp")
	_, err = carbQuery.GetAll(context, &carbIntakes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	context.Infof("Loaded %d carbs...", len(carbIntakes))

	injectionQuery := datastore.NewQuery("Injection").Ancestor(key).Filter("timestamp >=", lowerBound.Unix()).Filter("timestamp <=", upperBound.Unix()).Order("timestamp")
	_, err = injectionQuery.GetAll(context, &injections)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	context.Infof("Loaded %d injections...", len(injections))

	//	exercisesQuery := datastore.NewQuery("Exercise").Ancestor(key).Filter("timestamp >", lowerBound.Unix()).Filter("timestamp <", upperBound.Unix()).Order("-timestamp")
	//	_, err = exercisesQuery.GetAll(context, &exercises)
	//	if err != nil {
	//		return nil, nil, nil, nil, nil, err
	//	}
	//	context.Infof("Loaded %d exercises...", len(exercises))

	return userProfile, key, filteredReads, injections, carbIntakes, nil
}
