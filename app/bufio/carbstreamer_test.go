package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestWriteOfDayCarbBatch(t *testing.T) {
	statsWriter := new(statsCarbWriter)
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
	statsWriter := new(statsCarbWriter)
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
	statsWriter := new(statsCarbWriter)
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
