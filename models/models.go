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
	"appengine"
	"timeutils"
	"log"
)

const (
	UNDEFINED_READ = -1;
)

type TimeValue int64

type DataPoint struct {
	LocalTime string    `json:"label"`
	TimeValue TimeValue `json:"x"`
	Y         int       `json:"y"`
	Value     float32   `json:"value"`
}

type MeterRead struct {
	LocalTime string    `json:"label"`
	TimeValue TimeValue `json:"x"`
	Value     int       `json:"y"`
}

type Injection struct {
	LocalTime          string       `json:"label"`
	TimeValue          TimeValue    `json:"x"`
	Units              float32      `json:"units"`
	ReferenceReadValue int          `json:"y"`
}

type CarbIntake struct {
	LocalTime          string     `json:"label"`
	TimeValue          TimeValue  `json:"x"`
	Grams              float32    `json:"carbs"`
	ReferenceReadValue int        `json:"y"`
}

type Exercise struct {
	LocalTime         string      `json:"label"`
	TimeValue         TimeValue   `json:"unixtime"`
	DurationInMinutes int         `json:"duration"`
	// One of: light, medium, heavy
	Intensity         string      `json:"intensity"`
}

type ReadData struct {
	Email              string
	Name               string
	LastUpdated        time.Time
	ReadsBlobKey       appengine.BlobKey
	InjectionsBlobKey  appengine.BlobKey
	CarbIntakesBlobKey appengine.BlobKey
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
type InjectionSlice []Injection
type CarbIntakeSlice []CarbIntake

func (slice MeterReadSlice) Len() int {
	return len(slice)
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
	yValue = int(relativeTimePosition * float32(upperYValue - lowerYValue) + float32(lowerYValue))

	log.Printf("Extrapolated Y value [%d] from timeValue [%d] which was found between [%d]:[%d] and [%d]:[%d] with respective values [%d] and [%d]", yValue, timeValue, lowerIndex, lowerTimeValue, upperIndex, upperTimeValue, lowerYValue, upperYValue)
	return yValue
}
