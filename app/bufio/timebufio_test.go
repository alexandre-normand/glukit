package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestWriteOfTimedBatch(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterDuration(statsWriter, time.Hour*24)
	calibrations := make([]model.CalibrationRead, 25)

	ct, _ := time.Parse("18/04/2014 00:00", "02/01/2006 15:04")
	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		//log.Printf("time is %v", readTime)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 24 {
		t.Errorf("simple batch test failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("simple batch test failed: got a total of %d but expected %d", statsWriter.batchCount, 1)
	}
}
