package model

import (
	"github.com/alexandre-normand/glukit/app/util"
	"sort"
	"time"
)

type Time struct {
	Timestamp  int64  `json:"timestamp"`
	TimeZoneId string `json:"timezone"`
}

// GetTime gets the time of a Timestamp value
func (element Time) GetTime() (timeValue time.Time) {
	rawValue := time.Unix(int64(element.Timestamp/1000), 0)
	if location, err := util.GetOrLoadLocationForName(element.TimeZoneId); err != nil {
		return time.Unix(0, 0)
	} else {
		return rawValue.In(location)
	}

	return time.Unix(0, 0)
}

func (element Time) Format() (formatted string, err error) {
	timeValue := element.GetTime()
	if location, err := util.GetOrLoadLocationForName(element.TimeZoneId); err != nil {
		return "", err
	} else {
		return timeValue.In(location).Format(util.TIMEFORMAT), nil
	}
}

type TimeSlice []Time

func (slice TimeSlice) Len() int {
	return len(slice)
}

func (slice TimeSlice) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice TimeSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice TimeSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Timestamp / 1000
}

type Interface interface {
	sort.Interface
	GetEpochTime(i int) (epochTime int64)
}

// filter filters out any value that is outside of the lower and upper bounds. The two bounds are inclusive and the returned
// indexes are inclusive too.
func GetBoundariesOfElementsInRange(slice Interface, lowerBound, upperBound time.Time) (startIndex, endIndex int) {
	// Nothing to sort, return immediately
	if slice.Len() == 0 {
		return 0, slice.Len() - 1
	}

	arraySize := slice.Len()
	startIndex = 0
	endIndex = arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(slice)

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(slice.GetEpochTime(i), 0)
		if !elementTime.After(upperBound) {
			endIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(slice.GetEpochTime(i), 0)
		if !elementTime.Before(lowerBound) {
			startIndex = i
			break
		}
	}

	return startIndex, endIndex
}
