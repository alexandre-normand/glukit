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

type mealWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.Meal
}

type statsMealReadWriter struct {
	state *mealWriterState
}

func NewMealWriterState() *mealWriterState {
	s := new(mealWriterState)
	s.batches = make(map[int64][]apimodel.Meal)

	return s
}

func NewStatsMealReadWriter(s *mealWriterState) *statsMealReadWriter {
	w := new(statsMealReadWriter)
	w.state = s

	return w
}

func (w *statsMealReadWriter) WriteMealBatch(p []apimodel.Meal) (glukitio.MealBatchWriter, error) {
	log.Printf("WriteMealReadBatch with [%d] elements: %v", len(p), p)
	dayOfMeals := []apimodel.DayOfMeals{apimodel.NewDayOfMeals(p)}

	return w.WriteMealBatches(dayOfMeals)
}

func (w *statsMealReadWriter) WriteMealBatches(p []apimodel.DayOfMeals) (glukitio.MealBatchWriter, error) {
	log.Printf("WriteMealBatches with [%d] batches: %v", len(p), p)
	for i := range p {
		dayOfData := p[i]
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Meals[0].GetTime())
		w.state.total += len(dayOfData.Meals)
		w.state.batches[dayOfData.Meals[0].GetTime().Unix()] = dayOfData.Meals
	}

	log.Printf("WriteMealReadBatches with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsMealReadWriter) Flush() (glukitio.MealBatchWriter, error) {
	return w, nil
}

func TestWriteOfDayMealBatch(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealStreamerDuration(NewStatsMealReadWriter(state), apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w, _ = w.WriteMeal(apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfDayMealBatch failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfDayMealBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfDayMealBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestWriteOfDayMealBatchesInSingleCall(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealStreamerDuration(NewStatsMealReadWriter(state), apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	meals := make([]apimodel.Meal, 25)

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		meals[i] = apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)}
	}

	w, _ = w.WriteMeals(meals)
	w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfDayMealBatchesInSingleCall failed: got a total of %d but expected %d", state.total, 25)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfDayMealBatchesInSingleCall failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfDayMealBatchesInSingleCall failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfHourlyMealBatch(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealStreamerDuration(NewStatsMealReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteMeal(apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)})
	}

	if state.total != 12 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a total of %d but expected %d", state.total, 12)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 13 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyMealBatch failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfMultipleMealBatches(t *testing.T) {
	state := NewMealWriterState()
	w := NewMealStreamerDuration(NewStatsMealReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteMeal(apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a batchCount of %d but expected %d", state.batchCount, 3)
	}

	if state.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleMealBatches failed: got a writeCount of %d but expected %d", state.writeCount, 3)
	}
}

func TestMealStreamerWithBufferedIO(t *testing.T) {
	state := NewMealWriterState()
	bufferedWriter := bufio.NewMealWriterSize(NewStatsMealReadWriter(state), 2)
	w := NewMealStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteMeal(apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)})
		}
	}

	w, _ = w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestMealStreamerWithBufferedIO test failed: count not find first batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestMealStreamerWithBufferedIO test failed: count not find second batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestMealStreamerWithBufferedIO test failed: count not find third batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}

func TestMealBatchBoundaries(t *testing.T) {
	state := NewMealWriterState()
	bufferedWriter := bufio.NewMealWriterSize(NewStatsMealReadWriter(state), 2)
	w := NewMealStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 01:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteMeal(apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)})
		}
	}

	w, _ = w.Close()

	// Fist batch still starts with the first read which isn't a day boundary because we're just keeping track of an array of reads and
	// therefore will have the first read potentially not line up with the data
	firstBatchTime, _ := time.Parse("02/01/2006 15:04", "18/04/2014 01:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestMealBatchBoundaries test failed: count not find first batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	// Second batch starts at the truncated day boundary because we have a matching read that starts with it
	secondBatchTime, _ := time.Parse("02/01/2006 15:04", "19/04/2014 00:00")
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestMealBatchBoundaries test failed: count not find second batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	// Third batch starts at the truncated day boundary because we have a matching read that starts with it
	thirdBatchTime, _ := time.Parse("02/01/2006 15:04", "20/04/2014 00:00")
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestMealBatchBoundaries test failed: count not find third batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	// Fourth batch starts at the truncated day boundary because we have a matching read that starts with it
	fourthBatchTime, _ := time.Parse("02/01/2006 15:04", "21/04/2014 00:00")
	if _, ok := state.batches[fourthBatchTime.Unix()]; !ok {
		t.Errorf("TestMealBatchBoundaries test failed: could not find fourth batch starting with a read time of [%v]/ts[%d] in batches: [%v]", fourthBatchTime, fourthBatchTime.Unix(), state.batches)
	}
}
