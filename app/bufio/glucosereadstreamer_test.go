package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestWriteOfDayGlucoseReadBatch(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*24)
	glucoseReads := make([]model.GlucoseRead, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		glucoseReads[i] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteGlucoseReads(glucoseReads)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfDayGlucoseReadBatch failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfDayGlucoseReadBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfDayGlucoseReadBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}
}

func TestWriteOfHourlyGlucoseReadBatch(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*1)
	glucoseReads := make([]model.GlucoseRead, 13)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		glucoseReads[i] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteGlucoseReads(glucoseReads)

	if statsWriter.total != 12 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a total of %d but expected %d", statsWriter.total, 12)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 13 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfHourlyGlucoseReadBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}
}

func TestWriteOfMultipleGlucoseReadBatches(t *testing.T) {
	statsWriter := new(statsGlucoseReadWriter)
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*1)
	reads := make([]model.GlucoseRead, 25)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		reads[i] = model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75}
	}
	w.WriteGlucoseReads(reads)

	if statsWriter.total != 24 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 2 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 2)
	}

	// Flushing should trigger the trailing read to be written
	w.Flush()

	if statsWriter.total != 25 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a total of %d but expected %d", statsWriter.total, 13)
	}

	if statsWriter.batchCount != 3 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a batchCount of %d but expected %d", statsWriter.batchCount, 3)
	}

	if statsWriter.writeCount != 3 {
		t.Errorf("TestWriteOfMultipleGlucoseReadBatches failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 3)
	}
}
