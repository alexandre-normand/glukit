package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type DataStoreExerciseBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreExerciseBatchWriter creates a new ExerciseBatchWriter that persists to the datastore
func NewDataStoreExerciseBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreExerciseBatchWriter {
	w := new(DataStoreExerciseBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreExerciseBatchWriter) WriteExerciseBatches(p []model.DayOfExercises) (glukitio.ExerciseBatchWriter, error) {
	if _, err := StoreDaysOfExercises(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreExerciseBatchWriter) WriteExerciseBatch(p []model.Exercise) (glukitio.ExerciseBatchWriter, error) {
	dayOfExercises := make([]model.DayOfExercises, 1)
	dayOfExercises[0] = model.DayOfExercises{p}
	return w.WriteExerciseBatches(dayOfExercises)
}

func (w *DataStoreExerciseBatchWriter) Flush() (glukitio.ExerciseBatchWriter, error) {
	return w, nil
}
