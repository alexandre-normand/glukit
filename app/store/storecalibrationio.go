package store

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"context"
	"google.golang.org/appengine/datastore"
)

type DataStoreCalibrationBatchWriter struct {
	c context.Context
	k *datastore.Key
}

// NewDataStoreCalibrationBatchWriter creates a new CalibrationBatchWriter that persists to the datastore
func NewDataStoreCalibrationBatchWriter(context context.Context, userProfileKey *datastore.Key) *DataStoreCalibrationBatchWriter {
	w := new(DataStoreCalibrationBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreCalibrationBatchWriter) WriteCalibrationBatches(p []apimodel.DayOfCalibrationReads) (glukitio.CalibrationBatchWriter, error) {
	if _, err := StoreCalibrationReads(w.c, w.k, p); err != nil {
		return w, err
	} else {
		return w, nil
	}
}

func (w *DataStoreCalibrationBatchWriter) WriteCalibrationBatch(p []apimodel.CalibrationRead) (glukitio.CalibrationBatchWriter, error) {
	dayOfCalibrationReads := make([]apimodel.DayOfCalibrationReads, 1)
	dayOfCalibrationReads[0] = apimodel.NewDayOfCalibrationReads(p)
	return w.WriteCalibrationBatches(dayOfCalibrationReads)
}

func (w *DataStoreCalibrationBatchWriter) Flush() (glukitio.CalibrationBatchWriter, error) {
	return w, nil
}
