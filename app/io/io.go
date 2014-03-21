package io

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/alexandre-normand/glukit/app/model"
	"github.com/alexandre-normand/glukit/app/store"
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = errors.New("short write")

// CalibrationWriter is the interface that wraps the basic WriteCalibration method.
//
// WriteCalibrationBatch writes len(p) model.DayOfCalibrationReads from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type CalibrationBatchWriter interface {
	Flush() error
	WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error)
	WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error)
}

type DataStoreCalibrationBatchWriter struct {
	c appengine.Context
	k *datastore.Key
}

func NewDataStoreCalibrationBatchWriter(context appengine.Context, userProfileKey *datastore.Key) *DataStoreCalibrationBatchWriter {
	w := new(DataStoreCalibrationBatchWriter)
	w.c = context
	w.k = userProfileKey
	return w
}

func (w *DataStoreCalibrationBatchWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error) {
	if keys, err := store.StoreCalibrationReads(w.c, w.k, p); err != nil {
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
