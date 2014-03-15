/*
Package io provider buffered io to provide an efficient mecanism to accumulate data prior to physically persisting it.
*/
package bufio

import (
	"github.com/alexandre-normand/glukit/app/io"
	"github.com/alexandre-normand/glukit/app/model"
)

const (
	defaultBufSize = 200
)

type CalibrationWriter struct {
	buf []model.CalibrationRead
	n   int
	wr  io.CalibrationWriter
}

// WriteCalibration writes a single CalibrationRead
func (b *CalibrationWriter) WriteCalibration(calibrationRead *model.CalibrationRead) error {
	if b.err != nil {
		return b.err
	}
	if b.Available() <= 0 && b.Flush() != nil {
		return b.err
	}
	b.buf[b.n] = calibrationRead
	b.n++
	return nil
}

// WriteCalibrations writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *CalibrationWriter) WriteCalibrations(p []*model.CalibrationRead) (nn int, err error) {
	for len(calibrationReads) > b.Available() && b.err == nil {
		var n int
		n = copy(b.buf[b.n:], calibrationReads)
		b.n += n
		b.Flush()

		nn += n
		calibrationReads = calibrationReads[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], calibrationReads)
	b.n += n
	nn += n
	return nn, nil
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

// Flush writes any buffered data to the underlying io.Writer.
func (b *CalibrationWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}
