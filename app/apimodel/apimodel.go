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
