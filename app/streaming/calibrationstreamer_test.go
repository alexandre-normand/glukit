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

type calibrationWriterState struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]apimodel.CalibrationRead
}

type statsCalibrationReadWriter struct {
	state *calibrationWriterState
}

func NewCalibrationWriterState() *calibrationWriterState {
	s := new(calibrationWriterState)
	s.batches = make(map[int64][]apimodel.CalibrationRead)

	return s
}

func NewStatsCalibrationReadWriter(s *calibrationWriterState) *statsCalibrationReadWriter {
	w := new(statsCalibrationReadWriter)
	w.state = s

	return w
}

func (w *statsCalibrationReadWriter) WriteCalibrationBatch(p []apimodel.CalibrationRead) (glukitio.CalibrationBatchWriter, error) {
	log.Printf("WriteCalibrationReadBatch with [%d] elements: %v", len(p), p)
	dayOfCalibrationReads := []apimodel.DayOfCalibrationReads{apimodel.NewDayOfCalibrationReads(p)}

	return w.WriteCalibrationBatches(dayOfCalibrationReads)
}

func (w *statsCalibrationReadWriter) WriteCalibrationBatches(p []apimodel.DayOfCalibrationReads) (glukitio.CalibrationBatchWriter, error) {
	log.Printf("WriteCalibrationReadBatches with [%d] batches: %v", len(p), p)
	for i := range p {
		dayOfData := p[i]
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Reads[0].GetTime())
		w.state.total += len(dayOfData.Reads)
		w.state.batches[dayOfData.Reads[0].GetTime().Unix()] = dayOfData.Reads
	}

	log.Printf("WriteCalibrationReadBatches with total of %d", w.state.total)
	w.state.batchCount += len(p)
	w.state.writeCount++

	return w, nil
}

func (w *statsCalibrationReadWriter) Flush() (glukitio.CalibrationBatchWriter, error) {
	return w, nil
}

func TestWriteOfDayCalibrationBatch(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationReadStreamerDuration(NewStatsCalibrationReadWriter(state), time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w, _ = w.WriteCalibration(apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, apimodel.MG_PER_DL, float32(i)})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}
}

func TestWriteOfDayCalibrationReadBatchesInSingleCall(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationReadStreamerDuration(NewStatsCalibrationReadWriter(state), time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	reads := make([]apimodel.CalibrationRead, 25)

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		reads[i] = apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, apimodel.MG_PER_DL, float32(i)}
	}

	w, _ = w.WriteCalibrations(reads)
	w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfDayCalibrationReadBatchesInSingleCall failed: got a total of %d but expected %d", state.total, 25)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfDayCalibrationReadBatchesInSingleCall failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfDayCalibrationReadBatchesInSingleCall failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfHourlyCalibrationBatch(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationReadStreamerDuration(NewStatsCalibrationReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteCalibration(apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, apimodel.MG_PER_DL, float32(i)})
	}

	if state.total != 12 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a total of %d but expected %d", state.total, 12)
	}

	if state.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a batchCount of %d but expected %d", state.batchCount, 1)
	}

	if state.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a writeCount of %d but expected %d", state.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 13 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}
}

func TestWriteOfMultipleCalibrationBatches(t *testing.T) {
	state := NewCalibrationWriterState()
	w := NewCalibrationReadStreamerDuration(NewStatsCalibrationReadWriter(state), time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w, _ = w.WriteCalibration(apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, apimodel.MG_PER_DL, float32(i)})
	}

	if state.total != 24 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a total of %d but expected %d", state.total, 24)
	}

	if state.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a batchCount of %d but expected %d", state.batchCount, 2)
	}

	if state.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a writeCount of %d but expected %d", state.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w, _ = w.Flush()

	if state.total != 25 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a total of %d but expected %d", state.total, 13)
	}

	if state.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a batchCount of %d but expected %d", state.batchCount, 3)
	}

	if state.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a writeCount of %d but expected %d", state.writeCount, 3)
	}
}

func TestCalibrationStreamerWithBufferedIO(t *testing.T) {
	state := NewCalibrationWriterState()
	bufferedWriter := bufio.NewCalibrationWriterSize(NewStatsCalibrationReadWriter(state), 2)
	w := NewCalibrationReadStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w, _ = w.WriteCalibration(apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, apimodel.MG_PER_DL, float32(b*48 + i)})

		}
	}

	w, _ = w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := state.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := state.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := state.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), state.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
