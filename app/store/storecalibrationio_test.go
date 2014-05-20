package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/store"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleCalibrationBatch(t *testing.T) {
	r := make([]model.CalibrationRead, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreCalibrationBatchWriter(c, key)
	if _, err = w.WriteCalibrationBatch(r); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWriteOfCalibrationBatches(t *testing.T) {
	b := make([]model.DayOfCalibrationReads, 10)

	for i := 0; i < 10; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(j) * time.Hour)
			calibrations[j] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
		}
		b[i] = model.DayOfCalibrationReads{calibrations}
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	key := GetUserKey(c, "test@glukit.com")

	w := NewDataStoreCalibrationBatchWriter(c, key)
	if _, err = w.WriteCalibrationBatches(b); err != nil {
		t.Fatal(err)
	}
}
