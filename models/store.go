package models

import (
	"appengine/datastore"
	"fmt"
	"log"
	"strconv"
	"sysutils"
	"time"
	"timeutils"
)

func (dayOfReads *DayOfReads) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var endTime time.Time

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = property.Value.(time.Time)
			log.Printf("Parsing block of reads with starttime: %s", startTime)
		case columnName == "endTime":
			endTime = property.Value.(time.Time)
			log.Printf("Parsed block of reads with endtime: %s", endTime)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				sysutils.Propagate(err)
			}

			readTime := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to int64 and cast it as int
			value := int(columnValue.(int64))

			localTime := timeutils.LocalTimeWithDefaultTimezone(readTime)
			//context.Infof("Loading read with info: %s, %s, %d", localTime, readTime, value)
			read := GlucoseRead{localTime, TimeValue(readTime.Unix()), value}
			dayOfReads.Reads = append(dayOfReads.Reads, read)
		}
	}

	log.Printf("Loaded total of %d reads", len(dayOfReads.Reads))
	return nil
}

func (dayOfReads *DayOfReads) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfReads.Reads)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	reads := dayOfReads.Reads
	startTimestamp := int64(reads[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(reads[size-1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}

	for i := range reads {
		//context.Infof("Sending read to channel %s, %d, %d", reads[i].LocalTime, int64(reads[i].TimeValue), int64(reads[i].Value))
		readTimestamp := int64(reads[i].TimeValue)
		readOffset := readTimestamp - startTimestamp
		channel <- datastore.Property{
			Name:    strconv.FormatInt(readOffset, 10),
			Value:   int64(reads[i].Value),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfInjections *DayOfInjections) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var endTime time.Time

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = property.Value.(time.Time)
			log.Printf("Parsing block of injections with starttime: %s", startTime)
		case columnName == "endTime":
			endTime = property.Value.(time.Time)
			log.Printf("Parsed block of injections with endtime: %s", endTime)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				sysutils.Propagate(err)
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to float64 and we downcast to float32 (it's safe since we only up-casted it for
			// the store
			value := float32(columnValue.(float64))

			localTime := timeutils.LocalTimeWithDefaultTimezone(timestamp)
			injection := Injection{localTime, TimeValue(timestamp.Unix()), value, UNDEFINED_READ}
			dayOfInjections.Injections = append(dayOfInjections.Injections, injection)
		}
	}

	log.Printf("Loaded total of %d injections", len(dayOfInjections.Injections))
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
	startTimestamp := int64(injections[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(injections[size-1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}

	for i := range injections {
		timestamp := int64(injections[i].TimeValue)
		offset := timestamp - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   float64(injections[i].Units),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfCarbs *DayOfCarbs) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var endTime time.Time

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = property.Value.(time.Time)
			log.Printf("Parsing block of carbs with starttime: %s", startTime)
		case columnName == "endTime":
			endTime = property.Value.(time.Time)
			log.Printf("Parsed block of carbs with endtime: %s", endTime)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				sysutils.Propagate(err)
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We need to convert value to float64 and we downcast to float32 (it's safe since we only up-casted it for
			// the store
			value := float32(columnValue.(float64))

			localTime := timeutils.LocalTimeWithDefaultTimezone(timestamp)
			carb := Carb{localTime, TimeValue(timestamp.Unix()), value, UNDEFINED_READ}
			dayOfCarbs.Carbs = append(dayOfCarbs.Carbs, carb)
		}
	}

	log.Printf("Loaded total of %d carbs", len(dayOfCarbs.Carbs))
	return nil
}

func (dayOfCarbs *DayOfCarbs) Save(channel chan<- datastore.Property) error {
	defer close(channel)

	size := len(dayOfCarbs.Carbs)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	carbs := dayOfCarbs.Carbs
	startTimestamp := int64(carbs[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(carbs[size-1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}

	for i := range carbs {
		timestamp := int64(carbs[i].TimeValue)
		offset := timestamp - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   float64(carbs[i].Grams),
			NoIndex: true,
		}
	}

	return nil
}

func (dayOfExercises *DayOfExercises) Load(channel <-chan datastore.Property) error {
	var startTime time.Time
	var endTime time.Time

	for property := range channel {
		switch columnName, columnValue := property.Name, property.Value; {
		case columnName == "startTime":
			startTime = property.Value.(time.Time)
			log.Printf("Parsing block of exercises with starttime: %s", startTime)
		case columnName == "endTime":
			endTime = property.Value.(time.Time)
			log.Printf("Parsed block of exercises with endtime: %s", endTime)
		default:
			offsetInSeconds, err := strconv.ParseInt(columnName, 10, 64)
			if err != nil {
				sysutils.Propagate(err)
			}

			timestamp := time.Unix(startTime.Unix()+offsetInSeconds, 0)
			// We split the value string to extract the two components from it
			value := columnValue.(string)
			var duration int
			var intensity string
			fmt.Sscanf(value, EXERCISE_VALUE_FORMAT, &duration, &intensity)

			localTime := timeutils.LocalTimeWithDefaultTimezone(timestamp)
			exercise := Exercise{localTime, TimeValue(timestamp.Unix()), duration, intensity}
			dayOfExercises.Exercises = append(dayOfExercises.Exercises, exercise)
		}
	}

	log.Printf("Loaded total of %d exercises", len(dayOfExercises.Exercises))
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
	startTimestamp := int64(exercises[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(exercises[size-1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value: startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value: endTime,
	}

	for i := range exercises {
		timestamp := int64(exercises[i].TimeValue)
		offset := timestamp - startTimestamp
		// We need to store two values so we use a string and format each value inside of a single cell value
		channel <- datastore.Property{
			Name:    strconv.FormatInt(offset, 10),
			Value:   fmt.Sprintf(EXERCISE_VALUE_FORMAT, exercises[i].DurationInMinutes, exercises[i].Intensity),
			NoIndex: true,
		}
	}

	return nil
}
