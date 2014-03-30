package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleExerciseBatch(t *testing.T) {
	exercises := make([]model.Exercise, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		exercises[i] = model.Exercise{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, 15, "Light"}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreExerciseBatchWriter(c, key)
	n, err := w.WriteExerciseBatch(exercises)
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed, got batch write count of %d but expected %d", n, 1)
	}
}

func TestSimpleWriteOfExerciseBatches(t *testing.T) {
	b := make([]model.DayOfExercises, 10)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 10; i++ {
		exercises := make([]model.Exercise, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * time.Hour)
			exercises[j] = model.Exercise{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, 15, "Light"}
		}
		b[i] = model.DayOfExercises{exercises}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreExerciseBatchWriter(c, key)
	n, err := w.WriteExerciseBatches(b)
	if err != nil {
		t.Fatal(err)
	}

	if n != 10 {
		t.Errorf("TestSimpleWriteOfExerciseBatches failed, got batch write count of %d but expected %d", n, 10)
	}
}
