package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleInjectionBatch(t *testing.T) {
	injections := make([]model.Injection, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		injections[i] = model.Injection{model.Time{model.GetTimeMillis(readTime), "America/Los_Angeles"}, float32(i), "Levemir", "Basal"}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreInjectionBatchWriter(c, key)
	if _, err = w.WriteInjectionBatch(injections); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfInjectionBatches(t *testing.T) {
	b := make([]model.DayOfInjections, 10)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 10; i++ {
		injections := make([]model.Injection, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * time.Hour)
			injections[j] = model.Injection{model.Time{model.GetTimeMillis(readTime), "America/Los_Angeles"}, float32(1.5), "Levemir", "Basal"}
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
	if _, err = w.WriteInjectionBatches(b); err != nil {
		t.Fatal(err)
	}
}
