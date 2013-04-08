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
