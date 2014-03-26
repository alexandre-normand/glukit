package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type statsGlucoseReadWriter struct {
	total      int
	batchCount int
	writeCount int
}

func (w *statsGlucoseReadWriter) WriteGlucoseReadBatch(p []model.GlucoseRead) (n int, err error) {
	log.Printf("WriteGlucoseReadBatch with [%d] elements: %v", len(p), p)
	dayOfGlucoseReads := make([]model.DayOfGlucoseReads, 1)
	dayOfGlucoseReads[0] = model.DayOfGlucoseReads{p}

	return w.WriteGlucoseReadBatches(dayOfGlucoseReads)
}

func (w *statsGlucoseReadWriter) WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (n int, err error) {
	log.Printf("WriteGlucoseReadBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Reads)
	}
	log.Printf("WriteGlucoseReadBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func TestSimpleWriteOfSingleGlucoseReadBatch(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	batches := make([]model.DayOfGlucoseReads, 10)
	for i := 0; i < 10; i++ {
		GlucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			GlucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfGlucoseReads{GlucoseReads}
	}
	w.WriteGlucoseReadBatches(batches)
	w.Flush()

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a batchCount of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestIndividualGlucoseReadWrite(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	GlucoseReads := make([]model.GlucoseRead, 24)
	for j := 0; j < 24; j++ {
		GlucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
	}
	w.WriteGlucoseReadBatch(GlucoseReads)
	w.Flush()

	if statsWriter.total != 24 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a batchCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a writeCount of %d but expected %d", statsWriter.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneGlucoseReadBatch(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	batches := make([]model.DayOfGlucoseReads, 11)
	for i := 0; i < 11; i++ {
		GlucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			GlucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfGlucoseReads{GlucoseReads}
	}
	w.WriteGlucoseReadBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra GlucoseRead to be written
	w.Flush()

	if statsWriter.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a total of %d but expected %d", statsWriter.total, 264)
	}

	if statsWriter.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}

func TestWriteTwoFullGlucoseReadBatches(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	batches := make([]model.DayOfGlucoseReads, 20)
	for i := 0; i < 20; i++ {
		GlucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			GlucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfGlucoseReads{GlucoseReads}
	}
	w.WriteGlucoseReadBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra batch to be written
	w.Flush()

	if statsWriter.total != 480 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 20 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 20)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullGlucoseReadBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}
