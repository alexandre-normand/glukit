package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestWriteOfDayBatch(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterDuration(statsWriter, time.Hour*24)
	calibrations := make([]model.CalibrationRead, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayBatch failed: got a total of %d but expected %d", statsWriter.batchCount, 1)
	}
}

func TestWriteOfHourlyBatch(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterDuration(statsWriter, time.Hour*1)
	calibrations := make([]model.CalibrationRead, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyBatch failed: got a total of %d but expected %d", statsWriter.batchCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyBatch failed: got a total of %d but expected %d", statsWriter.batchCount, 2)
	}
}
