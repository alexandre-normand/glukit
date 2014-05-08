package streaming_test

import (
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	. "github.com/alexandre-normand/glukit/app/streaming"
	"log"
	"testing"
	"time"
)

type statsCarbWriter struct {
	total      int
	batchCount int
	writeCount int
	batches    map[int64][]model.Carb
}

func NewStatsCarbsWriter() *statsCarbWriter {
	w := new(statsCarbWriter)
	w.batches = make(map[int64][]model.Carb)

	return w
}

func (w *statsCarbWriter) WriteCarbBatch(p []model.Carb) (n int, err error) {
	log.Printf("WriteCarbBatch with [%d] elements: %v", len(p), p)

	return w.WriteCarbBatches([]model.DayOfCarbs{model.DayOfCarbs{p}})
}

func (w *statsCarbWriter) WriteCarbBatches(p []model.DayOfCarbs) (n int, err error) {
	log.Printf("WriteCarbBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Carbs)
		w.batches[dayOfData.Carbs[0].EpochTime] = dayOfData.Carbs
	}

	log.Printf("WriteCarbBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsCarbWriter) Flush() error {
	return nil
}

func TestWriteOfDayCarbBatch(t *testing.T) {
	statsWriter := NewStatsCarbsWriter()
	w := NewCarbStreamerDuration(statsWriter, time.Hour*24)
	carb := make([]model.Carb, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		carb[i] = model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteCarbs(carb)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfDayCarbBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestWriteOfHourlyCarbBatch(t *testing.T) {
	statsWriter := NewStatsCarbsWriter()
	w := NewCarbStreamerDuration(statsWriter, time.Hour*1)
	carb := make([]model.Carb, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		carb[i] = model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteCarbs(carb)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyCarbBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}
}

func TestWriteOfMultipleCarbBatches(t *testing.T) {
	statsWriter := NewStatsCarbsWriter()
	w := NewCarbStreamerDuration(statsWriter, time.Hour*1)
	carb := make([]model.Carb, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		carb[i] = model.Carb{model.Timestamp{"", readTime.Unix()}, float32(1.5), model.UNDEFINED_READ}
	}
	w.WriteCarbs(carb)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 25 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 3)
	}

	if statsWriter.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleCarbBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 3)
	}
}

func TestCarbStreamerWithBufferedIO(t *testing.T) {
	statsWriter := NewStatsCarbsWriter()
	bufferedWriter := bufio.NewCarbWriterSize(statsWriter, 2)
	w := NewCarbStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w.WriteCarb(model.Carb{model.Timestamp{"", readTime.Unix()}, float32(b*48 + i), model.UNDEFINED_READ})
		}
	}

	w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := statsWriter.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := statsWriter.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := statsWriter.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestCarbStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
