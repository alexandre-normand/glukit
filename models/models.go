package models

import (
	"time"
	"timeutils"
	"strconv"
	"sysutils"
	"appengine/datastore"
	"log"
	"fmt"
)

const (
	UNDEFINED_READ = -1;
	EXERCISE_VALUE_FORMAT = "%d,%s"
)

type TimeValue int64

type TrackingData struct {
	Mean             float64         `json:"mean"`
	Median           float64         `json:"median"`
	Deviation        float64         `json:"deviation"`
	Min      	    float64         `json:"min"`
	Max      	    float64         `json:"max"`
	Distribution     []Coordinate    `json:"distribution"`
}

type Coordinate struct {
	X    int   `json:"x"`
	Y    int   `json:"y"`
}

type DataPoint struct {
	LocalTime string    `json:"label"`
	TimeValue TimeValue `json:"x"`
	Y         int       `json:"y"`
	Value     float32   `json:"r"`
}

type GlucoseRead struct {
	LocalTime string    `json:"label" datastore:"localtime,noindex"`
	TimeValue TimeValue `json:"x" datastore:"timestamp"`
	Value     int       `json:"y" datastore:"value,noindex"`
}

// Used to model reads in a short-wide fashion.
type DayOfReads struct {
	Reads          []GlucoseRead
}

type DayOfInjections struct {
	Injections     []Injection
}

type DayOfCarbs struct {
	Carbs          []Carb
}

type DayOfExercises struct {
	Exercises      []Exercise
}

type Injection struct {
	LocalTime          string       `json:"label" datastore:"localtime,noindex"`
	TimeValue          TimeValue    `json:"x" datastore:"timestamp"`
	Units              float32      `json:"units" datastore:"units,noindex"`
	ReferenceReadValue int          `json:"y" datastore:"referenceReadValue,noindex"`
}

type Carb struct {
	LocalTime          string     `json:"label" datastore:"localtime,noindex"`
	TimeValue          TimeValue  `json:"x" datastore:"timestamp"`
	Grams              float32    `json:"carbs" datastore:"grams,noindex"`
	ReferenceReadValue int        `json:"y" datastore:"referenceReadValue,noindex"`
}

type Exercise struct {
	LocalTime         string      `json:"label" datastore:"localtime,noindex"`
	TimeValue         TimeValue   `json:"unixtime" datastore:"timestamp"`
	DurationInMinutes int         `json:"duration" datastore:"duration,noindex"`
	// One of: light, medium, heavy
	Intensity         string      `json:"intensity" datastore:"intensity,noindex"`
}

type FileImportLog struct {
	Id                string
	Md5Checksum       string
	LastDataProcessed time.Time
}

type GlukitUser struct {
	Email              string
	FirstName          string
	LastName           string
	DateOfBirth        time.Time
	DiabetesType       string
	Timezone           string
	LastUpdated        time.Time
	MostRecentRead     time.Time
}

type PointData interface {
	GetTime() time.Time
}

func (read GlucoseRead) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(read.LocalTime, timeutils.TIMEZONE)
	return value
}

func (exercise Exercise) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(exercise.LocalTime, timeutils.TIMEZONE)
	return value
}

func (carb Carb) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(carb.LocalTime, timeutils.TIMEZONE)
	return value
}

func (injection Injection) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(injection.LocalTime, timeutils.TIMEZONE)
	return value
}

type GlucoseReadSlice []GlucoseRead
type ReadStatsSlice []GlucoseRead
type InjectionSlice []Injection
type CarbSlice []Carb
type ExerciseSlice []Exercise
type CoordinateSlice []Coordinate

func (slice GlucoseReadSlice) Len() int {
	return len(slice)
}

func (slice GlucoseReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice GlucoseReadSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice GlucoseReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice GlucoseReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, slice[i].Value, float32(slice[i].Value)}
		dataPoints[i] = dataPoint
	}
	return dataPoints
}

func (slice ReadStatsSlice) Len() int {
	return len(slice)
}

func (slice ReadStatsSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice ReadStatsSlice) Less(i, j int) bool {
	return slice[i].Value < slice[j].Value
}

func (slice ReadStatsSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice InjectionSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, ExtrapolateYValueFromTime(matchingReads, slice[i].TimeValue), slice[i].Units}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}

func (slice CarbSlice) Len() int {
	return len(slice)
}

func (slice CarbSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice CarbSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CarbSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, ExtrapolateYValueFromTime(matchingReads, slice[i].TimeValue), slice[i].Grams}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}

func (slice ExerciseSlice) Len() int {
	return len(slice)
}

