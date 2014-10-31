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
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	r := make([]apimodel.GlucoseRead, NUM_READS_FOR_3_MONTHS)
	ct, _ := time.Parse("02/01/2006 15:04", "18/04/2014 00:00")
	for i := 0; i < NUM_READS_FOR_3_MONTHS; i++ {
		readTime := ct.Add(time.Duration(i) * time.Hour)
		r[i] = apimodel.GlucoseRead{apimodel.Time{apimodel.GetTimeMillis(readTime), "America/Los_Angeles"}, apimodel.MG_PER_DL, float32(80)}
	}

	a1cEstimate, err := engine.CalculateA1CEstimate(c, r)
	truncatedValue := truncateToOneDecimal(a1cEstimate.Value)
	if err != nil {
		t.Fatal(err)
	} else if truncatedValue != 4.4 {
		t.Errorf("TestA1cEstimationsFromFixedMeans failed: got an estimated a1c of %f but expected %f", truncatedValue, 4.4)
	}
}

func truncateToOneDecimal(value float64) float64 {
	return float64(int(value*10)) / 10
}
