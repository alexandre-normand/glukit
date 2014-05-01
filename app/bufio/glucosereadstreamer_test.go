package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"testing"
	"time"
)

func TestWriteOfDayGlucoseReadBatch(t *testing.T) {
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		w.WriteGlucoseRead(model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75})
	}

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
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 13; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w.WriteGlucoseRead(model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75})
	}

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
	statsWriter := NewStatsGlucoseReadWriter()
	w := NewGlucoseStreamerDuration(statsWriter, time.Hour*1)

	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")

	for i := 0; i < 25; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		w.WriteGlucoseRead(model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, 75})
	}

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

func TestStreamerWithBufferedIO(t *testing.T) {
	statsWriter := NewStatsGlucoseReadWriter()
	bufferedWriter := NewGlucoseReadWriterSize(statsWriter, 2)
	w := NewGlucoseStreamerDuration(bufferedWriter, time.Hour*24)

	ct, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")

	for b := 0; b < 3; b++ {
		for i := 0; i < 48; i++ {
			readTime := ct.Add(time.Duration(b*48+i) * 30 * time.Minute)
			w.WriteGlucoseRead(model.GlucoseRead{model.Timestamp{"", readTime.Unix()}, b*48 + i})
		}
	}

	w.Close()

	firstBatchTime, _ := time.Parse("02/01/2006 00:15", "18/04/2014 00:00")
	if value, ok := statsWriter.batches[firstBatchTime.Unix()]; !ok {
		t.Errorf("TestStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", firstBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	secondBatchTime := firstBatchTime.Add(time.Duration(24) * time.Hour)
	if value, ok := statsWriter.batches[secondBatchTime.Unix()]; !ok {
		t.Errorf("TestStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", secondBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}

	thirdBatchTime := firstBatchTime.Add(time.Duration(48) * time.Hour)
	if value, ok := statsWriter.batches[thirdBatchTime.Unix()]; !ok {
		t.Errorf("TestStreamerWithBufferedIO test failed: count not find a batch starting with a read time of [%v] in batches: [%v]", thirdBatchTime.Unix(), statsWriter.batches)
	} else {
		t.Logf("Value is [%s]", value)
	}
}
