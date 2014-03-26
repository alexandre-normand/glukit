package store

import (
	"appengine"
	"appengine/datastore"
	"github.com/alexandre-normand/glukit/app/model"
)

type DataStoreCalibrationBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

// NewDataStoreCalibrationBatchWriter creates a new CalibrationBatchWriter that persists to the datastore
func NewDataStoreCalibrationBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreCalibrationBatchWriter {
	w := new(DataStoreCalibrationBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreCalibrationBatchWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error) {
	if keys, err := StoreCalibrationReads(w.c, w.k, p); err != nil {
		return 0, err
	} else {
		return len(keys), nil
	}
}

func (w *DataStoreCalibrationBatchWriter) WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error) {
	dayOfCalibrationReads := make([]model.DayOfCalibrationReads, 1)
	dayOfCalibrationReads[0] = model.DayOfCalibrationReads{p}
	return w.WriteCalibrationBatches(dayOfCalibrationReads)
}

func (w *DataStoreCalibrationBatchWriter) Flush() error {
	return nil
}
