package glukitio

import (
	"errors"
	"github.com/alexandre-normand/glukit/app/model"
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = errors.New("short write")

// CalibrationBatchWriter is the interface that wraps the basic
// WriteCalibrationBatch and WriteCalibrationBatches methods.
//
// WriteCalibrationBatch writes len(p) model.CalibrationRead from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteCalibrationBatches writes len(p) model.DayOfCalibrationReads from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type CalibrationBatchWriter interface {
	WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error)
	WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error)
}

// GlucoseReadBatchWriter is the interface that wraps the basic
// WriteGlucoseReadBatch and WriteGlucoseReadBatches methods.
//
// WriteGlucoseReadBatch writes len(p) model.CalibrationRead from p to the
// underlying data stream. It returns the number of elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
//
// WriteGlucoseReadBatches writes len(p) model.DayOfCalibrationReads from p to the
// underlying data stream. It returns the number of batch elements written
// from p (0 <= n <= len(p)) and any error encountered that caused the
// write to stop early. Write must return a non-nil error if it returns n < len(p).
type GlucoseReadBatchWriter interface {
	WriteGlucoseReadBatch(p []model.GlucoseRead) (n int, err error)
	WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (n int, err error)
}
