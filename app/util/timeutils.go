package util

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"
)

const (
	// Default time format
	TIMEFORMAT = "2006-01-02 15:04:05 MST"

	// Default time format
	TIMEFORMAT_NO_TZ = "2006-01-02 15:04:05"

	// Time format used by the Google Drive api
	DRIVE_TIMEFORMAT = "2006-01-02T15:04:05.000Z"

	// Timezone for Dexcom interval time values
	INTERNAL_TIMEZONE = "GMT"

	// Let's make days end at 18h00
	HOUR_OF_END_OF_DAY = 18
)

var zoneNameRegexp = regexp.MustCompile("[+-](\\d){4}")

// Beginning of time should be unix epoch 0 but, to optimize some processing
// may iterate overtime starting at this value, we just define the notion
// of Glukit epoch time and have this value be set to something less far back
// but still before anything interesting happened in the Glukit world.
// This maps to 01 Jan 2004 00:00:00 GMT.
var GLUKIT_EPOCH_TIME = time.Unix(1072915200, 0)

var locationCache = make(map[string]*time.Location)

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

// GetTimeInSeconds parses a datetime string and returns its unix timestamp.
func GetTimeUTC(timeValue string) (time.Time, error) {
	// time values without timezone info are interpreted as UTC, which is perfect
	return time.Parse(TIMEFORMAT_NO_TZ, timeValue)

}

// GetEndOfDayBoundaryBefore returns the boundary of very last "end of day" before the given time.
// To give an example, if the given time is July 17th 8h00 PST, the boundary returned is going to be
// July 17th 06h00. If the time is July 17th 05h00 PST, the boundary returned is July 16th 06h00.
// Very important: The timeValue's location must be accurate!
func GetEndOfDayBoundaryBefore(timeValue time.Time) (latestEndOfDayBoundary time.Time) {
	if timeValue.Hour() < HOUR_OF_END_OF_DAY {
		// Rewind by one more day
		previousDay := timeValue.Add(time.Duration(-24 * time.Hour))
		latestEndOfDayBoundary = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), HOUR_OF_END_OF_DAY, 0, 0, 0, timeValue.Location())
	} else {
		latestEndOfDayBoundary = time.Date(timeValue.Year(), timeValue.Month(), timeValue.Day(), HOUR_OF_END_OF_DAY, 0, 0, 0, timeValue.Location())
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

// Returns the timevalue with its timezone set to UTC but without
// printing the timezone in the formatted string
func TimeInUTCNoTz(timevalue time.Time) (localTime string) {
	return timevalue.UTC().Format(TIMEFORMAT_NO_TZ)
}

func GetTimeWithImpliedLocation(localTime string, location *time.Location) (time.Time, error) {
	return time.ParseInLocation(TIMEFORMAT_NO_TZ, localTime, location)
}

// GetLocaltimeOffset returns the Fixed location extrapolated by calculating the offset
// of the localtime and the internal time in UTC
func GetLocaltimeOffset(localTime string, internalTime time.Time) (location *time.Location) {
	// Get the local time as if it was UTC (it's not)
	localTimeUTC, err := time.Parse(TIMEFORMAT_NO_TZ, localTime)
	if err != nil {
		Propagate(err)
	}

	// Get the difference between the internal time (actual UTC) and the local time
	durationOffset := localTimeUTC.Sub(internalTime)

	var truncatedDuration time.Duration
	if math.Signbit(durationOffset.Hours()) {
		minutesOffsetTruncated := int64(math.Ceil(durationOffset.Minutes()/15.-0.5) * 15.)
		truncatedDuration = time.Duration(minutesOffsetTruncated) * time.Minute
	} else {
		minutesOffsetTruncated := int64(math.Floor(durationOffset.Minutes()/15.+0.5) * 15.)
		truncatedDuration = time.Duration(minutesOffsetTruncated) * time.Minute
	}

	minutesOffsetPortion := float64((int64(truncatedDuration) - int64(truncatedDuration.Hours())*int64(time.Hour)) / int64(time.Minute))
	locationName := fmt.Sprintf("%+03d%02d", int(truncatedDuration.Hours()), int64(math.Abs(minutesOffsetPortion)))
	return time.FixedZone(locationName, int(durationOffset.Seconds()))
}

// GetLocalTimeInProperLocation returns the parsed local time with the location appropriately set as extrapolated
// by calculating the difference of the internal time vs the local time
func GetLocalTimeInProperLocation(localTime string, internalTime time.Time) (localTimeWithLocation time.Time) {
	location := GetLocaltimeOffset(localTime, internalTime)
	localTimeWithLocation, _ = time.ParseInLocation(TIMEFORMAT_NO_TZ, localTime, location)
	return
}

func GetOrLoadLocationForName(locationName string) (location *time.Location, err error) {
	if location, ok := locationCache[locationName]; !ok {
		location, err = time.LoadLocation(locationName)
		if err != nil {
			if !zoneNameRegexp.MatchString(locationName) {
				return nil, errors.New(fmt.Sprintf("Invalid location name, not a valid timezone location [%s]", locationName))
			} else {
				var hours, minutes int64
				fmt.Sscanf(locationName, "%+03d%02d", &hours, &minutes)
				offsetInMinutes := hours*int64(time.Duration(60)*time.Minute) + minutes
				offsetInSeconds := offsetInMinutes + minutes*int64(time.Duration(60)*time.Second)
				location = time.FixedZone(locationName, int(offsetInSeconds))
				locationCache[locationName] = location
			}
		}

		locationCache[locationName] = location
		return location, nil
	} else {
		return location, nil
	}
}
