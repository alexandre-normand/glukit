/*
Package apimodel provides structs that are used in the client facing API
*/
package apimodel

// Time represents a data point's time
type Time struct {
	Timestamp  int64  `json:"timestamp"`
	TimeZoneId string `json:"timezone"`
}

type Glucose struct {
	Time  Time    `json:"time"`
	Unit  string  `json:"unit"`
	Value float32 `json:"value"`
}

type Calibration struct {
	Time  Time    `json:"time"`
	Unit  string  `json:"unit"`
	Value float32 `json:"value"`
}

// Injection represents an insulin injection
type Injection struct {
	Time        Time    `json:"time"`
	Units       float32 `json:"units"`
	InsulinName string  `json:"insulinName"`
	InsulinType string  `json:"insulinType"`
}

// Meal represents a meal's nutrients
type Meal struct {
	Time          Time    `json:"time"`
	Carbohydrates float32 `json:"carbohydrates"`
	Proteins      float32 `json:"proteins"`
	Fat           float32 `json:"fat"`
	SaturatedFat  float32 `json:"saturatedFat"`
}

// Exercise represents a session of exercise
type Exercise struct {
	Time            Time   `json:"time"`
	DurationMinutes int    `json:"durationInMinutes"`
	Intensity       string `json:"intensity"`
	Description     string `json:"description"`
}
