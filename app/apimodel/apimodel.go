/*
Package apimodel provides structs that are used in the client facing API
*/
package apimodel

type Glucose struct {
	InternalTime string `json:"InternalTime" xml:"InternalTime,attr"`
	DisplayTime  string `json:"DisplayTime" xml:"DisplayTime,attr"`
	Value        int    `json:"Value" xml:"Value,attr"`
}

type Calibration struct {
	InternalTime string `json:"InternalTime" xml:"InternalTime,attr"`
	DisplayTime  string `json:"DisplayTime" xml:"DisplayTime,attr"`
	Value        int    `json:"Value" xml:"Value,attr"`
}

// Event represents the event structure that holds all events. This includes injections, carbs and exercise.
type Event struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	EventTime    string `xml:"EventTime,attr"`
	EventType    string `xml:"EventType,attr"`
	Description  string `xml:"Decription,attr"`
}

// Injection represents an insulin injection
type Injection struct {
	EventTimestamp
	Units       float32 `json:"units"`
	InsulinName string  `json:"insulinName"`
	InsulinType string  `json:"insulinType"`
}

// DataTime represents hold a timestamp and a localtime string
type EventTimestamp struct {
	RecordedTime string `json:"recordedTime"`
	InternalTime string `json:"internalTime"`
	EventTime    string `json:"eventTime"`
}

// Meal represents a meal's nutrients
type Meal struct {
	EventTimestamp
	Carbohydrates float32 `json:"carbohydrates"`
	Proteins      float32 `json:"proteins"`
	Fat           float32 `json:"fat"`
	SaturatedFat  float32 `json:"saturatedFat"`
}

// Exercise represents a session of exercise
type Exercise struct {
	EventTimestamp
	DurationMinutes float32 `json:"durationInMinutes"`
	Intensity       string  `json:"intensity"`
	Description     string  `json:"description"`
}
