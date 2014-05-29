package util_test

import (
	. "github.com/alexandre-normand/glukit/app/util"
	"testing"
	"time"
)

func TestGetTimezone(t *testing.T) {
	internalTime, _ := time.Parse(TIMEFORMAT_NO_TZ, "2014-04-18 00:00:00")

	location := GetLocaltimeOffset("2014-04-18 01:28:00", internalTime)
	t.Logf("Location is [%v]", location)

	if location.String() != "+0130" {
		t.Errorf("Expected location [+0130] but got [%s]", location.String())
	}

	if _, err := GetOrLoadLocationForName(location.String()); err != nil {
		t.Errorf("Invalid location [%s]: [%v]", location.String(), err)
	}
}

func TestGetTimezoneNegativeOffset(t *testing.T) {
	internalTime, _ := time.Parse(TIMEFORMAT_NO_TZ, "2014-04-18 09:00:00")

	location := GetLocaltimeOffset("2014-04-18 02:02:00", internalTime)
	t.Logf("Location is [%v]", location)

	if location.String() != "-0700" {
		t.Errorf("Expected location [-0700] but got [%s]", location.String())
	}

	if _, err := GetOrLoadLocationForName(location.String()); err != nil {
		t.Errorf("Invalid location [%s]: [%v]", location.String(), err)
	}
}

func TestKnownLocationLoading(t *testing.T) {
	location := "America/Montreal"
	if _, err := GetOrLoadLocationForName(location); err != nil {
		t.Errorf("Should be a valid location location [%s] but got error: [%v]", location, err)
	}
}

func TestGetTimeWithImpliedLocation(t *testing.T) {
	locationName := "America/Los_Angeles"
	location, err := GetOrLoadLocationForName(locationName)
	if err != nil {
		t.Errorf("Should be a valid location location [%s] but got error: [%v]", locationName, err)
	}

	timeValue, err := GetTimeWithImpliedLocation("2014-04-23 13:33:57", location)
	if err != nil {
		t.Errorf("Should be able to parse time but got error: [%v]", err)
	}

	expected := int64(1398285237)
	if timeValue.Unix() != expected {
		t.Errorf("Expected timestamp [%d] but got [%d]", expected, timeValue.Unix())
	}
}
