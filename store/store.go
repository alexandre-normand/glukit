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

func StoreUserData(updatedAt time.Time, u *user.User, writer http.ResponseWriter, context appengine.Context, reads []models.MeterRead) (key *datastore.Key, err error) {
	blobKey, err := StoreReadsBlob(context, reads)
	if err != nil {
		utils.Propagate(err)
	}

	//modifiedDate, error := utils.ParseGoogleDriveDate(file.ModifiedDate)
	//log.Printf("Modified date is: %s", modifiedDate.String())
	readData := models.ReadData{Email: u.Email,
		Name: u.String(),
		LastUpdated: updatedAt,
		ReadsBlobKey: blobKey}


	key, error := datastore.Put(context, GetKey(context, u), &readData)
	if error != nil {
		utils.Propagate(error)
	}

	return key, nil
}

func StoreReadsBlob(context appengine.Context, reads []models.MeterRead) (blobKey appengine.BlobKey, err error) {
	var k appengine.BlobKey
	w, err := blobstore.Create(context, "application/json")
	if err != nil {
		return k, err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(reads)
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

func GetKey(context appengine.Context, u *user.User) (key *datastore.Key) {
	return datastore.NewKey(context, "ReadData", u.Email, 0, nil)
}

func GetUserData(context appengine.Context, u *user.User) (readData *models.ReadData, reads []models.MeterRead, err error) {
	readData = new(models.ReadData)
	key := GetKey(context, u)
	log.Printf("Fetching readData for key: %s", key.String())
	error := datastore.Get(context, key, readData)
	if error != nil {
		return nil, nil, error
	}

	reads, err = FetchReadsBlob(context, readData.ReadsBlobKey)
	if err != nil {
		return nil, nil, err
	}

	return readData, reads, nil
}
