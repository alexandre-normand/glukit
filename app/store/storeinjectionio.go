package store

import (
	"appengine"
	"appengine/datastore"
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

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatches(p []model.DayOfInjections) (n int, err error) {
	if keys, err := StoreDaysOfInjections(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreInjectionBatchWriter) WriteInjectionBatch(p []model.Injection) (n int, err error) {
	dayOfInjections := make([]model.DayOfInjections, 1)
	dayOfInjections[0] = model.DayOfInjections{p}
	return w.WriteInjectionBatches(dayOfInjections)
}

func (w *DataStoreInjectionBatchWriter) Flush() error {
	return nil
}
