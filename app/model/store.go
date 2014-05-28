package model

import (
	"appengine/datastore"
	"fmt"
	"strconv"
	"time"
)

func (dayOfReads *DayOfGlucoseReads) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var unit string
	var locationName string

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = columnValue.(time.Time)
		case columnName == "unit":
			unit = columnValue.(string)
		case columnName == "timezone":
			locationName = columnValue.(string)
		case columnName == "endTime":
			// We ignore it on load
			_ = columnValue.(time.Time)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				return err
			}

			readTime := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to int64 and cast it as int
			value := columnValue.(float64)

			read := GlucoseRead{Time{GetTimeMillis(readTime), locationName}, unit, float32(value)}
			dayOfReads.Reads = append(dayOfReads.Reads, read)
		}
	}

	return nil
}

func (dayOfReads *DayOfGlucoseReads) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfReads.Reads)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	reads := dayOfReads.Reads
	startTimestamp := reads[0].GetTime().Unix()
	startTime := reads[0].GetTime()
	endTime := reads[size-1].GetTime()

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}
	channel <- datastore.Property{
		Name:  "timezone",
		Value: reads[0].Time.TimeZoneId,
	}
	channel <- datastore.Property{
		Name:  "unit",
		Value: reads[0].Unit,
	}

	for _, read := range reads {
		readOffset := read.GetTime().Unix() - startTimestamp
		channel <- datastore.Property{
			Name:    strconv.FormatInt(readOffset, 10),
			Value:   float64(read.Value),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfReads *DayOfCalibrationReads) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var unit string
	var locationName string

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = columnValue.(time.Time)
		case columnName == "endTime":
			// We ignore it on load
			_ = columnValue.(time.Time)
		case columnName == "unit":
			unit = columnValue.(string)
		case columnName == "timezone":
			locationName = columnValue.(string)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				return err
			}

			readTime := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to int64 and cast it as int
			value := columnValue.(float64)

			read := CalibrationRead{Time{GetTimeMillis(readTime), locationName}, unit, float32(value)}
			dayOfReads.Reads = append(dayOfReads.Reads, read)
		}
	}

	return nil
}

func (dayOfReads *DayOfCalibrationReads) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfReads.Reads)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	reads := dayOfReads.Reads
	startTimestamp := reads[0].GetTime().Unix()
	startTime := reads[0].GetTime()
	endTime := reads[size-1].GetTime()

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}
	channel <- datastore.Property{
		Name:  "timezone",
		Value: reads[0].Time.TimeZoneId,
	}
	channel <- datastore.Property{
		Name:  "unit",
		Value: reads[0].Unit,
	}

	for _, read := range reads {
		readOffset := read.GetTime().Unix() - startTimestamp
		channel <- datastore.Property{
			Name:    strconv.FormatInt(readOffset, 10),
			Value:   float64(read.Value),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfInjections *DayOfInjections) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var locationName string

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = columnValue.(time.Time)
		case columnName == "endTime":
			// We ignore it on load
			_ = columnValue.(time.Time)
		case columnName == "timezone":
			locationName = columnValue.(string)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				return err
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to float64 and we downcast to float32 (it's safe since we only up-casted it for
			// the store
			value := float32(columnValue.(float64))

			injection := Injection{Time{GetTimeMillis(timestamp), locationName}, value, "Not implemented", "Not implemented"}
			dayOfInjections.Injections = append(dayOfInjections.Injections, injection)
		}
	}

	return nil
}

func (dayOfInjections *DayOfInjections) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfInjections.Injections)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	injections := dayOfInjections.Injections
	startTimestamp := injections[0].GetTime().Unix()
	startTime := injections[0].GetTime()
	endTime := injections[size-1].GetTime()

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}
	channel <- datastore.Property{
		Name:  "timezone",
		Value: injections[0].Time.TimeZoneId,
	}

	for _, injection := range injections {
		offset := injection.GetTime().Unix() - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   float64(injection.Units),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfMeals *DayOfMeals) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var locationName string

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = columnValue.(time.Time)
		case columnName == "endTime":
			// We ignore it on load
			_ = columnValue.(time.Time)
		case columnName == "timezone":
			locationName = columnValue.(string)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				return err
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to float64 and we downcast to float32 (it's safe since we only up-casted it for
			// the store
			carbs := float32(columnValue.(float64))

			meal := Meal{Time{GetTimeMillis(timestamp), locationName}, carbs, 0., 0., 0.}
			dayOfMeals.Meals = append(dayOfMeals.Meals, meal)
		}
	}

	return nil
}

func (dayOfMeals *DayOfMeals) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfMeals.Meals)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	meals := dayOfMeals.Meals
	startTimestamp := meals[0].GetTime().Unix()
	startTime := meals[0].GetTime()
	endTime := meals[size-1].GetTime()

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}
	channel <- datastore.Property{
		Name:  "timezone",
		Value: meals[0].Time.TimeZoneId,
	}

	for _, meal := range meals {
		offset := meal.GetTime().Unix() - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   float64(meal.Carbohydrates),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfExercises *DayOfExercises) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var locationName string

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = columnValue.(time.Time)
		case columnName == "endTime":
			// We ignore it on load
			_ = columnValue.(time.Time)
		case columnName == "timezone":
			locationName = columnValue.(string)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				return err
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We split the value string to extract the two components from it
			value := columnValue.(string)
			var duration int
			var intensity string
			fmt.Sscanf(value, EXERCISE_VALUE_FORMAT, &duration, &intensity)

			exercise := Exercise{Time{GetTimeMillis(timestamp), locationName}, duration, intensity, ""}
			dayOfExercises.Exercises = append(dayOfExercises.Exercises, exercise)
		}
	}

	return nil
}

func (dayOfExercises *DayOfExercises) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfExercises.Exercises)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	exercises := dayOfExercises.Exercises
	startTimestamp := exercises[0].GetTime().Unix()
	startTime := exercises[0].GetTime()
	endTime := exercises[size-1].GetTime()

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}
	channel <- datastore.Property{
		Name:  "timezone",
		Value: exercises[0].Time.TimeZoneId,
	}

	for _, exercise := range exercises {
		offset := exercise.GetTime().Unix() - startTimestamp
		// We need to store two values so we use a string and format each value inside of a single cell value
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   fmt.Sprintf(EXERCISE_VALUE_FORMAT, exercise.DurationMinutes, exercise.Intensity),
			NoIndex: true,
		}
	}

	return nil
}
