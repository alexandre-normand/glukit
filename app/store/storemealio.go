package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/glukitio"
)

type DataStoreMealBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreMealBatchWriter creates a new MealBatchWriter that persists to the datastore
func NewDataStoreMealBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreMealBatchWriter {
	w := new(DataStoreMealBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreMealBatchWriter) WriteMealBatches(p []apimodel.DayOfMeals) (glukitio.MealBatchWriter, error) {
	if _, err := StoreDaysOfMeals(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreMealBatchWriter) WriteMealBatch(p []apimodel.Meal) (glukitio.MealBatchWriter, error) {
	dayOfMeals := make([]apimodel.DayOfMeals, 1)
	dayOfMeals[0] = apimodel.NewDayOfMeals(p)
	return w.WriteMealBatches(dayOfMeals)
}

func (w *DataStoreMealBatchWriter) Flush() (glukitio.MealBatchWriter, error) {
	return w, nil
}
