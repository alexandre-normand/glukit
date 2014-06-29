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

type exerciseWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Exercise
}

type statsExerciseReadWriter struct {
	state *exerciseWriterState
}

func NewExerciseWriterState() *exerciseWriterState {
	s := new(exerciseWriterState)
	s.batches = make(map[int64][]apimodel.Exercise)

	return s
}

func NewStatsExerciseReadWriter(s *exerciseWriterState) *statsExerciseReadWriter {
	w := new(statsExerciseReadWriter)
	w.state = s

	return w
}

func (w *statsExerciseReadWriter) WriteExerciseBatch(p []apimodel.Exercise) (glukitio.ExerciseBatchWriter, error) {
	log.Printf("WriteExerciseReadBatch with [%d] elements: %v", len(p), p)
	dayOfExercises := []apimodel.DayOfExercises{apimodel.NewDayOfExercises(p)}

	return w.WriteExerciseBatches(dayOfExercises)
}

func (w *statsExerciseReadWriter) WriteExerciseBatches(p []apimodel.DayOfExercises) (glukitio.ExerciseBatchWriter, error) {
	log.Printf("WriteExerciseBatches with [%d] batches: %v", len(p), p)
	for i := range p {
		dayOfData := p[i]
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Exercises[0].GetTime())
		w.state.total += len(dayOfData.Exercises)
		w.state.batches[dayOfData.Exercises[0].GetTime().Unix()] = dayOfData.Exercises
	}

	log.Printf("WriteExerciseReadBatches with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsExerciseReadWriter) Flush() (glukitio.ExerciseBatchWriter, error) {
	return w, nil
}

func TestWriteOfDayExerciseBatch(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseStreamerDuration(NewStatsExerciseReadWriter(state), time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w, _ = w.WriteExercise(apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestWriteOfDayExerciseBatchesInSingleCall(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseStreamerDuration(NewStatsExerciseReadWriter(state), time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	exercises := make([]apimodel.Exercise, 25)

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		exercises[i] = apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"}
	}

	w, _ = w.WriteExercises(exercises)
	w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfDayExerciseBatchesInSingleCall failed: got a total of %d but expected %d", state.total, 25)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfDayExerciseBatchesInSingleCall failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfDayExerciseBatchesInSingleCall failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfHourlyExerciseBatch(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseStreamerDuration(NewStatsExerciseReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteExercise(apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"})
	}

	if state.total != 12 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a total of %d but expected %d", state.total, 12)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 13 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfMultipleExerciseBatches(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseStreamerDuration(NewStatsExerciseReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteExercise(apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a batchCount of %d but expected %d", state.batchCount, 3)
	}

	if state.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a writeCount of %d but expected %d", state.writeCount, 3)
	}
}

func TestExerciseStreamerWithBufferedIO(t *testing.T) {
	state := NewExerciseWriterState()
	bufferedWriter := bufio.NewExerciseWriterSize(NewStatsExerciseReadWriter(state), 2)
	w := NewExerciseStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteExercise(apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, b*48 + i, "Light", "details"})
		}
	}

	w, _ = w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestExerciseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestExerciseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestExerciseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
