package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/goauth2/oauth"
	"testing"
	"time"
)

var TEST_USER = "test@glukit.com"

func setup(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var oauthToken oauth.Token
	user := model.GlukitUser{TEST_USER, "", "", time.Now(),
		"", "", util.GLUKIT_EPOCH_TIME, model.UNDEFINED_GLUCOSE_READ, oauthToken, oauthToken.RefreshToken,
		model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, "", time.Now()}

	_, err = StoreUserProfile(c, time.Unix(1000, 0), user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfSingleGlucoseReadBatch(t *testing.T) {
	setup(t)

	r := make([]model.GlucoseRead, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, TEST_USER)

	w := NewDataStoreGlucoseReadBatchWriter(c, key)
	n, err := w.WriteGlucoseReadBatch(r)
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed, got batch write count of %d but expected %d", n, 1)
	}
}

func TestSimpleWriteOfGlucoseReadBatches(t *testing.T) {
	setup(t)

	b := make([]model.DayOfGlucoseReads, 10)

	for i := 0; i < 10; i++ {
		calibrations := make([]model.GlucoseRead, 24)
		ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(j) * time.Hour)
			calibrations[j] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75}
		}
		b[i] = model.DayOfGlucoseReads{calibrations}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, TEST_USER)

	w := NewDataStoreGlucoseReadBatchWriter(c, key)
	n, err := w.WriteGlucoseReadBatches(b)
	if err != nil {
		t.Fatal(err)
	}

	if n != 10 {
		t.Errorf("TestSimpleWriteOfGlucoseReadBatches failed, got batch write count of %d but expected %d", n, 10)
	}
}
