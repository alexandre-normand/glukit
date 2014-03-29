package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/util"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleInjectionBatch(t *testing.T) {
	injections := make([]model.Injection, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		injections[i] = model.Injection{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreInjectionBatchWriter(c, key)
	n, err := w.WriteInjectionBatch(injections)
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed, got batch write count of %d but expected %d", n, 1)
	}
}

func TestSimpleWriteOfInjectionBatches(t *testing.T) {
	b := make([]model.DayOfInjections, 10)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 10; i++ {
		injections := make([]model.Injection, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * time.Hour)
			injections[j] = model.Injection{model.Timestamp{readTime.Format(util.TIMEFORMAT_NO_TZ), readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
		}
		b[i] = model.DayOfInjections{injections}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreInjectionBatchWriter(c, key)
	n, err := w.WriteInjectionBatches(b)
	if err != nil {
		t.Fatal(err)
	}

	if n != 10 {
		t.Errorf("TestSimpleWriteOfInjectionBatches failed, got batch write count of %d but expected %d", n, 10)
	}
}
