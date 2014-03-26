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

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatches(p []model.DayOfCalibrationReads) (n int, err error) {
	if keys, err := StoreCalibrationReads(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreGlucoseReadBatchWriter) WriteGlucoseReadBatch(p []model.CalibrationRead) (n int, err error) {
	dayOfCalibrationReads := make([]model.DayOfCalibrationReads, 1)
	dayOfCalibrationReads[0] = model.DayOfCalibrationReads{p}
	return w.WriteGlucoseReadBatches(dayOfCalibrationReads)
}

func (w *DataStoreGlucoseReadBatchWriter) Flush() error {
	return nil
}
