package engine_test

import (
	"appengine/aetest"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/engine"
	"testing"
	"time"
)

const (
	NUM_READS_FOR_3_MONTHS = 288 * 90
)

func TestA1cEstimationsFromFixedMeans(t *testing.T) {
	testA1CEstimateFromFixedAverage(t, 65, 4.0)
	testA1CEstimateFromFixedAverage(t, 79, 4.4)
	testA1CEstimateFromFixedAverage(t, 90, 4.7)
	testA1CEstimateFromFixedAverage(t, 101, 5.0)
	testA1CEstimateFromFixedAverage(t, 158, 6.6)
	testA1CEstimateFromFixedAverage(t, 403, 13.5)
}

func testA1CEstimateFromFixedAverage(t *testing.T, average float32, expectedA1C float64) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	r := make([]apimodel.GlucoseRead, NUM_READS_FOR_3_MONTHS)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < NUM_READS_FOR_3_MONTHS; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, average}
	}

	a1cEstimate, err := engine.CalculateA1CEstimate(c, r)
	roundedValue := roundToOneDecimal(a1cEstimate.Value)
	if err != nil {
		t.Fatal(err)
	} else if roundedValue != expectedA1C {
		t.Errorf("TestA1cEstimationsFromFixedMeans failed: got an estimated a1c of %f but expected %f", roundedValue, expectedA1C)
	}
}

func roundToOneDecimal(value float64) float64 {
	return float64(int((value+0.05)*10)) / 10
}
