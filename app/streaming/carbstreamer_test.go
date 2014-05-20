package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type carbWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.Carb
}

type statsCarbReadWriter struct {
	state *carbWriterState
}

func NewCarbWriterState() *carbWriterState {
	s := new(carbWriterState)
	s.batches = make(map[int64][]model.Carb)

	return s
}

func NewStatsCarbReadWriter(s *carbWriterState) *statsCarbReadWriter {
	w := new(statsCarbReadWriter)
	w.state = s

	return w
}

func (w *statsCarbReadWriter) WriteCarbBatch(p []model.Carb) (glukitio.CarbBatchWriter, error) {
	log.Printf("WriteCarbReadBatch with [%d] elements: %v", len(p), p)
	dayOfCarbs := []model.DayOfCarbs{model.DayOfCarbs{p}}

	return w.WriteCarbBatches(dayOfCarbs)
}

func (w *statsCarbReadWriter) WriteCarbBatches(p []model.DayOfCarbs) (glukitio.CarbBatchWriter, error) {
	log.Printf("WriteCarbBatches with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Carbs[0].GetTime())
		w.state.total += len(dayOfData.Carbs)
		w.state.batches[dayOfData.Carbs[0].EpochTime] = dayOfData.Carbs
	}

	log.Printf("WriteCarbReadBatches with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsCarbReadWriter) Flush() (glukitio.CarbBatchWriter, error) {
	return w, nil
}

func TestWriteOfDayCarbBatch(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbStreamerDuration(NewStatsCarbReadWriter(state), time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w, _ = w.WriteCarb(model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestWriteOfHourlyCarbBatch(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbStreamerDuration(NewStatsCarbReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteCarb(model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ})
	}

	if state.total != 12 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a total of %d but expected %d", state.total, 12)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 13 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfMultipleCarbBatches(t *testing.T) {
	state := NewCarbWriterState()
	w := NewCarbStreamerDuration(NewStatsCarbReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteCarb(model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a batchCount of %d but expected %d", state.batchCount, 3)
	}

	if state.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a writeCount of %d but expected %d", state.writeCount, 3)
	}
}

func TestCarbStreamerWithBufferedIO(t *testing.T) {
	state := NewCarbWriterState()
	bufferedWriter := bufio.NewCarbWriterSize(NewStatsCarbReadWriter(state), 2)
	w := NewCarbStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteCarb(model.Carb{model.Timestamp{"", readTime.Unix()}, float32(b*48 + i), model.UNDEFINED_READ})
		}
	}

	w, _ = w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
