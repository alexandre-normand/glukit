package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleCarbBatch(t *testing.T) {
	carbs := make([]model.Carb, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		carbs[i] = model.Carb{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreCarbBatchWriter(c, key)
	if _, err = w.WriteCarbBatch(carbs); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfCarbBatches(t *testing.T) {
	b := make([]model.DayOfCarbs, 10)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 10; i++ {
		carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * time.Hour)
			carbs[j] = model.Carb{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
		}
		b[i] = model.DayOfCarbs{carbs}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreCarbBatchWriter(c, key)
	if _, err = w.WriteCarbBatches(b); err != nil {
		t.Fatal(err)
	}
}
