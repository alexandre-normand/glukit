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
		t.Errorf("Invalid location [%s]: [%v]", location, err)
	}
}
