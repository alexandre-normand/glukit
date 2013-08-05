package util

import (
	"log"
	"time"
)

const (
	// Default time format
	TIMEFORMAT = "2006-01-02 15:04:05 MST"

	// Default time format
	TIMEFORMAT_NO_TZ = "2006-01-02 15:04:05"

	// Time format used by the Google Drive api
	DRIVE_TIMEFORMAT = "2006-01-02T15:04:05.000Z"
	// Default timezone
	TIMEZONE = "PST"
	// Timezone for Dexcom interval time values
	INTERNAL_TIMEZONE = "GMT"
)

// The default location is assumed to be America/Los_Angeles because I use
// my data and we're in this timezone. This will have to change before
// this goes live
var TIMEZONE_LOCATION, _ = time.LoadLocation("America/Los_Angeles")

var UTC_LOCATION, _ = time.LoadLocation("UTC")
var BEGINNING_OF_TIME = time.Unix(0, 0)

// ParseGoogleDriveDate parses a Google Drive API time value
func ParseGoogleDriveDate(value string) (timeValue time.Time, err error) {
	return time.Parse(DRIVE_TIMEFORMAT, value)
}

// GetTimeInSeconds parses a datetime string and returns its unix timestamp.
func GetTimeInSeconds(timeValue string) (value int64) {
	if timeValue, err := time.Parse(TIMEFORMAT, timeValue+" "+INTERNAL_TIMEZONE); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}

// ParseTime parses a datetime string as a time.Time value. This currently assumes the default
// location of "America/Los_Angeles". I know, it's terrible
func ParseTime(timeValue string, timezoneString string) (value time.Time, err error) {
	if value, err = time.Parse(TIMEFORMAT, timeValue+" "+timezoneString); err == nil {
		value = time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(),
			value.Nanosecond(), TIMEZONE_LOCATION)
	}

	return value, err
}

// GetEndOfDayBoundaryBefore returns the boundary of very last "end of day" before the given time.
// To give an example, if the given time is July 17th 8h00 PST, the boundary returned is going to be
// July 17th 06h00. If the time is July 17th 05h00 PST, the boundary returned is July 16th 06h00.
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

// Returns the timevalue with its timezone set to the default TIMEZONE_LOCATION
func TimeWithDefaultTimezone(timevalue time.Time) (localTime string) {
	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT)
}

// Returns the timevalue with its timezone set to the default TIMEZONE_LOCATION but without
// printing the timezone in the formatted string
func TimeInDefaultTimezoneNoTz(timevalue time.Time) (localTime string) {
	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT_NO_TZ)
}

// Returns the timevalue with its timezone set to UTC but without
// printing the timezone in the formatted string
func TimeInUTCNoTz(timevalue time.Time) (localTime string) {
	return timevalue.In(UTC_LOCATION).Format(TIMEFORMAT_NO_TZ)
}
