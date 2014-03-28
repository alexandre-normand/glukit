package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
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

func TestSimpleWriteOfSingleExerciseBatch(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseWriterSize(statsWriter, 10)
	batches := make([]model.DayOfExercises, 10)
	for i := 0; i < 10; i++ {
		exercises := make([]model.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = model.Exercise{model.Timestamp{"", 0}, 10, "Light"}
		}
		batches[i] = model.DayOfExercises{exercises}
	}
	w.WriteExerciseBatches(batches)
	w.Flush()

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a batchCount of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestIndividualExerciseWrite(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseWriterSize(statsWriter, 10)
	exercises := make([]model.Exercise, 24)
	for j := 0; j < 24; j++ {
		exercises[j] = model.Exercise{model.Timestamp{"", 0}, 10, "Light"}
	}
	w.WriteExerciseBatch(exercises)
	w.Flush()

	if statsWriter.total != 24 {
		t.Errorf("TestIndividualExerciseWrite failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestIndividualExerciseWrite failed: got a batchCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestIndividualExerciseWrite failed: got a writeCount of %d but expected %d", statsWriter.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneExerciseBatch(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseWriterSize(statsWriter, 10)
	batches := make([]model.DayOfExercises, 11)
	for i := 0; i < 11; i++ {
		exercises := make([]model.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = model.Exercise{model.Timestamp{"", 0}, 10, "Light"}
		}
		batches[i] = model.DayOfExercises{exercises}
	}
	w.WriteExerciseBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra Exercise to be written
	w.Flush()

	if statsWriter.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a total of %d but expected %d", statsWriter.total, 264)
	}

	if statsWriter.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}

func TestWriteTwoFullExerciseBatches(t *testing.T) {
	statsWriter := new(statsExerciseWriter)
	w := NewExerciseWriterSize(statsWriter, 10)
	batches := make([]model.DayOfExercises, 20)
	for i := 0; i < 20; i++ {
		exercises := make([]model.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = model.Exercise{model.Timestamp{"", 0}, 10, "Light"}
		}
		batches[i] = model.DayOfExercises{exercises}
	}
	w.WriteExerciseBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestWriteTwoFullExerciseBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra batch to be written
	w.Flush()

	if statsWriter.total != 480 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 20 {
		t.Errorf("TestWriteTwoFullExerciseBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 20)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}
