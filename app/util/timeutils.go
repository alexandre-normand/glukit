package util

import (
	"fmt"
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

// Beginning of time should be unix epoch 0 but, to optimize some processing
// may iterate overtime starting at this value, we just define the notion
// of Glukit epoch time and have this value be set to something less far back
// but still before anything interesting happened in the Glukit world.
// This maps to 01 Jan 2004 00:00:00 GMT.
var GLUKIT_EPOCH_TIME = time.Unix(1072915200, 0)

// ParseGoogleDriveDate parses a Google Drive API time value
func ParseGoogleDriveDate(value string) (timeValue time.Time, err error) {
	return time.Parse(DRIVE_TIMEFORMAT, value)
}

// GetTimeInSeconds parses a datetime string and returns its unix timestamp.
func GetTimeInSeconds(timeValue string) (value int64) {
	// time values without timezone info are interpreted as UTC, which is perfect
	if timeValue, err := time.Parse(TIMEFORMAT_NO_TZ, timeValue); err == nil {
		return timeValue.Unix()
	} else {
		log.Printf("Error parsing string", err)
	}
	return 0
}

// GetEndOfDayBoundaryBefore returns the boundary of very last "end of day" before the given time.
// To give an example, if the given time is July 17th 8h00 PST, the boundary returned is going to be
// July 17th 06h00. If the time is July 17th 05h00 PST, the boundary returned is July 16th 06h00.
// Very important: The timeValue's location must be accurate!
func GetEndOfDayBoundaryBefore(timeValue time.Time) (latestEndOfDayBoundary time.Time) {
	if timeValue.Hour() < 6 {
		// Rewind by one more day
		previousDay := timeValue.Add(time.Duration(-24 * time.Hour))
		latestEndOfDayBoundary = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, timeValue.Location())
	} else {
		latestEndOfDayBoundary = time.Date(timeValue.Year(), timeValue.Month(), timeValue.Day(), 6, 0, 0, 0, timeValue.Location())
	}

	return latestEndOfDayBoundary
}

// GetMidnightUTCBefore returns the boundary of very last occurence of midnight before the given time.
// To give an example, if the given time is July 17th 2h00 UTC, the boundary returned is going to be
// July 17th 00h00. If the time is July 16th 23h00 PST, the boundary returned is July 16th 00h00.
func GetMidnightUTCBefore(timeValue time.Time) (latestMidnightBoundary time.Time) {
	timeInUTC := timeValue.UTC()
	latestMidnightBoundary = time.Date(timeInUTC.Year(), timeInUTC.Month(), timeInUTC.Day(), 0, 0, 0, 0, time.UTC)
	return latestMidnightBoundary
}

// Returns the timevalue with its timezone set to the default TIMEZONE_LOCATION
// func TimeWithDefaultTimezone(timevalue time.Time) (localTime string) {
// 	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT)
// }

// Returns the timevalue with its timezone set to the default TIMEZONE_LOCATION but without
// printing the timezone in the formatted string
// func TimeInDefaultTimezoneNoTz(timevalue time.Time) (localTime string) {
// 	return timevalue.In(TIMEZONE_LOCATION).Format(TIMEFORMAT_NO_TZ)
// }

// Returns the timevalue with its timezone set to UTC but without
// printing the timezone in the formatted string
func TimeInUTCNoTz(timevalue time.Time) (localTime string) {
	return timevalue.UTC().Format(TIMEFORMAT_NO_TZ)
}

// GetLocaltimeOffset returns the Fixed location extrapolated by calculating the offset
// of the localtime and the internal time in UTC
func GetLocaltimeOffset(localTime string, internalTime time.Time) (location *time.Location) {
	// Get the local time as if it was UTC (it's not)
	localTimeUTC, _ := time.Parse(TIMEFORMAT_NO_TZ, localTime)

	// Get the difference between the internal time (actual UTC) and the local time
	durationOffset := localTimeUTC.Sub(internalTime)

	locationName := fmt.Sprintf("%+03d%02d", int64(durationOffset.Hours()), (int64(durationOffset)-int64(durationOffset.Hours())*int64(time.Hour))/int64(time.Minute))
	return time.FixedZone(locationName, int(durationOffset.Seconds()))
}

// GetLocalTimeInProperLocation returns the parsed local time with the location appropriately set as extrapolated
// by calculating the difference of the internal time vs the local time
func GetLocalTimeInProperLocation(localTime string, internalTime time.Time) (localTimeWithLocation time.Time) {
	location := GetLocaltimeOffset(localTime, internalTime)
	localTimeWithLocation, _ = time.ParseInLocation(TIMEFORMAT_NO_TZ, localTime, location)
	return
}