func (slice ExerciseSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice ExerciseSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice ExerciseSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, ExtrapolateYValueFromTime(matchingReads, slice[i].TimeValue), float32(slice[i].DurationInMinutes)}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}

func (slice CoordinateSlice) Len() int {
	return len(slice)
}

func (slice CoordinateSlice) Get(i int) int {
	return slice[i].X
}

func (slice CoordinateSlice) Less(i, j int) bool {
	return slice[i].X < slice[j].X
}

func (slice CoordinateSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func ExtrapolateYValueFromTime(reads []GlucoseRead, timeValue TimeValue) (yValue int) {
	lowerIndex := -1
	upperIndex := -1
	for i := range reads {
		if reads[i].TimeValue > timeValue {
			lowerIndex = i - 1;
			upperIndex = i
			break;
		}
	}

	lowerTimeValue := reads[lowerIndex].TimeValue
	upperTimeValue := reads[upperIndex].TimeValue
	lowerYValue := reads[lowerIndex].Value
	upperYValue := reads[upperIndex].Value

	relativeTimePosition := float32((timeValue - lowerTimeValue))/float32((upperTimeValue - lowerTimeValue))
	yValue = int(relativeTimePosition*float32(upperYValue - lowerYValue) + float32(lowerYValue))

	return yValue
}

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

			readTime := time.Unix(startTime.Unix() + offsetInSeconds, 0)
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

func (dayOfReads *DayOfReads) Save(channel chan <- datastore.Property) error {
	defer close(channel)

	size := len(dayOfReads.Reads)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	reads := dayOfReads.Reads
	startTimestamp := int64(reads[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(reads[size - 1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value:  startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value:  endTime,
	}

	for i := range reads {
		//context.Infof("Sending read to channel %s, %d, %d", reads[i].LocalTime, int64(reads[i].TimeValue), int64(reads[i].Value))
		readTimestamp := int64(reads[i].TimeValue)
		readOffset := readTimestamp - startTimestamp
		channel <- datastore.Property {
			Name: strconv.FormatInt(readOffset, 10),
			Value: int64(reads[i].Value),
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

			timestamp := time.Unix(startTime.Unix() + offsetInSeconds, 0)
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

func (dayOfInjections *DayOfInjections) Save(channel chan <- datastore.Property) error {
	defer close(channel)

	size := len(dayOfInjections.Injections)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	injections := dayOfInjections.Injections
	startTimestamp := int64(injections[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(injections[size - 1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value:  startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value:  endTime,
	}

	for i := range injections {
		timestamp := int64(injections[i].TimeValue)
		offset := timestamp - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property {
			Name: strconv.FormatInt(offset, 10),
			Value: float64(injections[i].Units),
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

			timestamp := time.Unix(startTime.Unix() + offsetInSeconds, 0)
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

func (dayOfCarbs *DayOfCarbs) Save(channel chan <- datastore.Property) error {
	defer close(channel)

	size := len(dayOfCarbs.Carbs)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	carbs := dayOfCarbs.Carbs
	startTimestamp := int64(carbs[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(carbs[size - 1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value:  startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value:  endTime,
	}

	for i := range carbs {
		timestamp := int64(carbs[i].TimeValue)
		offset := timestamp - startTimestamp
		// The datastore only supports float64 so we up-cast it to float64
		channel <- datastore.Property {
			Name: strconv.FormatInt(offset, 10),
			Value: float64(carbs[i].Grams),
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

			timestamp := time.Unix(startTime.Unix() + offsetInSeconds, 0)
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

func (dayOfExercises *DayOfExercises) Save(channel chan <- datastore.Property) error {
	defer close(channel)

	size := len(dayOfExercises.Exercises)

	// Nothing to do if the slice has zero elements
	if size == 0 {
		return nil
	}
	exercises := dayOfExercises.Exercises
	startTimestamp := int64(exercises[0].TimeValue)
	startTime := time.Unix(startTimestamp, 0)
	endTimestamp := int64(exercises[size - 1].TimeValue)
	endTime := time.Unix(endTimestamp, 0)

	channel <- datastore.Property{
		Name:  "startTime",
		Value:  startTime,
	}
	channel <- datastore.Property{
		Name:  "endTime",
		Value:  endTime,
	}

	for i := range exercises {
		timestamp := int64(exercises[i].TimeValue)
		offset := timestamp - startTimestamp
		// We need to store two values so we use a string and format each value inside of a single cell value
		channel <- datastore.Property {
			Name: strconv.FormatInt(offset, 10),
			Value: fmt.Sprintf(EXERCISE_VALUE_FORMAT, exercises[i].DurationInMinutes, exercises[i].Intensity),
			NoIndex: true,
		}
	}

	return nil
}
