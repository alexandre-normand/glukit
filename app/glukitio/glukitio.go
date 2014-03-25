package glukitio

import (
	"errors"
	"github.com/alexandre-normand/glukit/app/model"
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
	// TODO get rid of the Flush on the Writer. Only bufio should care about the Flush
	Flush() error
	WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error)
	WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error)
}
