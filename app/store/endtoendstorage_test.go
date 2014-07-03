package store_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"testing"
	"time"
)

func TestEndToEndMergeOfBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreGlucoseReadBatchWriter(c, key)
	bufferedWriter := bufio.NewGlucoseReadWriterSize(w, 2)
	s := streaming.NewGlucoseStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	// Write first chunk
	r := make([]apimodel.GlucoseRead, 25)
	firstChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 01:00")
	for i := 0; i < 25; i++ {
		readTime := firstChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
	}
	s, _ = s.WriteGlucoseReads(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ := time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ := time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err := store.GetGlucoseReads(c, TEST_USER, lowerBound, upperBound)

	if err != nil {
		t.Fatal(err)
	}

	if reads == nil {
		t.Errorf("Reads is nil. Expected retrieval with results")
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First reads of batch doesn't match expected time. Expected [%v], got [%v]", firstChunkStart, reads[0].GetTime())
	}

	if len(reads) != 25 {
		t.Errorf("Expected [25] reads but got [%d]: %v", len(reads), reads)
	}

	// Simulate a second batch of writes that overlaps partly with the second batch and confirm that the total
	// number of reads matches the expected (i.e. no duplicates)
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 22:00")
	r = make([]apimodel.GlucoseRead, 25)
	for i := 0; i < 25; i++ {
		readTime := secondChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
	}
	s, _ = s.WriteGlucoseReads(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ = time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ = time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err = store.GetGlucoseReads(c, TEST_USER, lowerBound, upperBound)
	if err != nil {
		t.Fatal(err)
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First reads of batch doesn't match expected time. Expected [%v], got [%v]", secondChunkStart, reads[0].GetTime())
	}

	if len(reads) != 29 {
		t.Errorf("Expected [29] reads but got [%d]", len(reads))
	}
}
