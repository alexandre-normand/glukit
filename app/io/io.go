package io

// CalibrationWriter is the interface that wraps the basic WriteCalibration method.
//
// WriteCalibration writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
type Writer interface {
	Write(p []*CalibrationRead) (n int, err error)
}
