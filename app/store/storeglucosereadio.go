package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/model"
)

type DataStoreGlucoseReadBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreGlucoseReadBatchWriter creates a new GlucoseReadBatchWriter that persists to the datastore
func NewDataStoreGlucoseReadBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreGlucoseReadBatchWriter {
	w := new(DataStoreGlucoseReadBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (n int, err error) {
	if keys, err := StoreDaysOfReads(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatch(p []model.GlucoseRead) (n int, err error) {
	dayOfGlucoseReads := make([]model.DayOfGlucoseReads, 1)
	dayOfGlucoseReads[0] = model.DayOfGlucoseReads{p}
	return w.WriteGlucoseReadBatches(dayOfGlucoseReads)
}

func (w *DataStoreGlucoseReadBatchWriter) Flush() error {
	return nil
}
