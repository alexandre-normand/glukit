// The util package contains general functions to support the other packages (timeutils, sysutils)
package util

import (
	"log"
	"time"
)

const (
	TIMEFORMAT        = "2006-01-02 15:04:05 MST"
	DRIVE_TIMEFORMAT  = "2006-01-02T15:04:05.000Z"
	TIMEZONE          = "PST"
	INTERNAL_TIMEZONE = "GMT"
)

var TIMEZONE_LOCATION, _ = time.LoadLocation("America/Los_Angeles")
var BEGINNING_OF_TIME = time.Unix(0, 0)

func ParseGoogleDriveDate(value string) (timeValue time.Time, err error) {
	return time.Parse(DRIVE_TIMEFORMAT, value)
}

func GetTimeInSeconds(timeValue string) (value int64) {
	if timeValue, err := time.Parse(TIMEFORMAT, timeValue+" "+INTERNAL_TIMEZONE); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}

func ParseTime(timeValue string, timezoneString string) (value time.Time, err error) {
	if value, err = time.Parse(TIMEFORMAT, timeValue+" "+timezoneString); err == nil {
		value = time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(),
			value.Nanosecond(), TIMEZONE_LOCATION)
	}

	return value, err
}

func GetEndOfDayBoundaryBefore(timeValue time.Time) (latestEndOfDayBoundary time.Time) {
	if timeValue.Hour() < 6 {
		// Rewind by one more day
		previousDay := timeValue.Add(time.Duration(-24 * time.Hour))
		latestEndOfDayBoundary = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	} else {
		latestEndOfDayBoundary = time.Date(timeValue.Year(), timeValue.Month(), timeValue.Day(), 6, 0, 0, 0, TIMEZONE_LOCATION)
	}

	return latestEndOfDayBoundary
}

func LocalTimeWithDefaultTimezone(timevalue time.Time) (localTime string) {
	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT)
}
