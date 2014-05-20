package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/glukitio"
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

func (w *DataStoreCarbBatchWriter) WriteCarbBatches(p []model.DayOfCarbs) (glukitio.CarbBatchWriter, error) {
	if _, err := StoreDaysOfCarbs(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreCarbBatchWriter) WriteCarbBatch(p []model.Carb) (glukitio.CarbBatchWriter, error) {
	dayOfCarbs := make([]model.DayOfCarbs, 1)
	dayOfCarbs[0] = model.DayOfCarbs{p}
	return w.WriteCarbBatches(dayOfCarbs)
}

func (w *DataStoreCarbBatchWriter) Flush() (glukitio.CarbBatchWriter, error) {
	return w, nil
}
