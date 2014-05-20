package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type carbWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.Carb
}

type statsCarbWriter struct {
	state *carbWriterState
}

func NewCarbWriterState() *carbWriterState {
	s := new(carbWriterState)
	s.batches = make(map[int64][]model.Carb)

	return s
}

func NewStatsCarbWriter(s *carbWriterState) *statsCarbWriter {
	w := new(statsCarbWriter)
	w.state = s

	return w
}

func (w *statsCarbWriter) WriteCarbBatch(p []model.Carb) (glukitio.CarbBatchWriter, error) {
	log.Printf("WriteCarbBatch with [%d] elements: %v", len(p), p)

	return w.WriteCarbBatches([]model.DayOfCarbs{model.DayOfCarbs{p}})
}

func (w *statsCarbWriter) WriteCarbBatches(p []model.DayOfCarbs) (glukitio.CarbBatchWriter, error) {
	log.Printf("WriteCarbBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.state.total += len(dayOfData.Carbs)
		w.state.batches[dayOfData.Carbs[0].EpochTime] = dayOfData.Carbs
	}
	log.Printf("WriteCarbBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsCarbWriter) Flush() (glukitio.CarbBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleCarbBatch(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbWriterSize(NewStatsCarbWriter(state), 10)
	batches := make([]model.DayOfCarbs, 10)
	for i := 0; i < 10; i++ {
		carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{carbs}
	}
	newWriter, _ := w.WriteCarbBatches(batches)
	w = newWriter.(*BufferedCarbBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualCarbWrite(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbWriterSize(NewStatsCarbWriter(state), 10)
	carbs := make([]model.Carb, 24)
	for j := 0; j < 24; j++ {
		carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
	}
	newWriter, _ := w.WriteCarbBatch(carbs)
	w = newWriter.(*BufferedCarbBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualCarbWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualCarbWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualCarbWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneCarbBatch(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbWriterSize(NewStatsCarbWriter(state), 10)
	batches := make([]model.DayOfCarbs, 11)
	for i := 0; i < 11; i++ {
		carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{carbs}
	}
	newWriter, _ := w.WriteCarbBatches(batches)
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra Carb to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a total of %d but expected %d", state.total, 264)
	}

	if state.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteTwoFullCarbBatches(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbWriterSize(NewStatsCarbWriter(state), 10)
	batches := make([]model.DayOfCarbs, 20)
	for i := 0; i < 20; i++ {
		carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{carbs}
	}
	newWriter, _ := w.WriteCarbBatches(batches)
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 240 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestWriteTwoFullCarbBatches test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra batch to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedCarbBatchWriter)

	if state.total != 480 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 20 {
		t.Errorf("TestWriteTwoFullCarbBatches test: got a batchCount of %d but expected %d", state.batchCount, 20)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}
