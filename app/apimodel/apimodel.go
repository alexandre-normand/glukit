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
