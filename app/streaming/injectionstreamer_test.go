package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type injectionWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Injection
}

type statsInjectionReadWriter struct {
	state *injectionWriterState
}

func NewInjectionWriterState() *injectionWriterState {
	s := new(injectionWriterState)
	s.batches = make(map[int64][]apimodel.Injection)

	return s
}

func NewStatsInjectionReadWriter(s *injectionWriterState) *statsInjectionReadWriter {
	w := new(statsInjectionReadWriter)
	w.state = s

	return w
}

func (w *statsInjectionReadWriter) WriteInjectionBatch(p []apimodel.Injection) (glukitio.InjectionBatchWriter, error) {
	log.Printf("WriteInjectionReadBatch with [%d] elements: %v", len(p), p)
	dayOfInjections := []apimodel.DayOfInjections{apimodel.NewDayOfInjections(p)}

	return w.WriteInjectionBatches(dayOfInjections)
}

func (w *statsInjectionReadWriter) WriteInjectionBatches(p []apimodel.DayOfInjections) (glukitio.InjectionBatchWriter, error) {
	log.Printf("WriteInjectionBatches with [%d] batches: %v", len(p), p)
	for i := range p {
		dayOfData := p[i]
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Injections[0].GetTime())
		w.state.total += len(dayOfData.Injections)
		w.state.batches[dayOfData.Injections[0].GetTime().Unix()] = dayOfData.Injections
	}

	log.Printf("WriteInjectionReadBatches with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsInjectionReadWriter) Flush() (glukitio.InjectionBatchWriter, error) {
	return w, nil
}

func TestWriteOfDayInjectionBatch(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionStreamerDuration(NewStatsInjectionReadWriter(state), apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w, _ = w.WriteInjection(apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfDayInjectionBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestWriteOfDayInjectionBatchesInSingleCall(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionStreamerDuration(NewStatsInjectionReadWriter(state), apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	injections := make([]apimodel.Injection, 25)

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		injections[i] = apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"}
	}

	w, _ = w.WriteInjections(injections)
	w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfDayInjectionBatchesInSingleCall failed: got a total of %d but expected %d", state.total, 25)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfDayInjectionBatchesInSingleCall failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfDayInjectionBatchesInSingleCall failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfHourlyInjectionBatch(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionStreamerDuration(NewStatsInjectionReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteInjection(apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"})
	}

	if state.total != 12 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a total of %d but expected %d", state.total, 12)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 13 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyInjectionBatch failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfMultipleInjectionBatches(t *testing.T) {
	state := NewInjectionWriterState()
	w := NewInjectionStreamerDuration(NewStatsInjectionReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteInjection(apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a batchCount of %d but expected %d", state.batchCount, 3)
	}

	if state.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleInjectionBatches failed: got a writeCount of %d but expected %d", state.writeCount, 3)
	}
}

func TestInjectionStreamerWithBufferedIO(t *testing.T) {
	state := NewInjectionWriterState()
	bufferedWriter := bufio.NewInjectionWriterSize(NewStatsInjectionReadWriter(state), 2)
	w := NewInjectionStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteInjection(apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(b*48 + i), "Humalog", "Bolus"})
		}
	}

	w, _ = w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestInjectionStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
