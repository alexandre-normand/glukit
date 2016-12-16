package store_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	. "github.com/alexandre-normand/glukit/app/store"
	"google.golang.org/appengine/aetest"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleExerciseBatch(t *testing.T) {
	exercises := make([]apimodel.Exercise, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		exercises[i] = apimodel.Exercise{apimodel.Time{readTime.Unix(), "America/Los_Angeles"}, i, "Light", "details"}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreExerciseBatchWriter(c, key)
	if _, err = w.WriteExerciseBatch(exercises); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfExerciseBatches(t *testing.T) {
	b := make([]apimodel.DayOfExercises, 10)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 10; i++ {
		exercises := make([]apimodel.Exercise, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * time.Hour)
			exercises[j] = apimodel.Exercise{apimodel.Time{readTime.Unix(), "America/Los_Angeles"}, j, "Light", "details"}
		}
		b[i] = apimodel.NewDayOfExercises(exercises)
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreExerciseBatchWriter(c, key)
	if _, err = w.WriteExerciseBatches(b); err != nil {
		t.Fatal(err)
	}
}
