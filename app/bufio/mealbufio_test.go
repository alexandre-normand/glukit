package bufio_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"log"
	"testing"
)

type mealWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Meal
}

type statsMealWriter struct {
	state *mealWriterState
}

func NewMealWriterState() *mealWriterState {
	s := new(mealWriterState)
	s.batches = make(map[int64][]apimodel.Meal)

	return s
}

func NewStatsMealWriter(s *mealWriterState) *statsMealWriter {
	w := new(statsMealWriter)
	w.state = s

	return w
}

func (w *statsMealWriter) WriteMealBatch(p []apimodel.Meal) (glukitio.MealBatchWriter, error) {
	log.Printf("WriteMealBatch with [%d] elements: %v", len(p), p)

	return w.WriteMealBatches([]apimodel.DayOfMeals{apimodel.DayOfMeals{p}})
}

func (w *statsMealWriter) WriteMealBatches(p []apimodel.DayOfMeals) (glukitio.MealBatchWriter, error) {
	log.Printf("WriteMealBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.state.total += len(dayOfData.Meals)
		w.state.batches[dayOfData.Meals[0].GetTime().Unix()] = dayOfData.Meals
	}
	log.Printf("WriteMealBatch with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsMealWriter) Flush() (glukitio.MealBatchWriter, error) {
	return w, nil
}

func TestSimpleWriteOfSingleMealBatch(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealWriterSize(NewStatsMealWriter(state), 10)
	batches := make([]apimodel.DayOfMeals, 10)
	for i := 0; i < 10; i++ {
		meals := make([]apimodel.Meal, 24)
		for j := 0; j < 24; j++ {
			meals[j] = apimodel.Meal{apimodel.Time{0, "America/Montreal"}, float32(j), float32(j + 1), float32(j + 2), float32(j + 3)}
		}
		batches[i] = apimodel.DayOfMeals{meals}
	}
	newWriter, _ := w.WriteMealBatches(batches)
	w = newWriter.(*BufferedMealBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleMealBatch failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleMealBatch failed: got a batchCount of %d but expected %d", state.total, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleMealBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestIndividualMealWrite(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealWriterSize(NewStatsMealWriter(state), 10)
	meals := make([]apimodel.Meal, 24)
	for j := 0; j < 24; j++ {
		meals[j] = apimodel.Meal{apimodel.Time{0, "America/Montreal"}, float32(j), float32(j + 1), float32(j + 2), float32(j + 3)}
	}
	newWriter, _ := w.WriteMealBatch(meals)
	w = newWriter.(*BufferedMealBatchWriter)
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 24 {
		t.Errorf("TestIndividualMealWrite failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestIndividualMealWrite failed: got a batchCount of %d but expected %d", state.total, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestIndividualMealWrite failed: got a writeCount of %d but expected %d", state.batchCount, 1)
	}
}

func TestSimpleWriteLargerThanOneMealBatch(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealWriterSize(NewStatsMealWriter(state), 10)
	batches := make([]apimodel.DayOfMeals, 11)
	for i := 0; i < 11; i++ {
		meals := make([]apimodel.Meal, 24)
		for j := 0; j < 24; j++ {
			meals[j] = apimodel.Meal{apimodel.Time{0, "America/Montreal"}, float32(j), float32(j + 1), float32(j + 2), float32(j + 3)}
		}
		batches[i] = apimodel.DayOfMeals{meals}
	}
	newWriter, _ := w.WriteMealBatches(batches)
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra Meal to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test failed: got a total of %d but expected %d", state.total, 264)
	}

	if state.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test: got a batchCount of %d but expected %d", state.batchCount, 11)
	}

	if state.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneMealBatch test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}

func TestWriteTwoFullMealBatches(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealWriterSize(NewStatsMealWriter(state), 10)
	batches := make([]apimodel.DayOfMeals, 20)
	for i := 0; i < 20; i++ {
		meals := make([]apimodel.Meal, 24)
		for j := 0; j < 24; j++ {
			meals[j] = apimodel.Meal{apimodel.Time{0, "America/Montreal"}, float32(j), float32(j + 1), float32(j + 2), float32(j + 3)}
		}
		batches[i] = apimodel.DayOfMeals{meals}
	}
	newWriter, _ := w.WriteMealBatches(batches)
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 240 {
		t.Errorf("TestWriteTwoFullMealBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 10 {
		t.Errorf("TestWriteTwoFullMealBatches test: got a batchCount of %d but expected %d", state.batchCount, 10)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteTwoFullMealBatches test failed: got a writeCount of %d but expected %d", state.total, 1)
	}

	// Flushing should cause the extra batch to be written
	newWriter, _ = w.Flush()
	w = newWriter.(*BufferedMealBatchWriter)

	if state.total != 480 {
		t.Errorf("TestWriteTwoFullMealBatches test failed: got a total of %d but expected %d", state.total, 240)
	}

	if state.batchCount != 20 {
		t.Errorf("TestWriteTwoFullMealBatches test: got a batchCount of %d but expected %d", state.batchCount, 20)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteTwoFullMealBatches test failed: got a writeCount of %d but expected %d", state.total, 2)
	}
}
