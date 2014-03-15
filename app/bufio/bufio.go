/*
Package io provider buffered io to provide an efficient mecanism to accumulate data prior to physically persisting it.
*/
package bufio

import (
	"github.com/alexandre-normand/glukit/app/model"
)

const (
	defaultBufSize = 200
)

type CalibrationWriter struct {
	buf []model.CalibrationRead
}

func (w *CalibrationWriter) WriteCalibration(calibrationRead *model.CalibrationRead) {
	// TODO: implement
}

func (w *CalibrationWriter) WriteCalibrations(calibrationReads []*model.CalibrationRead) {
	// TODO: implement
}

// NewWriterSize returns a new Writer whose buffer has the specified
// size.
func NewWriterSize(size int) *CalibrationWriter {
	if size <= 0 {
		size = defaultBufSize
	}
	w := new(CalibrationWriter)
	w.buf = make([]model.CalibrationRead, size)

	return w
}
