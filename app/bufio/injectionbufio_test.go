package bufio_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"log"
	"testing"
)

type injectionWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Injection
}

type statsInjectionWriter struct {
	state *injectionWriterState
}

func NewInjectionWriterState() *injectionWriterState {
	s := new(injectionWriterState)
	s.batches = make(map[int64][]apimodel.Injection)

	return s
}

func NewStatsInjectionWriter(s *injectionWriterState) *statsInjectionWriter {
	w := new(statsInjectionWriter)
	w.state = s

	return w
}

func (w *statsInjectionWriter) WriteInjectionBatch(p []apimodel.Injection) (glukitio.InjectionBatchWriter, error) {
	log.Printf("WriteInjectionBatch with [%d] elements: %v", len(p), p)

	return w.WriteInjectionBatches([]apimodel.DayOfInjections{apimodel.DayOfInjections{p}})
}

func (w *statsInjectionWriter) WriteInjectionBatches(p []apimodel.DayOfInjections) (glukitio.InjectionBatchWriter, error) {
	log.Printf("WriteInjectionBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.state.total += len(dayOfData.Injections)
		w.state.batches[dayOfData.Injections[0].GetTime().Unix()] = dayOfData.Injections
	}
	log.Printf("WriteInjectionBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsInjectionWriter) Flush() (glukitio.InjectionBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleInjectionBatch(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionWriterSize(NewStatsInjectionWriter(state), 10)
	batches := make([]apimodel.DayOfInjections, 10)
	for i := 0; i < 10; i++ {
		injections := make([]apimodel.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = apimodel.Injection{apimodel.Time{0, "America/Montreal"}, float32(j), "Humalog", "Bolus"}
		}
		batches[i] = apimodel.DayOfInjections{injections}
	}
	newWriter, _ := w.WriteInjectionBatches(batches)
	w = newWriter.(*BufferedInjectionBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleInjectionBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualInjectionWrite(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionWriterSize(NewStatsInjectionWriter(state), 10)
	injections := make([]apimodel.Injection, 24)
	for j := 0; j < 24; j++ {
		injections[j] = apimodel.Injection{apimodel.Time{0, "America/Montreal"}, float32(j), "Humalog", "Bolus"}
	}
	newWriter, _ := w.WriteInjectionBatch(injections)
	w = newWriter.(*BufferedInjectionBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualInjectionWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualInjectionWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualInjectionWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneInjectionBatch(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionWriterSize(NewStatsInjectionWriter(state), 10)
	batches := make([]apimodel.DayOfInjections, 11)
	for i := 0; i < 11; i++ {
		injections := make([]apimodel.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = apimodel.Injection{apimodel.Time{0, "America/Montreal"}, float32(j), "Humalog", "Bolus"}
		}
		batches[i] = apimodel.DayOfInjections{injections}
	}
	newWriter, _ := w.WriteInjectionBatches(batches)
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra Injection to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a total of %d but expected %d", state.total, 264)
	}

	if state.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneInjectionBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteTwoFullInjectionBatches(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionWriterSize(NewStatsInjectionWriter(state), 10)
	batches := make([]apimodel.DayOfInjections, 20)
	for i := 0; i < 20; i++ {
		injections := make([]apimodel.Injection, 24)
		for j := 0; j < 24; j++ {
			injections[j] = apimodel.Injection{apimodel.Time{0, "America/Montreal"}, float32(j), "Humalog", "Bolus"}
		}
		batches[i] = apimodel.DayOfInjections{injections}
	}
	newWriter, _ := w.WriteInjectionBatches(batches)
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 240 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestWriteTwoFullInjectionBatches test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra batch to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedInjectionBatchWriter)

	if state.total != 480 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 20 {
		t.Errorf("TestWriteTwoFullInjectionBatches test: got a batchCount of %d but expected %d", state.batchCount, 20)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteTwoFullInjectionBatches test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}
