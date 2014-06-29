package bufio_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"log"
	"testing"
)

type exerciseWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Exercise
}

type statsExerciseWriter struct {
	state *exerciseWriterState
}

func NewExerciseWriterState() *exerciseWriterState {
	s := new(exerciseWriterState)
	s.batches = make(map[int64][]apimodel.Exercise)

	return s
}

func NewStatsExerciseWriter(s *exerciseWriterState) *statsExerciseWriter {
	w := new(statsExerciseWriter)
	w.state = s

	return w
}

func (w *statsExerciseWriter) WriteExerciseBatch(p []apimodel.Exercise) (glukitio.ExerciseBatchWriter, error) {
	log.Printf("WriteExerciseBatch with [%d] elements: %v", len(p), p)

	return w.WriteExerciseBatches([]apimodel.DayOfExercises{apimodel.NewDayOfExercises(p)})
}

func (w *statsExerciseWriter) WriteExerciseBatches(p []apimodel.DayOfExercises) (glukitio.ExerciseBatchWriter, error) {
	log.Printf("WriteExerciseBatch with [%d] batches: %v", len(p), p)
	for i := range p {
		dayOfData := p[i]
		w.state.total += len(dayOfData.Exercises)
		w.state.batches[dayOfData.Exercises[0].GetTime().Unix()] = dayOfData.Exercises
	}
	log.Printf("WriteExerciseBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsExerciseWriter) Flush() (glukitio.ExerciseBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleExerciseBatch(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseWriterSize(NewStatsExerciseWriter(state), 10)
	batches := make([]apimodel.DayOfExercises, 10)
	for i := 0; i < 10; i++ {
		exercises := make([]apimodel.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = apimodel.Exercise{apimodel.Time{0, "America/Montreal"}, j, "Light", "details"}
		}
		batches[i] = apimodel.NewDayOfExercises(exercises)
	}
	newWriter, _ := w.WriteExerciseBatches(batches)
	w = newWriter.(*BufferedExerciseBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleExerciseBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualExerciseWrite(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseWriterSize(NewStatsExerciseWriter(state), 10)
	exercises := make([]apimodel.Exercise, 24)
	for j := 0; j < 24; j++ {
		exercises[j] = apimodel.Exercise{apimodel.Time{0, "America/Montreal"}, j, "Light", "details"}
	}
	newWriter, _ := w.WriteExerciseBatch(exercises)
	w = newWriter.(*BufferedExerciseBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualExerciseWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualExerciseWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualExerciseWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneExerciseBatch(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseWriterSize(NewStatsExerciseWriter(state), 10)
	batches := make([]apimodel.DayOfExercises, 11)
	for i := 0; i < 11; i++ {
		exercises := make([]apimodel.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = apimodel.Exercise{apimodel.Time{0, "America/Montreal"}, j, "Light", "details"}
		}
		batches[i] = apimodel.NewDayOfExercises(exercises)
	}
	newWriter, _ := w.WriteExerciseBatches(batches)
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra Exercise to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a total of %d but expected %d", state.total, 264)
	}

	if state.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneExerciseBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteTwoFullExerciseBatches(t *testing.T) {
	state := NewExerciseWriterState()
	w := NewExerciseWriterSize(NewStatsExerciseWriter(state), 10)
	batches := make([]apimodel.DayOfExercises, 20)
	for i := 0; i < 20; i++ {
		exercises := make([]apimodel.Exercise, 24)
		for j := 0; j < 24; j++ {
			exercises[j] = apimodel.Exercise{apimodel.Time{0, "America/Montreal"}, j, "Light", "details"}
		}
		batches[i] = apimodel.NewDayOfExercises(exercises)
	}
	newWriter, _ := w.WriteExerciseBatches(batches)
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 240 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestWriteTwoFullExerciseBatches test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra batch to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedExerciseBatchWriter)

	if state.total != 480 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 20 {
		t.Errorf("TestWriteTwoFullExerciseBatches test: got a batchCount of %d but expected %d", state.batchCount, 20)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteTwoFullExerciseBatches test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}
