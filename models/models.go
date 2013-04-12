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

type ReadData struct {
	Email         string
	Name          string
	LastUpdated   time.Time
	ReadsBlobKey  appengine.BlobKey
}

type MeterReadSlice []MeterRead

func (slice MeterReadSlice) Len() int {
  return len(slice)
}

func (slice MeterReadSlice) Less(i, j int) bool {
  return slice[i].TimeValue < slice[j].TimeValue
}

func (slice MeterReadSlice) Swap(i, j int) {
  slice[i], slice[j] = slice[j], slice[i]
}
