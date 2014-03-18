package bufio

import (
	"github.com/alexandre-normand/glukit/app/io"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"time"
)

const (
	bufferSize = 86400
)

type TimeBufferedCalibrationWriter struct {
	buf []model.CalibrationRead
	n   int
	wr  io.CalibrationWriter
	t   time.Time
	d   time.Duration
	err error
}

// WriteCalibrations writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *TimeBufferedCalibrationWriter) WriteCalibrations(p []model.CalibrationRead) (nn int, err error) {
	// Special case, we don't have a recorded value yet so we start our
	// buffer with the date of the first element
	if b.t.IsZero() {
		last := p[0]
		b.t = last.GetTime().Truncate(b.d)
		log.Printf("Initialized time to %v", b.t)
	}

	i := 0
	for j, c := 0, p[0]; j < len(p) && b.err == nil; j, c = j+1, p[i] {
		t := c.GetTime()
		if t.Sub(b.t) >= b.d {
			var n int
			n = copy(b.buf[b.n:], p[i:j])
			b.n += n
			b.Flush()

			nn += n
			// Move beginning of next batch
			i = j
		}
	}

	if b.err != nil {
		return nn, b.err
	}

	return nn, nil
}

// NewWriterSize returns a new Writer whose buffer has the specified
// size.
func NewWriterDuration(wr io.CalibrationWriter, bufferLength time.Duration) *TimeBufferedCalibrationWriter {
	// Is it already a Writer?
	b, ok := wr.(*TimeBufferedCalibrationWriter)
	if ok {
		return b
	}

	w := new(TimeBufferedCalibrationWriter)
	w.buf = make([]model.CalibrationRead, bufferSize)
	w.wr = wr
	w.d = bufferLength

	return w
}

// Flush writes any buffered data to the underlying io.Writer.
func (b *TimeBufferedCalibrationWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

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

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *TimeBufferedCalibrationWriter) Buffered() int {
	return b.n
}
