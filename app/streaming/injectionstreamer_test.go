package streaming_test

import (
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

func (w *statsInjectionWriter) Flush() error {
	return nil
}

func TestWriteOfDayInjectionBatch(t *testing.T) {
	statsWriter := new(statsInjectionWriter)
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
	statsWriter := new(statsInjectionWriter)
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
	statsWriter := new(statsInjectionWriter)
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
