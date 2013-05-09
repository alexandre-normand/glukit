package timeutils

import (
	"time"
	"log"
)


const (
	TIMEFORMAT        = "2006-01-02 15:04:05 MST"
	DRIVE_TIMEFORMAT  = "2006-01-02T15:04:05.000Z"
	TIMEZONE          = "PST"
	INTERNAL_TIMEZONE = "GMT"
)

var TIMEZONE_LOCATION, _ = time.LoadLocation("America/Los_Angeles")

func ParseGoogleDriveDate(value string) (timeValue time.Time, err error) {
	return time.Parse(DRIVE_TIMEFORMAT, value)
}

func GetTimeInSeconds(timeValue string, timezoneString string) (value int64) {
	if timeValue, err := ParseTime(timeValue, timezoneString); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}

func ParseTime(timeValue string, timezoneString string) (value time.Time, err error) {
	if value, err = time.Parse(TIMEFORMAT, timeValue + " " + timezoneString); err == nil {
		value = time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(),
			value.Nanosecond(), TIMEZONE_LOCATION)
	}

	return value, err;
}

func LocalTimeWithDefaultTimezone(timevalue time.Time) (localTime string) {
	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT)
}
