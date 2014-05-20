package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type DataStoreInjectionBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreInjectionBatchWriter creates a new InjectionBatchWriter that persists to the datastore
func NewDataStoreInjectionBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreInjectionBatchWriter {
	w := new(DataStoreInjectionBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatches(p []model.DayOfInjections) (glukitio.InjectionBatchWriter, error) {
	if _, err := StoreDaysOfInjections(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatch(p []model.Injection) (glukitio.InjectionBatchWriter, error) {
	dayOfInjections := make([]model.DayOfInjections, 1)
	dayOfInjections[0] = model.DayOfInjections{p}
	return w.WriteInjectionBatches(dayOfInjections)
}

func (w *DataStoreInjectionBatchWriter) Flush() (glukitio.InjectionBatchWriter, error) {
	return w, nil
}
