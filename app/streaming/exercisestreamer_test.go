package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type statsExerciseWriter struct {
	total      int
	batchCount int
	writeCount int
}

func (w *statsExerciseWriter) WriteExerciseBatch(p []model.Exercise) (n int, err error) {
	log.Printf("WriteExerciseBatch with [%d] elements: %v", len(p), p)
	dayOfExercises := make([]model.DayOfExercises, 1)
	dayOfExercises[0] = model.DayOfExercises{p}

	return w.WriteExerciseBatches(dayOfExercises)
}

func (w *statsExerciseWriter) WriteExerciseBatches(p []model.DayOfExercises) (n int, err error) {
	log.Printf("WriteExerciseBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Exercises)
	}
	log.Printf("WriteExerciseBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsExerciseWriter) Flush() error {
	return nil
}

func TestWriteOfDayExerciseBatch(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseStreamerDuration(statsWriter, time.Hour*24)
	exercises := make([]model.Exercise, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		exercises[i] = model.Exercise{model.Timestamp{"", readTime.Unix()}, 10, "Light"}
	}
	w.WriteExercises(exercises)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfDayExerciseBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestWriteOfHourlyExerciseBatch(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseStreamerDuration(statsWriter, time.Hour*1)
	exercises := make([]model.Exercise, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		exercises[i] = model.Exercise{model.Timestamp{"", readTime.Unix()}, 10, "Light"}
	}
	w.WriteExercises(exercises)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyExerciseBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}
}

func TestWriteOfMultipleExerciseBatches(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseStreamerDuration(statsWriter, time.Hour*1)
	exercises := make([]model.Exercise, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		exercises[i] = model.Exercise{model.Timestamp{"", readTime.Unix()}, 10, "Light"}
	}
	w.WriteExercises(exercises)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 25 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 3)
	}

	if statsWriter.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleExerciseBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 3)
	}
}
