package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type statsCalibrationReadWriter struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.CalibrationRead
}

func NewStatsCalibrationReadWriter() *statsCalibrationReadWriter {
	w := new(statsCalibrationReadWriter)
	w.batches = make(map[int64][]model.CalibrationRead)

	return w
}

func (w *statsCalibrationReadWriter) WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error) {
	log.Printf("WriteCalibrationBatch with [%d] elements: %v", len(p), p)
	dayOfCalibrations := []model.DayOfCalibrationReads{model.DayOfCalibrationReads{p}}

	return w.WriteCalibrationBatches(dayOfCalibrations)
}

func (w *statsCalibrationReadWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error) {
	log.Printf("WriteCalibrationBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		log.Printf("Persisting batch with start date of [%v]", dayOfData.Reads[0].GetTime())
		w.total += len(dayOfData.Reads)
		w.batches[dayOfData.Reads[0].EpochTime] = dayOfData.Reads
	}

	log.Printf("WriteCalibrationBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsCalibrationReadWriter) Flush() error {
	return nil
}

func TestWriteOfDayCalibrationBatch(t *testing.T) {
	statsWriter := NewStatsCalibrationReadWriter()
	w := NewCalibrationReadStreamerDuration(statsWriter, time.Hour*24)
	calibrations := make([]model.CalibrationRead, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfDayCalibrationBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestWriteOfHourlyCalibrationBatch(t *testing.T) {
	statsWriter := NewStatsCalibrationReadWriter()
	w := NewCalibrationReadStreamerDuration(statsWriter, time.Hour*1)
	calibrations := make([]model.CalibrationRead, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyCalibrationBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}
}

func TestWriteOfMultipleCalibrationBatches(t *testing.T) {
	statsWriter := NewStatsCalibrationReadWriter()
	w := NewCalibrationReadStreamerDuration(statsWriter, time.Hour*1)
	calibrations := make([]model.CalibrationRead, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		calibrations[i] = model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteCalibrations(calibrations)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 25 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 3)
	}

	if statsWriter.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleCalibrationBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 3)
	}
}

func TestCalibrationStreamerWithBufferedIO(t *testing.T) {
	statsWriter := NewStatsCalibrationReadWriter()
	bufferedWriter := bufio.NewCalibrationWriterSize(statsWriter, 2)
	w := NewCalibrationReadStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w.WriteCalibration(model.CalibrationRead{model.Timestamp{"", readTime.Unix()}, b*48 + i})
		}
	}

	w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := statsWriter.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := statsWriter.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := statsWriter.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestGlucoseStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
