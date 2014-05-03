package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
	"time"
)

type statsGlucoseReadWriter struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.GlucoseRead
}

func NewStatsGlucoseReadWriter() *statsGlucoseReadWriter {
	w := new(statsGlucoseReadWriter)
	w.batches = make(map[int64][]model.GlucoseRead)

	return w
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
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Reads[0].GetTime())
		w.total += len(dayOfData.Reads)
		w.batches[dayOfData.Reads[0].EpochTime] = dayOfData.Reads
	}

	log.Printf("WriteGlucoseReadBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsGlucoseReadWriter) Flush() error {
	return nil
}

func TestSimpleWriteOfSingleGlucoseReadBatch(t *testing.T) {
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	batches := make([]model.DayOfGlucoseReads, 10)
	for i := 0; i < 10; i++ {
		glucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			glucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfGlucoseReads{glucoseReads}
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
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	glucoseReads := make([]model.GlucoseRead, 24)
	for j := 0; j < 24; j++ {
		glucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
	}
	w.WriteGlucoseReadBatch(glucoseReads)
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
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseReadWriterSize(statsWriter, 10)
	batches := make([]model.DayOfGlucoseReads, 19)
	for i := 0; i < 19; i++ {
		glucoseReads := make([]model.GlucoseRead, 24)
		for j := 0; j < 24; j++ {
			glucoseReads[j] = model.GlucoseRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfGlucoseReads{glucoseReads}
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

	if statsWriter.total != 456 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a total of %d but expected %d", statsWriter.total, 456)
	}

	if statsWriter.batchCount != 19 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneGlucoseReadBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}

func TestWriteOverTwoFullGlucoseReadBatches(t *testing.T) {
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseReadWriterSize(statsWriter, 2)
	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		glucoseReads := make([]model.GlucoseRead, 48)

		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			glucoseReads[i] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, b*48 + i}
		}

		w.WriteGlucoseReadBatch(glucoseReads)
	}

	w.Flush()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := statsWriter.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := statsWriter.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := statsWriter.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestWriteOverTwoFullGlucoseReadBatches test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
