package store

import (
	"encoding/json"
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"time"
	"utils"
	"log"
	"models"
	"timeutils"
)

func StoreUserProfile(context appengine.Context, updatedAt time.Time, userProfile models.GlukitUser) (key *datastore.Key, err error) {
	key, error := datastore.Put(context, GetUserKey(context, userProfile.Email), &userProfile)
	if error != nil {
		utils.Propagate(error)
	}

	return key, nil
}

func StoreReads(context appengine.Context, userProfileKey *datastore.Key, reads []models.MeterRead) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(reads))
	for i := range reads {
		//log.Printf("Prepare store for read (%s, %d, %d)", reads[i].LocalTime, reads[i].TimeValue, reads[i].Value)
		elementKeys[i] = datastore.NewKey(context, "MeterRead", "", int64(reads[i].TimeValue), userProfileKey)
	}

	log.Printf("Emitting a PutMulti with %d keys for all %d reads", len(elementKeys), len(reads))
	key, error := datastore.PutMulti(context, elementKeys, reads)
	if error != nil {
		utils.Propagate(error)
	}

	// Get the time of the batch's last read and update the most recent read timestamp if necessary
	userProfile, err := GetUserProfile(context, userProfileKey)
	if err != nil {
		utils.Propagate(err)
	}

	lastRead := reads[len(reads) - 1]
	if userProfile.MostRecentRead.Before(lastRead.GetTime()) {
		log.Printf("Updating most recent read date to %s", lastRead.GetTime())
		userProfile.MostRecentRead = lastRead.GetTime()
		_, err := StoreUserProfile(context, time.Now(), *userProfile)
		if err != nil {
			utils.Propagate(err)
		}
	}

	return key, nil
}

func StoreInjections(context appengine.Context, userProfileKey *datastore.Key, injections []models.Injection) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(injections))
	for i := range injections {
		elementKeys[i] = datastore.NewKey(context, "Injection", "", int64(injections[i].TimeValue), userProfileKey)
	}

	log.Printf("Emitting a PutMulti with %d keys for all %d injections", len(elementKeys), len(injections))
	key, error := datastore.PutMulti(context, elementKeys, injections)
	if error != nil {
		utils.Propagate(error)
	}

	return key, nil
}

func StoreCarbs(context appengine.Context, userProfileKey *datastore.Key, carbs []models.CarbIntake) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(carbs))
	for i := range carbs {
		elementKeys[i] = datastore.NewKey(context, "CarbIntake", "", int64(carbs[i].TimeValue), userProfileKey)
	}

	log.Printf("Emitting a PutMulti with %d keys for all %d carbs", len(elementKeys), len(carbs))
	key, error := datastore.PutMulti(context, elementKeys, carbs)
	if error != nil {
		utils.Propagate(error)
	}

	return key, nil
}

func StoreExerciseData(context appengine.Context, userProfileKey *datastore.Key, exercises []models.Exercise) (key[] *datastore.Key, err error) {
	elementKeys := make([] *datastore.Key, len(exercises))
	for i := range exercises {
		elementKeys[i] = datastore.NewKey(context, "Exercise", "", int64(exercises[i].TimeValue), userProfileKey)
	}

	log.Printf("Emitting a PutMulti with %d keys for all %d exercises", len(elementKeys), len(exercises))
	key, error := datastore.PutMulti(context, elementKeys, exercises)
	if error != nil {
		utils.Propagate(error)
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
	log.Printf("Fetching user profile for key: %s", key.String())
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

	readQuery := datastore.NewQuery("MeterRead").Ancestor(key).Filter("timestamp >", lowerBound.Unix()).Filter("timestamp <", upperBound.Unix()).Order("timestamp")
	_, err = readQuery.GetAll(context, &reads)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	//log.Printf("Loaded %d reads...", len(reads))

	carbQuery := datastore.NewQuery("CarbIntake").Ancestor(key).Filter("timestamp >", lowerBound.Unix()).Filter("timestamp <", upperBound.Unix()).Order("timestamp")
	_, err = carbQuery.GetAll(context, &carbIntakes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	//log.Printf("Loaded %d carbs...", len(carbIntakes))

	injectionQuery := datastore.NewQuery("Injection").Ancestor(key).Filter("timestamp >", lowerBound.Unix()).Filter("timestamp <", upperBound.Unix()).Order("timestamp")
	_, err = injectionQuery.GetAll(context, &injections)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	//log.Printf("Loaded %d injections...", len(injections))

//	exercisesQuery := datastore.NewQuery("Exercise").Ancestor(key).Filter("timestamp >", lowerBound.Unix()).Filter("timestamp <", upperBound.Unix()).Order("-timestamp")
//	_, err = exercisesQuery.GetAll(context, &exercises)
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//	log.Printf("Loaded %d exercises...", len(exercises))

	return userProfile, key, reads, injections, carbIntakes, nil
}
