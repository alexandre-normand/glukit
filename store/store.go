package store

import (
	"encoding/json"
	"appengine"
	"appengine/user"
	"appengine/blobstore"
	"appengine/datastore"
	"time"
	"net/http"
	"utils"
	"log"
	"models"
)

func StoreUserData(updatedAt time.Time, u *user.User, writer http.ResponseWriter, context appengine.Context, reads []models.MeterRead, injections []models.Injection, carbIntakes []models.CarbIntake) (key *datastore.Key, err error) {
	readsBlobKey, err := StoreReadsBlob(context, reads)
	if err != nil {
		utils.Propagate(err)
	}

	injectionsBlobKey, err := StoreReadsBlob(context, injections)
	if err != nil {
		utils.Propagate(err)
	}

	carbIntakesBlobKey, err := StoreReadsBlob(context, carbIntakes)
	if err != nil {
		utils.Propagate(err)
	}

	//modifiedDate, error := utils.ParseGoogleDriveDate(file.ModifiedDate)
	//log.Printf("Modified date is: %s", modifiedDate.String())
	readData := models.ReadData{Email: u.Email,
		Name: u.String(),
		LastUpdated: updatedAt,
		ReadsBlobKey: readsBlobKey,
		InjectionsBlobKey: injectionsBlobKey,
		CarbIntakesBlobKey: carbIntakesBlobKey}


	key, error := datastore.Put(context, GetKey(context, u), &readData)
	if error != nil {
		utils.Propagate(error)
	}

	return key, nil
}

func StoreReadsBlob(context appengine.Context, data interface{}) (blobKey appengine.BlobKey, err error) {
	var k appengine.BlobKey
	w, err := blobstore.Create(context, "application/json")
	if err != nil {
		return k, err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		return k, err
	}
	err = w.Close()
	if err != nil {
		return k, err
	}

	return w.Key()
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

func GetKey(context appengine.Context, u *user.User) (key *datastore.Key) {
	return datastore.NewKey(context, "ReadData", u.Email, 0, nil)
}

func GetUserData(context appengine.Context, u *user.User) (readData *models.ReadData, reads []models.MeterRead, injections []models.Injection, carbIntakes []models.CarbIntake, err error) {
	readData = new(models.ReadData)
	key := GetKey(context, u)
	log.Printf("Fetching readData for key: %s", key.String())
	error := datastore.Get(context, key, readData)
	if error != nil {
		return nil, nil, nil, nil, error
	}

	reads, err = FetchReadsBlob(context, readData.ReadsBlobKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	injections, err = FetchInjectionsBlob(context, readData.InjectionsBlobKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	carbIntakes, err = FetchCarbIntakesBlob(context, readData.CarbIntakesBlobKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return readData, reads, injections, carbIntakes, nil
}
