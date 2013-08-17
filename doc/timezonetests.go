package main

import (
	"fmt"
	"time"
)

// GetLocaltimeOffset returns the Fixed location extrapolated by calculating the offset
// of the localtime and the internal time in UTC
func GetLocaltimeOffset(localTime string, internalTime time.Time) (location *time.Location) {
	// Get the local time as if it was UTC (it's not)
	localTimeUTC, _ := time.Parse("2006-01-02 15:04:05", localTime)
	fmt.Printf("Local time is %s and internal is %s", localTimeUTC, internalTime)

	// Get the difference between the internal time (actual UTC) and the local time
	durationOffset := localTimeUTC.Sub(internalTime)

	locationName := fmt.Sprintf("%+03d%02d", int64(durationOffset.Hours()), (int64(durationOffset)-int64(durationOffset.Hours())*int64(time.Hour))/int64(time.Minute))
	fmt.Printf("name is %s\n", locationName)
	return time.FixedZone(locationName, int(durationOffset.Seconds()))
}

func main() {
	internalTime, err := time.Parse("2006-01-02 15:04:05", "2013-04-24 00:00:00")
	fmt.Printf("Internal time is %s, err is %v\n", internalTime, err)
	location := GetLocaltimeOffset("2013-04-24 02:00:00", internalTime)
	fmt.Printf("Timezone offset is %v", location)

	fmt.Printf("internal time converted to local is %s", internalTime.In(location).Format("2006-01-02 15:04:05"))
}
