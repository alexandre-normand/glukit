package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type statsCalibrationWriter struct {
	total      int
	batchCount int
}

func (w *statsCalibrationWriter) WriteCalibrations(p []model.CalibrationRead) (n int, err error) {
	log.Printf("Got some calibrations")
	w.total += len(p)
	w.batchCount++

	return len(p), nil
}

func TestBufferSimple(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterSize(statsWriter, 10)
	calibrations := make([]model.CalibrationRead, 10)
	for i := 0; i < 10; i++ {
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", 0}, 75}
	}
	w.WriteCalibrations(calibrations)
	w.Flush()

	if statsWriter.total != 10 {
		t.Errorf("simple batch test failed: got a total of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("simple batch test failed: got a total of %d but expected %d", statsWriter.batchCount, 1)
	}
}
