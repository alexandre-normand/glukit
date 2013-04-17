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
)

type MeterRead struct {
	LocalTime string   `json:"label"`
	TimeValue int64    `json:"x"`
	Value     int      `json:"y"`
}

type Injection struct {
	LocalTime string   `json:"label"`
	TimeValue int64    `json:"unixtime"`
	Value     float32  `json:"units"`
}

type CarbIntake struct {
	LocalTime string   `json:"label"`
	TimeValue int64    `json:"unixtime"`
	Value     int      `json:"carbs"`
}

type Exercise struct {
	LocalTime         string   `json:"label"`
	TimeValue         int64    `json:"unixtime"`
	DurationInMinutes int      `json:"duration"`
	// One of: light, medium, heavy
	Intensity         string   `json:"intensity"`
}

type ReadData struct {
	Email              string
	Name               string
	LastUpdated        time.Time
	ReadsBlobKey       appengine.BlobKey
	InjectionsBlobKey  appengine.BlobKey
	CarbIntakesBlobKey appengine.BlobKey
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

func (slice InjectionSlice) Len() int {
	return len(slice)
}

func (slice InjectionSlice) Less(i, j int) bool {
	return slice[i].TimeValue < slice[j].TimeValue
}

func (slice InjectionSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
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
