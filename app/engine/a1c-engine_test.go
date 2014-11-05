package engine_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/engine"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
	"github.com/alexandre-normand/glukit/app/streaming"
	"github.com/alexandre-normand/glukit/app/util"
	"github.com/alexandre-normand/glukit/lib/goauth2/oauth"
	"log"
	"sort"
	"testing"
	"time"
)

const (
	NUM_READS_FOR_3_MONTHS = 288 * 91
	TEST_USER              = "test@glukit.com"
)

func TestA1cEstimationsFromFixedMeans(t *testing.T) {
	testA1CEstimateFromFixedAverage(t, 65, 4.0)
	testA1CEstimateFromFixedAverage(t, 79, 4.4)
	testA1CEstimateFromFixedAverage(t, 90, 4.7)
	testA1CEstimateFromFixedAverage(t, 101, 5.0)
	testA1CEstimateFromFixedAverage(t, 158, 6.6)
	testA1CEstimateFromFixedAverage(t, 403, 13.5)
}

func TestCalculationWithInsufficientCoverage(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	r := make([]apimodel.GlucoseRead, 288*89)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < 288*89; i++ {
		readTime := ct.Add(time.Duration(i*5) * time.Minute)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(80)}
	}

	a1cEstimate, err := engine.CalculateA1CEstimate(c, r)
	if err == nil {
		t.Errorf("TestCalculationWithInsufficientCoverage failed: should return error when coverage is insufficient but calculated a1c of [%v]", a1cEstimate)
	}
}

func testA1CEstimateFromFixedAverage(t *testing.T, average float32, expectedA1C float64) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	a1cEstimate, err := engine.CalculateA1CEstimate(c, generateReadsWithFixedAverage(average, time.Now()))
	if err != nil {
		t.Fatal(err)
	} else if roundedValue := roundToOneDecimal(a1cEstimate.Value); roundedValue != expectedA1C {
		t.Errorf("TestA1cEstimationsFromFixedMeans failed: got an estimated a1c of %f but expected %f", roundedValue, expectedA1C)
	}
}

func generateReadsWithFixedAverage(average float32, upperDate time.Time) []apimodel.GlucoseRead {
	r := make([]apimodel.GlucoseRead, NUM_READS_FOR_3_MONTHS)

	for i := 0; i < NUM_READS_FOR_3_MONTHS; i++ {
		readTime := upperDate.Add((time.Duration(i*-5) * time.Minute))
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, average}
	}

	sortedReads := apimodel.GlucoseReadSlice(r)
	sort.Sort(sortedReads)
	log.Printf("Generated reads from [%s] to [%s]", r[0].GetTime(), r[len(r)-1].GetTime())
	return r
}

func setupTestData(t *testing.T, average float32, upperDate time.Time) (c aetest.Context, glukitUser *model.GlukitUser, key *datastore.Key) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	var oauthToken oauth.Token
	user := model.GlukitUser{TEST_USER, "", "", upperDate,
		"", "", util.GLUKIT_EPOCH_TIME, apimodel.UNDEFINED_GLUCOSE_READ, oauthToken, oauthToken.RefreshToken,
		model.UNDEFINED_SCORE, model.UNDEFINED_SCORE, false, "", upperDate}

	key, err = store.StoreUserProfile(c, upperDate, user)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Initialized [%s] with key [%v]", TEST_USER, key)

	dataStoreWriter := store.NewDataStoreGlucoseReadBatchWriter(c, key)
	batchingWriter := bufio.NewGlucoseReadWriterSize(dataStoreWriter, store.GLUKIT_SCORE_PUT_MULTI_SIZE)
	glucoseReadStreamer := streaming.NewGlucoseStreamerDuration(batchingWriter, apimodel.DAY_OF_DATA_DURATION)

	r := generateReadsWithFixedAverage(average, upperDate)
	glucoseReadStreamer, err = glucoseReadStreamer.WriteGlucoseReads(r)
	if err != nil {
		t.Fatal(err)
	}

	glucoseReadStreamer, err = glucoseReadStreamer.Close()
	if err != nil {
		t.Fatal(err)
	}

	return c, &user, key
}

func TestFetchAndEstimateFlow(t *testing.T) {
	upperDate, _ := time.Parse(util.TIMEFORMAT_NO_TZ, "2014-04-18 00:00:00")

	c, glukitUser, _ := setupTestData(t, 79, upperDate)
	defer c.Close()

	a1cEstimate, err := engine.EstimateA1C(c, glukitUser, upperDate)
	if err != nil {
		t.Fatal(err)
	}
	roundedValue := roundToOneDecimal(a1cEstimate.Value)
	if roundedValue != 4.4 {

		t.Errorf("TestFetchAndEstimateFlow failed: got an estimated a1c of %f but expected %f", roundedValue, 4.4)
	}

}

func roundToOneDecimal(value float64) float64 {
	return float64(int((value+0.05)*10)) / 10
}
