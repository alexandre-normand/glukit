package store_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"log"
	"testing"
	"time"
)

var TEST_USER = "test@glukit.com"

func setup(t *testing.T) (c aetest.Context, key *datastore.Key) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	user := model.GlukitUser{TEST_USER, "", "", time.Now(),
		"", "", util.GLUKIT_EPOCH_TIME, apimodel.UNDEFINED_GLUCOSE_READ,
		model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, "", time.Now()}

	key, err = StoreUserProfile(c, time.Unix(1000, 0), user)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Initialized [%s] with key [%v]", TEST_USER, key)

	return c, key
}

func TestSimpleWriteOfSingleGlucoseReadBatch(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	r := make([]apimodel.GlucoseRead, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
	}

	w := NewDataStoreGlucoseReadBatchWriter(c, key)
	if _, err := w.WriteGlucoseReadBatch(r); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfGlucoseReadBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	b := make([]apimodel.DayOfGlucoseReads, 10)

	for i := 0; i < 10; i++ {
		reads := make([]apimodel.GlucoseRead, 24)
		ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(j) * time.Hour)
			reads[j] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i*24 + j)}
		}
		b[i] = apimodel.NewDayOfGlucoseReads(reads)
	}

	w := NewDataStoreGlucoseReadBatchWriter(c, key)
	if _, err := w.WriteGlucoseReadBatches(b); err != nil {
		t.Fatal(err)
	}
}
