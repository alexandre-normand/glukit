package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
	"time"
)

type glucoseWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.GlucoseRead
}

type statsGlucoseReadWriter struct {
	state *glucoseWriterState
}

func NewGlucoseWriterState() *glucoseWriterState {
	s := new(glucoseWriterState)
	s.batches = make(map[int64][]model.GlucoseRead)

	return s
}

func NewStatsGlucoseReadWriter(s *glucoseWriterState) *statsGlucoseReadWriter {
	w := new(statsGlucoseReadWriter)
	w.state = s

	return w
}

func (w *statsGlucoseReadWriter) WriteGlucoseReadBatch(p []model.GlucoseRead) (glukitio.GlucoseReadBatchWriter, error) {
	log.Printf("WriteGlucoseReadBatch with [%d] elements: %v", len(p), p)

	return w.WriteGlucoseReadBatches([]model.DayOfGlucoseReads{model.DayOfGlucoseReads{p}})
}

func (w *statsGlucoseReadWriter) WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (glukitio.GlucoseReadBatchWriter, error) {
	log.Printf("WriteGlucoseReadBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.state.total += len(dayOfData.Reads)
		log.Printf("Adding batch with time [%v]", dayOfData.Reads[0].GetTime())
		w.state.batches[dayOfData.Reads[0].GetTime().Unix()] = dayOfData.Reads
	}

	log.Printf("WriteGlucoseReadBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsGlucoseReadWriter) Flush() (glukitio.GlucoseReadBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleGlucoseReadBatch(t *testing.T) {
	state := NewGlucoseWriterState()
	w := NewGlucoseReadWriterSize(NewStatsGlucoseReadWriter(state), 10)
	batches := make([]model.DayOfGlucoseReads, 10)
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	for i := 0; i < 10; i++ {
		glucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * 1 * time.Hour)
			glucoseReads[j] = model.GlucoseRead{model.Time{model.GetTimeMillis(readTime), "America/Montreal"}, model.MG_PER_DL, float32(j)}
		}
		batches[i] = model.DayOfGlucoseReads{glucoseReads}
	}
	newWriter, _ := w.WriteGlucoseReadBatches(batches)
	w = newWriter.(*BufferedGlucoseReadBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedGlucoseReadBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleGlucoseReadBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualGlucoseReadWrite(t *testing.T) {
	state := NewGlucoseWriterState()
	w := NewGlucoseReadWriterSize(NewStatsGlucoseReadWriter(state), 10)
	glucoseReads := make([]model.GlucoseRead, 24)
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	for j := 0; j < 24; j++ {
		readTime := ct.Add(time.Duration(j) * 1 * time.Hour)
		glucoseReads[j] = model.GlucoseRead{model.Time{model.GetTimeMillis(readTime), "America/Montreal"}, model.MG_PER_DL, float32(j)}
	}
	newWriter, _ := w.WriteGlucoseReadBatch(glucoseReads)
	w = newWriter.(*BufferedGlucoseReadBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedGlucoseReadBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualGlucoseReadWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneGlucoseReadBatch(t *testing.T) {
	state := NewGlucoseWriterState()
	w := NewGlucoseReadWriterSize(NewStatsGlucoseReadWriter(state), 10)
	batches := make([]model.DayOfGlucoseReads, 19)
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	for i := 0; i < 19; i++ {
		glucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			readTime := ct.Add(time.Duration(i*24+j) * 1 * time.Hour)
			glucoseReads[j] = model.GlucoseRead{model.Time{model.GetTimeMillis(readTime), "America/Montreal"}, model.MG_PER_DL, float32(i*24 + j)}
		}
		batches[i] = model.DayOfGlucoseReads{glucoseReads}
	}
	newWriter, _ := w.WriteGlucoseReadBatches(batches)
	w = newWriter.(*BufferedGlucoseReadBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra GlucoseRead to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedGlucoseReadBatchWriter)

	if state.total != 456 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a total of %d but expected %d", state.total, 456)
	}

	if state.batchCount != 19 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteOverTwoFullGlucoseReadBatches(t *testing.T) {
	state := NewGlucoseWriterState()
	w := NewGlucoseReadWriterSize(NewStatsGlucoseReadWriter(state), 2)
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		glucoseReads := make([]model.GlucoseRead, 48)

		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			glucoseReads[i] = model.GlucoseRead{model.Time{model.GetTimeMillis(readTime), "America/Montreal"}, model.MG_PER_DL, float32(b*48 + i)}
		}

		newWriter, _ := w.WriteGlucoseReadBatch(glucoseReads)
		w = newWriter.(*BufferedGlucoseReadBatchWriter)
	}

	newWriter, _ := w.Flush()
	w = newWriter.(*BufferedGlucoseReadBatchWriter)

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
