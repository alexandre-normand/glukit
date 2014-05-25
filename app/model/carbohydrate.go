package model

const (
	CARB_TAG = "Meals"
)

// Meal is the data structure that represents a meal of food intake. Only carbohydrates
// are fully supported at the moment.
type Meal struct {
	Time          Time    `json:"time"`
	Mealohydrates float32 `json:"carbohydrates"`
	Proteins      float32 `json:"proteins"`
	Fat           float32 `json:"fat"`
	SaturatedFat  float32 `json:"saturatedFat"`
}

// This holds an array of injections for a whole day
type DayOfMeals struct {
	Meals []Meal
}

type MealSlice []Meal

func (slice MealSlice) Len() int {
	return len(slice)
}

func (slice MealSlice) Less(i, j int) bool {
	return slice[i].Timestamp.EpochTime < slice[j].Timestamp.EpochTime
}

func (slice MealSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice MealSlice) GetEpochTime(i int) (epochTime int64) {
	return slice[i].Timestamp.EpochTime
}

// ToDataPointSlice converts an MealSlice into a generic DataPoint array
func (slice MealSlice) ToDataPointSlice(matchingReads []GlucoseRead) (dataPoints []DataPoint) {
	dataPoints = make([]DataPoint, len(slice))
	for i := range slice {
		dataPoint := DataPoint{slice[i].Timestamp.LocalTime, slice[i].Timestamp.EpochTime,
			linearInterpolateY(matchingReads, slice[i].Timestamp), slice[i].Grams, CARB_TAG}
		dataPoints[i] = dataPoint
	}

	return dataPoints
}
