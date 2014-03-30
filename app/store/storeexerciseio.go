package store

import (
	"appengine"
	"appengine/datastore"
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

func (w *DataStoreExerciseBatchWriter) WriteExerciseBatches(p []model.DayOfExercises) (n int, err error) {
	if keys, err := StoreDaysOfExercises(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreExerciseBatchWriter) WriteExerciseBatch(p []model.Exercise) (n int, err error) {
	dayOfExercises := make([]model.DayOfExercises, 1)
	dayOfExercises[0] = model.DayOfExercises{p}
	return w.WriteExerciseBatches(dayOfExercises)
}

func (w *DataStoreExerciseBatchWriter) Flush() error {
	return nil
}
