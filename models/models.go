/**
 * Created with IntelliJ IDEA.
 * User: skippyjon
 * Date: 2013-04-07
 * Time: 8:26 PM
 * To change this template use File | Settings | File Templates.
 */
package models

import (
	"time"
	"timeutils"
	"log"
)

const (
	UNDEFINED_READ = -1;
)

type TimeValue int64

type TrackingData struct {
	Mean            float64         `json:"mean"`
	Median          float64         `json:"median"`
	Deviation       float64         `json:"deviation"`
	Min      	   float64         `json:"min"`
	Max      	   float64         `json:"max"`
	Distribution    []Coordinate    `json:"distribution"`
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

type MeterRead struct {
	LocalTime string    `json:"label" datastore:"localtime,noindex"`
	TimeValue TimeValue `json:"x" datastore:"timestamp"`
	Value     int       `json:"y" datastore:"value,noindex"`
}

type Injection struct {
	LocalTime          string       `json:"label" datastore:"localtime,noindex"`
	TimeValue          TimeValue    `json:"x" datastore:"timestamp"`
	Units              float32      `json:"units" datastore:"units,noindex"`
	ReferenceReadValue int          `json:"y" datastore:"referenceReadValue,noindex"`
}

type CarbIntake struct {
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

type GlukitUser struct {
	Email              string
	FirstName          string
	LastName           string
	DateOfBirth        time.Time
	DiabetesType       string
	Timezone           string
	LastUpdated        time.Time
}

type PointData interface {
	GetTime() time.Time
}

func (read MeterRead) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(read.LocalTime, timeutils.TIMEZONE)
	return value
}

func (exercise Exercise) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(exercise.LocalTime, timeutils.TIMEZONE)
	return value
}

func (carbIntake CarbIntake) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(carbIntake.LocalTime, timeutils.TIMEZONE)
	return value
}

func (injection Injection) GetTime() (timeValue time.Time) {
	value, _ := timeutils.ParseTime(injection.LocalTime, timeutils.TIMEZONE)
	return value
}

type MeterReadSlice []MeterRead
type ReadStatsSlice []MeterRead
type InjectionSlice []Injection
type CarbIntakeSlice []CarbIntake
type CoordinateSlice []Coordinate

func (slice MeterReadSlice) Len() int {
	return len(slice)
}

func (slice MeterReadSlice) Get(i int) float64 {
	return float64(slice[i].Value)
}

func (slice MeterReadSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice MeterReadSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice MeterReadSlice) ToDataPointSlice() (dataPoints []DataPoint) {
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

func (slice InjectionSlice) ToDataPointSlice(matchingReads []MeterRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, ExtrapolateYValueFromTime(matchingReads, slice[i].TimeValue), slice[i].Units}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}

func (slice CarbIntakeSlice) Len() int {
	return len(slice)
}

func (slice CarbIntakeSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice CarbIntakeSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CarbIntakeSlice) ToDataPointSlice(matchingReads []MeterRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].LocalTime, slice[i].TimeValue, ExtrapolateYValueFromTime(matchingReads, slice[i].TimeValue), slice[i].Grams}
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

func ExtrapolateYValueFromTime(reads []MeterRead, timeValue TimeValue) (yValue int) {
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

	log.Printf("Extrapolated Y value [%d] from timeValue [%d] which was found between [%d]:[%d] and [%d]:[%d] with respective values [%d] and [%d]", yValue, timeValue, lowerIndex, lowerTimeValue, upperIndex, upperTimeValue, lowerYValue, upperYValue)
	return yValue
}
