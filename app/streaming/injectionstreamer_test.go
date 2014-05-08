package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type statsInjectionWriter struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.Injection
}

func NewStatsInjectionsWriter() *statsInjectionWriter {
	w := new(statsInjectionWriter)
	w.batches = make(map[int64][]model.Injection)

	return w
}

func (w *statsInjectionWriter) WriteInjectionBatch(p []model.Injection) (n int, err error) {
	log.Printf("WriteInjectionBatch with [%d] elements: %v", len(p), p)

	return w.WriteInjectionBatches([]model.DayOfInjections{model.DayOfInjections{p}})
}

func (w *statsInjectionWriter) WriteInjectionBatches(p []model.DayOfInjections) (n int, err error) {
	log.Printf("WriteInjectionBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Injections)
		w.batches[dayOfData.Injections[0].EpochTime] = dayOfData.Injections
	}

	log.Printf("WriteInjectionBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsInjectionWriter) Flush() error {
	return nil
}

func TestWriteOfDayInjectionBatch(t *testing.T) {
	statsWriter := NewStatsInjectionsWriter()
	w := NewInjectionStreamerDuration(statsWriter, time.Hour*24)
	injections := make([]model.Injection, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		injections[i] = model.Injection{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteInjections(injections)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestWriteOfHourlyInjectionBatch(t *testing.T) {
	statsWriter := NewStatsInjectionsWriter()
	w := NewInjectionStreamerDuration(statsWriter, time.Hour*1)
	injections := make([]model.Injection, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		injections[i] = model.Injection{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteInjections(injections)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}
}

func TestWriteOfMultipleInjectionBatches(t *testing.T) {
	statsWriter := NewStatsInjectionsWriter()
	w := NewInjectionStreamerDuration(statsWriter, time.Hour*1)
	injections := make([]model.Injection, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		injections[i] = model.Injection{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteInjections(injections)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 25 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 3)
	}

	if statsWriter.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 3)
	}
}

func TestInjectionStreamerWithBufferedIO(t *testing.T) {
	statsWriter := NewStatsInjectionsWriter()
	bufferedWriter := bufio.NewInjectionWriterSize(statsWriter, 2)
	w := NewInjectionStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w.WriteInjection(model.Injection{model.Timestamp{"", readTime.Unix()}, float32(b*48 + i), model.UNDEFINED_READ})
		}
	}

	w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := statsWriter.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := statsWriter.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := statsWriter.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
