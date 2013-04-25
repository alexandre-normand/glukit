package datautils

import (
	"models"
	"timeutils"
	"time"
	"strconv"
)

// Assumes reads are ordered by time
func GetLastDayOfData(meterReads []models.MeterRead, injections []models.Injection, carbIntakes []models.CarbIntake) (lastDayOfReads []models.MeterRead, lastDayOfInjections []models.Injection, lastDayOfCarbIntakes []models.CarbIntake) {
	dataSize := len(meterReads)

	lastValue := meterReads[dataSize - 1]
	lastTime, _ := timeutils.ParseTime(lastValue.LocalTime, timeutils.TIMEZONE)
	var upperBound time.Time;
	if (lastTime.Hour() < 6) {
		// Rewind by one more day
		previousDay := lastTime.Add(time.Duration(-24*time.Hour))
		upperBound = time.Date(previousDay.Year(), previousDay.Month(), previousDay.Day(), 6, 0, 0, 0, timeutils.TIMEZONE_LOCATION)
	} else {
		upperBound = time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 6, 0, 0, 0, timeutils.TIMEZONE_LOCATION)
	}
	lowerBound := upperBound.Add(time.Duration(-24*time.Hour))

	return GetLastDayOfReads(meterReads, lowerBound, upperBound), GetLastDayOfInjections(injections, lowerBound, upperBound), GetLastDayOfCarbIntakes(carbIntakes, lowerBound, upperBound)
}

func GetLastDayOfReads(data []models.MeterRead, lowerBound, upperBound time.Time) (filteredData []models.MeterRead) {
	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	for i := arraySize - 1; i > 0; i-- {
		element := models.PointData(data[i])
		elementTime := element.GetTime()
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex:endOfDayIndex + 1]
}

func GetLastDayOfInjections(data []models.Injection, lowerBound, upperBound time.Time) (filteredData []models.Injection) {
	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	for i := arraySize - 1; i > 0; i-- {
		element := models.PointData(data[i])
		elementTime := element.GetTime()
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex:endOfDayIndex + 1]
}

func GetLastDayOfCarbIntakes(data []models.CarbIntake, lowerBound, upperBound time.Time) (filteredData []models.CarbIntake) {
	arraySize := len(data)
	startOfDayIndex := -1
	endOfDayIndex := -1

	for i := arraySize - 1; i > 0; i-- {
		element := models.PointData(data[i])
		elementTime := element.GetTime()
		if endOfDayIndex < 0 && elementTime.Before(upperBound) {
			endOfDayIndex = i
		}

		if startOfDayIndex < 0 && elementTime.Before(lowerBound) {
			startOfDayIndex = i + 1
		}
	}

	return data[startOfDayIndex:endOfDayIndex + 1]
}

func BuildHistogram(reads []models.MeterRead) (histogram map[string] int) {
	histogram = make(map[string] int)

	for i := range reads {
		currentReadValue := strconv.Itoa(reads[i].Value)
		currentValue, found := histogram[currentReadValue]
		if found {
			histogram[currentReadValue] = currentValue + 1
		} else {
			histogram[currentReadValue] = 1
		}
	}

	return histogram
}
