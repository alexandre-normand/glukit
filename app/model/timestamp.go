package model

import (
	"sort"
	"time"
)

// Timestamp represents hold a timestamp and a localtime string
type Timestamp struct {
	LocalTime string `json:"label" datastore:"localtime,noindex"`
	EpochTime int64  `json:"x" datastore:"timestamp"`
}

// GetTime gets the time of a Timestamp value
func (element Timestamp) GetTime() (timeValue time.Time) {
	value := time.Unix(int64(element.EpochTime), 0)
	return value
}

type TimestampSlice []Timestamp

func (slice TimestampSlice) Len() int {
	return len(slice)
}

func (slice TimestampSlice) Less(i, j int) bool {
	return slice[i].EpochTime < slice[j].EpochTime
}

func (slice TimestampSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// filter filters out any value that is outside of the lower and upper bounds. The two bounds are inclusive.
func (slice TimestampSlice) Filter(lowerBound, upperBound time.Time) (startIndex, endIndex int) {
	// Nothing to sort, return immediately
	if len(slice) == 0 {
		return 0, slice.Len() - 1
	}

	arraySize := slice.Len()
	startIndex = 0
	endIndex = arraySize - 1

	// Sort might not be strictly needed depending on the ordering of the datastore loading but since there doesn't
	// seem to be any warranty, sorting seems like a good idea
	sort.Sort(slice)

	for i := arraySize - 1; i > 0; i-- {
		elementTime := time.Unix(slice[i].EpochTime, 0)
		if !elementTime.After(upperBound) {
			endIndex = i
			break
		}
	}

	for i := 0; i < arraySize; i++ {
		elementTime := time.Unix(slice[i].EpochTime, 0)
		if !elementTime.Before(lowerBound) {
			startIndex = i
			break
		}
	}

	return startIndex, endIndex
}
