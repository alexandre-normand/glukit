package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type statsCarbWriter struct {
	total      int
	batchCount int
	writeCount int
}

func (w *statsCarbWriter) WriteCarbBatch(p []model.Carb) (n int, err error) {
	log.Printf("WriteCarbBatch with [%d] elements: %v", len(p), p)
	dayOfCarbs := make([]model.DayOfCarbs, 1)
	dayOfCarbs[0] = model.DayOfCarbs{p}

	return w.WriteCarbBatches(dayOfCarbs)
}

func (w *statsCarbWriter) WriteCarbBatches(p []model.DayOfCarbs) (n int, err error) {
	log.Printf("WriteCarbBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Carbs)
	}
	log.Printf("WriteCarbBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func TestSimpleWriteOfSingleCarbBatch(t *testing.T) {
	statsWriter := new(statsCarbWriter)
	w := NewCarbWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCarbs, 10)
	for i := 0; i < 10; i++ {
		Carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			Carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{Carbs}
	}
	w.WriteCarbBatches(batches)
	w.Flush()

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a batchCount of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleCarbBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestIndividualCarbWrite(t *testing.T) {
	statsWriter := new(statsCarbWriter)
	w := NewCarbWriterSize(statsWriter, 10)
	Carbs := make([]model.Carb, 24)
	for j := 0; j < 24; j++ {
		Carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteCarbBatch(Carbs)
	w.Flush()

	if statsWriter.total != 24 {
		t.Errorf("TestIndividualCarbWrite failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestIndividualCarbWrite failed: got a batchCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestIndividualCarbWrite failed: got a writeCount of %d but expected %d", statsWriter.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneCarbBatch(t *testing.T) {
	statsWriter := new(statsCarbWriter)
	w := NewCarbWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCarbs, 11)
	for i := 0; i < 11; i++ {
		Carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			Carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{Carbs}
	}
	w.WriteCarbBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra Carb to be written
	w.Flush()

	if statsWriter.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a total of %d but expected %d", statsWriter.total, 264)
	}

	if statsWriter.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneCarbBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}

func TestWriteTwoFullCarbBatches(t *testing.T) {
	statsWriter := new(statsCarbWriter)
	w := NewCarbWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCarbs, 20)
	for i := 0; i < 20; i++ {
		Carbs := make([]model.Carb, 24)
		for j := 0; j < 24; j++ {
			Carbs[j] = model.Carb{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfCarbs{Carbs}
	}
	w.WriteCarbBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestWriteTwoFullCarbBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra batch to be written
	w.Flush()

	if statsWriter.total != 480 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 20 {
		t.Errorf("TestWriteTwoFullCarbBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 20)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullCarbBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}
