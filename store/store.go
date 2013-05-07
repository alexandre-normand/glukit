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
		elementKeys[i] = datastore.NewKey(context, "MeterRead", "", int64(reads[i].TimeValue), userProfileKey)
	}

	log.Printf("Emitting a PutMulti with %d keys for all %d reads", len(elementKeys), len(reads))
	key, error := datastore.PutMulti(context, elementKeys, reads)
	if error != nil {
		utils.Propagate(error)
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

func GetUserData(context appengine.Context, email string) (userProfile *models.GlukitUser, key *datastore.Key, reads []models.MeterRead, injections []models.Injection, carbIntakes []models.CarbIntake, err error) {
	userProfile = new(models.GlukitUser)
	key = GetUserKey(context, email)
	log.Printf("Fetching readData for key: %s", key.String())
	error := datastore.Get(context, key, userProfile)
	if error != nil {
		return nil, nil, nil, nil, nil, error
	}

	reads = make([]models.MeterRead, 288)
	readQuery := datastore.NewQuery("MeterRead").Ancestor(key).Filter("timestamp >", 0).Order("-timestamp")
	_, err = readQuery.GetAll(context, &reads)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	carbIntakes = make([]models.CarbIntake, 288)
	carbQuery := datastore.NewQuery("CarbIntake").Ancestor(key).Filter("timestamp >", 0).Order("-timestamp")
	_, err = carbQuery.GetAll(context, &carbIntakes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	injections = make([]models.Injection, 288)
	injectionQuery := datastore.NewQuery("Injection").Ancestor(key).Filter("timestamp >", 0).Order("-timestamp")
	_, err = injectionQuery.GetAll(context, &injections)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return userProfile, key, reads, injections, carbIntakes, nil
}
