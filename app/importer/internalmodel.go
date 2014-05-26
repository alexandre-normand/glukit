package importer

import (
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/util"
	"regexp"
	"strconv"
)

type Glucose struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	Value        string `xml:"Value,attr"`
}

type Calibration struct {
	InternalTime string `json:"InternalTime" xml:"InternalTime,attr"`
	DisplayTime  string `json:"DisplayTime" xml:"DisplayTime,attr"`
	Value        string `json:"Value" xml:"Value,attr"`
}

// Event represents the event structure that holds all events. This includes injections, carbs and exercise.
type Event struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	EventTime    string `xml:"EventTime,attr"`
	EventType    string `xml:"EventType,attr"`
	Description  string `xml:"Decription,attr"`
}

// DataTime represents hold a timestamp and a localtime string
type EventTimestamp struct {
	RecordedTime string `json:"recordedTime"`
	InternalTime string `json:"internalTime"`
	EventTime    string `json:"eventTime"`
}

var mmolValueRegExp = regexp.MustCompile("\\d\\.\\d\\d")

func convertXmlGlucoseRead(read Glucose) (*model.GlucoseRead, error) {
	// Convert display/internal to timestamp with timezone extracted
	if timeUTC, err := util.GetTimeUTC(read.InternalTime); err != nil {
		return nil, err
	} else {
		timestamp := timeUTC.Unix()
		timeLocation := util.GetLocaltimeOffset(read.DisplayTime, timeUTC)

		unit := getUnitFromValue(read.Value)
		if value, err := strconv.ParseFloat(read.Value, 32); err != nil {
			return nil, err
		} else {
			return &model.GlucoseRead{model.Time{timestamp, timeLocation.String()}, unit, float32(value)}, nil
		}
	}
}

func getUnitFromValue(value string) (unit string) {
	unit = model.MG_PER_DL
	if mmolValueRegExp.MatchString(value) {
		unit = model.MMOL_PER_L
	}

	return unit
}

func convertXmlCalibrationRead(calibration Calibration) (*model.CalibrationRead, error) {
	// Convert display/internal to timestamp with timezone extracted
	if timeUTC, err := util.GetTimeUTC(calibration.InternalTime); err != nil {
		return nil, err
	} else {
		timestamp := timeUTC.Unix()
		timeLocation := util.GetLocaltimeOffset(calibration.DisplayTime, timeUTC)

		unit := getUnitFromValue(calibration.Value)
		if value, err := strconv.ParseFloat(calibration.Value, 32); err != nil {
			return nil, err
		} else {
			return &model.CalibrationRead{model.Time{timestamp, timeLocation.String()}, unit, float32(value)}, nil
		}

	}
}
