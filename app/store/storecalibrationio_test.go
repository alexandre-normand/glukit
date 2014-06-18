package store_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/apimodel"
	. "github.com/alexandre-normand/glukit/app/store"
	"testing"
	"time"
)

func TestSimpleWriteOfSingleCalibrationBatch(t *testing.T) {
	r := make([]apimodel.CalibrationRead, 25)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.CalibrationRead{apimodel.Time{readTime.Unix(), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
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
	b := make([]apimodel.DayOfCalibrationReads, 10)

	for i := 0; i < 10; i++ {
		calibrations := make([]apimodel.CalibrationRead, 24)
		ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(j) * time.Hour)
			calibrations[j] = apimodel.CalibrationRead{apimodel.Time{readTime.Unix(), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(j)}
		}
		b[i] = apimodel.NewDayOfCalibrationReads(calibrations)
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
