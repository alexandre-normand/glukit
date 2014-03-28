package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type statsInjectionWriter struct {
	total      int
	batchCount int
	writeCount int
}

func (w *statsInjectionWriter) WriteInjectionBatch(p []model.Injection) (n int, err error) {
	log.Printf("WriteInjectionBatch with [%d] elements: %v", len(p), p)
	dayOfInjections := make([]model.DayOfInjections, 1)
	dayOfInjections[0] = model.DayOfInjections{p}

	return w.WriteInjectionBatches(dayOfInjections)
}

func (w *statsInjectionWriter) WriteInjectionBatches(p []model.DayOfInjections) (n int, err error) {
	log.Printf("WriteInjectionBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Injections)
	}
	log.Printf("WriteInjectionBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func TestSimpleWriteOfSingleInjectionBatch(t *testing.T) {
	statsWriter := new(statsInjectionWriter)
	w := NewInjectionWriterSize(statsWriter, 10)
	batches := make([]model.DayOfInjections, 10)
	for i := 0; i < 10; i++ {
		injections := make([]model.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = model.Injection{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfInjections{injections}
	}
	w.WriteInjectionBatches(batches)
	w.Flush()

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a batchCount of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestIndividualInjectionWrite(t *testing.T) {
	statsWriter := new(statsInjectionWriter)
	w := NewInjectionWriterSize(statsWriter, 10)
	injections := make([]model.Injection, 24)
	for j := 0; j < 24; j++ {
		injections[j] = model.Injection{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteInjectionBatch(injections)
	w.Flush()

	if statsWriter.total != 24 {
		t.Errorf("TestIndividualInjectionWrite failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestIndividualInjectionWrite failed: got a batchCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestIndividualInjectionWrite failed: got a writeCount of %d but expected %d", statsWriter.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneInjectionBatch(t *testing.T) {
	statsWriter := new(statsInjectionWriter)
	w := NewInjectionWriterSize(statsWriter, 10)
	batches := make([]model.DayOfInjections, 11)
	for i := 0; i < 11; i++ {
		injections := make([]model.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = model.Injection{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfInjections{injections}
	}
	w.WriteInjectionBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra Injection to be written
	w.Flush()

	if statsWriter.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a total of %d but expected %d", statsWriter.total, 264)
	}

	if statsWriter.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}

func TestWriteTwoFullInjectionBatches(t *testing.T) {
	statsWriter := new(statsInjectionWriter)
	w := NewInjectionWriterSize(statsWriter, 10)
	batches := make([]model.DayOfInjections, 20)
	for i := 0; i < 20; i++ {
		injections := make([]model.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = model.Injection{model.Timestamp{"", 0}, float32(1.5), model.UNDEFINED_READ}
		}
		batches[i] = model.DayOfInjections{injections}
	}
	w.WriteInjectionBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestWriteTwoFullInjectionBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra batch to be written
	w.Flush()

	if statsWriter.total != 480 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 20 {
		t.Errorf("TestWriteTwoFullInjectionBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 20)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}
