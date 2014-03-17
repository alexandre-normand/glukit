/*
Package io provider buffered io to provide an efficient mecanism to accumulate data prior to physically persisting it.
*/
package bufio

import (
	"github.com/alexandre-normand/glukit/app/io"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
)

const (
	defaultBufSize = 200
)

type CalibrationWriter struct {
	buf []model.CalibrationRead
	n   int
	wr  io.CalibrationWriter
	err error
}

// WriteCalibration writes a single CalibrationRead
func (b *CalibrationWriter) WriteCalibration(calibrationRead model.CalibrationRead) error {
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
func (b *CalibrationWriter) WriteCalibrations(p []model.CalibrationRead) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		n = copy(b.buf[b.n:], p)
		b.n += n
		b.Flush()

		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

// NewWriterSize returns a new Writer whose buffer has the specified
// size.
func NewWriterSize(wr io.CalibrationWriter, size int) *CalibrationWriter {
	// Is it already a Writer?
	b, ok := wr.(*CalibrationWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(CalibrationWriter)
	w.buf = make([]model.CalibrationRead, size)
	w.wr = wr

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
	log.Printf("b is [%v], wr is [%v]", b.buf, b.wr)
	n, err := b.wr.WriteCalibrations(b.buf[0:b.n])
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

// Available returns how many bytes are unused in the buffer.
func (b *CalibrationWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *CalibrationWriter) Buffered() int {
	return b.n
}
