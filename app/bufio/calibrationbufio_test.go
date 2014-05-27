package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type calibrationWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.CalibrationRead
}

type statsCalibrationWriter struct {
	state *calibrationWriterState
}

func NewCalibrationWriterState() *calibrationWriterState {
	s := new(calibrationWriterState)
	s.batches = make(map[int64][]model.CalibrationRead)

	return s
}

func NewStatsCalibrationWriter(s *calibrationWriterState) *statsCalibrationWriter {
	w := new(statsCalibrationWriter)
	w.state = s

	return w
}

func (w *statsCalibrationWriter) WriteCalibrationBatch(p []model.CalibrationRead) (glukitio.CalibrationBatchWriter, error) {
	log.Printf("WriteCalibrationBatch with [%d] elements: %v", len(p), p)

	return w.WriteCalibrationBatches([]model.DayOfCalibrationReads{model.DayOfCalibrationReads{p}})
}

func (w *statsCalibrationWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (glukitio.CalibrationBatchWriter, error) {
	log.Printf("WriteCalibrationBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.state.total += len(dayOfData.Reads)
		w.state.batches[dayOfData.Reads[0].GetTime().Unix()] = dayOfData.Reads
	}
	log.Printf("WriteCalibrationBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsCalibrationWriter) Flush() (glukitio.CalibrationBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleCalibrationBatch(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationWriterSize(NewStatsCalibrationWriter(state), 10)
	batches := make([]model.DayOfCalibrationReads, 10)
	for i := 0; i < 10; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Time{0, "America/Montreal"}, model.MG_PER_DL, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	newWriter, _ := w.WriteCalibrationBatches(batches)
	w = newWriter.(*BufferedCalibrationBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleCalibrationBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleCalibrationBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleCalibrationBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualCalibrationWrite(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationWriterSize(NewStatsCalibrationWriter(state), 10)
	calibrations := make([]model.CalibrationRead, 24)
	for j := 0; j < 24; j++ {
		calibrations[j] = model.CalibrationRead{model.Time{0, "America/Montreal"}, model.MG_PER_DL, 75}
	}
	newWriter, _ := w.WriteCalibrationBatch(calibrations)
	w = newWriter.(*BufferedCalibrationBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualCalibrationWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualCalibrationWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualCalibrationWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneCalibrationBatch(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationWriterSize(NewStatsCalibrationWriter(state), 10)
	batches := make([]model.DayOfCalibrationReads, 11)
	for i := 0; i < 11; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Time{0, "America/Montreal"}, model.MG_PER_DL, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	newWriter, _ := w.WriteCalibrationBatches(batches)
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra calibration to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test failed: got a total of %d but expected %d", state.total, 264)
	}

	if state.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneCalibrationBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteTwoFullCalibrationBatches(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationWriterSize(NewStatsCalibrationWriter(state), 10)
	batches := make([]model.DayOfCalibrationReads, 20)
	for i := 0; i < 20; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Time{0, "America/Montreal"}, model.MG_PER_DL, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	newWriter, _ := w.WriteCalibrationBatches(batches)
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 240 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra batch to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCalibrationBatchWriter)

	if state.total != 480 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 20 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test: got a batchCount of %d but expected %d", state.batchCount, 20)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteTwoFullCalibrationBatches test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}
