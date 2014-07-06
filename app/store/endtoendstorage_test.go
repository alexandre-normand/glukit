package store_test

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"testing"
	"time"
)

func TestEndToEndMergeOfReadBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreGlucoseReadBatchWriter(c, key)
	bufferedWriter := bufio.NewGlucoseReadWriterSize(w, 5)
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
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 05:00")
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

func TestEndToEndMergeOfCalibrationBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreCalibrationBatchWriter(c, key)
	bufferedWriter := bufio.NewCalibrationWriterSize(w, 5)
	s := streaming.NewCalibrationReadStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	// Write first chunk
	r := make([]apimodel.CalibrationRead, 25)
	firstChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 01:00")
	for i := 0; i < 25; i++ {
		readTime := firstChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
	}
	s, _ = s.WriteCalibrations(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ := time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ := time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err := store.GetCalibrations(c, TEST_USER, lowerBound, upperBound)

	if err != nil {
		t.Fatal(err)
	}

	if reads == nil {
		t.Errorf("Reads is nil. Expected retrieval with results")
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First calibration of batch doesn't match expected time. Expected [%v], got [%v]", firstChunkStart, reads[0].GetTime())
	}

	if len(reads) != 25 {
		t.Errorf("Expected [25] calibrations but got [%d]: %v", len(reads), reads)
	}

	// Simulate a second batch of writes that overlaps partly with the second batch and confirm that the total
	// number of reads matches the expected (i.e. no duplicates)
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 05:00")
	r = make([]apimodel.CalibrationRead, 25)
	for i := 0; i < 25; i++ {
		readTime := secondChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.CalibrationRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(i)}
	}
	s, _ = s.WriteCalibrations(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ = time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ = time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err = store.GetCalibrations(c, TEST_USER, lowerBound, upperBound)
	if err != nil {
		t.Fatal(err)
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First calibration of batch doesn't match expected time. Expected [%v], got [%v]", secondChunkStart, reads[0].GetTime())
	}

	if len(reads) != 29 {
		t.Errorf("Expected [29] reads but got [%d]", len(reads))
	}
}

func TestEndToEndMergeOfInjectionBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreInjectionBatchWriter(c, key)
	bufferedWriter := bufio.NewInjectionWriterSize(w, 5)
	s := streaming.NewInjectionStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	// Write first chunk
	r := make([]apimodel.Injection, 25)
	firstChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 01:00")
	for i := 0; i < 25; i++ {
		readTime := firstChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"}
	}
	s, _ = s.WriteInjections(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ := time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ := time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err := store.GetInjections(c, TEST_USER, lowerBound, upperBound)

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
		t.Errorf("Expected [25] injections but got [%d]: %v", len(reads), reads)
	}

	// Simulate a second batch of writes that overlaps partly with the second batch and confirm that the total
	// number of reads matches the expected (i.e. no duplicates)
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 05:00")
	r = make([]apimodel.Injection, 25)
	for i := 0; i < 25; i++ {
		readTime := secondChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Injection{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), "Humalog", "Bolus"}
	}
	s, _ = s.WriteInjections(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ = time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ = time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err = store.GetInjections(c, TEST_USER, lowerBound, upperBound)
	if err != nil {
		t.Fatal(err)
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First injection of batch doesn't match expected time. Expected [%v], got [%v]", secondChunkStart, reads[0].GetTime())
	}

	if len(reads) != 29 {
		t.Errorf("Expected [29] injections but got [%d]", len(reads))
	}
}

func TestEndToEndMergeOfExerciseBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreExerciseBatchWriter(c, key)
	bufferedWriter := bufio.NewExerciseWriterSize(w, 5)
	s := streaming.NewExerciseStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	// Write first chunk
	r := make([]apimodel.Exercise, 25)
	firstChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 01:00")
	for i := 0; i < 25; i++ {
		readTime := firstChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"}
	}
	s, _ = s.WriteExercises(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ := time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ := time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err := store.GetExercises(c, TEST_USER, lowerBound, upperBound)

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
		t.Errorf("Expected [25] injections but got [%d]: %v", len(reads), reads)
	}

	// Simulate a second batch of writes that overlaps partly with the second batch and confirm that the total
	// number of reads matches the expected (i.e. no duplicates)
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 05:00")
	r = make([]apimodel.Exercise, 25)
	for i := 0; i < 25; i++ {
		readTime := secondChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Exercise{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, i, "Light", "details"}
	}
	s, _ = s.WriteExercises(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ = time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ = time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err = store.GetExercises(c, TEST_USER, lowerBound, upperBound)
	if err != nil {
		t.Fatal(err)
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First injection of batch doesn't match expected time. Expected [%v], got [%v]", secondChunkStart, reads[0].GetTime())
	}

	if len(reads) != 29 {
		t.Errorf("Expected [29] injections but got [%d]", len(reads))
	}
}

func TestEndToEndMergeOfMealBatches(t *testing.T) {
	c, key := setup(t)
	defer c.Close()

	w := store.NewDataStoreMealBatchWriter(c, key)
	bufferedWriter := bufio.NewMealWriterSize(w, 5)
	s := streaming.NewMealStreamerDuration(bufferedWriter, apimodel.DAY_OF_DATA_DURATION)

	// Write first chunk
	r := make([]apimodel.Meal, 25)
	firstChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 01:00")
	for i := 0; i < 25; i++ {
		readTime := firstChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)}
	}
	s, _ = s.WriteMeals(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ := time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ := time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err := store.GetMeals(c, TEST_USER, lowerBound, upperBound)

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
		t.Errorf("Expected [25] injections but got [%d]: %v", len(reads), reads)
	}

	// Simulate a second batch of writes that overlaps partly with the second batch and confirm that the total
	// number of reads matches the expected (i.e. no duplicates)
	secondChunkStart, _ := time.Parse("02/01/2006 15:04", "18/04/2015 05:00")
	r = make([]apimodel.Meal, 25)
	for i := 0; i < 25; i++ {
		readTime := secondChunkStart.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.Meal{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Montreal"}, float32(i), float32(i + 1), float32(i + 2), float32(i + 3)}
	}
	s, _ = s.WriteMeals(r)
	s, _ = s.Flush()
	s, _ = s.Close()

	lowerBound, _ = time.Parse("02/01/2006 15:04", "18/04/2015 00:00")
	upperBound, _ = time.Parse("02/01/2006 15:04", "20/04/2015 02:00")
	reads, err = store.GetMeals(c, TEST_USER, lowerBound, upperBound)
	if err != nil {
		t.Fatal(err)
	}

	if reads[0].GetTime().Unix() != firstChunkStart.Unix() {
		t.Errorf("First injection of batch doesn't match expected time. Expected [%v], got [%v]", secondChunkStart, reads[0].GetTime())
	}

	if len(reads) != 29 {
		t.Errorf("Expected [29] injections but got [%d]", len(reads))
	}
}
