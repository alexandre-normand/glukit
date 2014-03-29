package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/model"
)

type DataStoreCarbBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreCarbBatchWriter creates a new CarbBatchWriter that persists to the datastore
func NewDataStoreCarbBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreCarbBatchWriter {
	w := new(DataStoreCarbBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreCarbBatchWriter) WriteCarbBatches(p []model.DayOfCarbs) (n int, err error) {
	if keys, err := StoreDaysOfCarbs(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreCarbBatchWriter) WriteCarbBatch(p []model.Carb) (n int, err error) {
	dayOfCarbs := make([]model.DayOfCarbs, 1)
	dayOfCarbs[0] = model.DayOfCarbs{p}
	return w.WriteCarbBatches(dayOfCarbs)
}

func (w *DataStoreCarbBatchWriter) Flush() error {
	return nil
}
