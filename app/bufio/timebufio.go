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
	wr  io.CalibrationBatchWriter
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
	}

	i := 0
	for j, c := range p {
		t := c.GetTime()
		if t.Sub(b.t) >= b.d {
			var n int
			n = copy(b.buf[b.n:], p[i:j])
			b.n += n
			b.flush()

			nn += n

			// If the flush resulted in an error, abort the write
			if b.err != nil {
				return nn, b.err
			}

			// Move beginning of next batch
			i = j
		}
	}

	n := copy(b.buf[b.n:], p[i:])
	b.n += n
	nn += n

	return nn, nil
}

// NewWriterSize returns a new Writer whose buffer has the specified
// size.
func NewWriterDuration(wr io.CalibrationBatchWriter, bufferLength time.Duration) *TimeBufferedCalibrationWriter {
	w := new(TimeBufferedCalibrationWriter)
	w.buf = make([]model.CalibrationRead, bufferSize)
	w.wr = wr
	w.d = bufferLength

	return w
}

// Flush writes any buffered data to the underlying io.Writer as a batch.
func (b *TimeBufferedCalibrationWriter) flush() error {
	log.Printf("in TimeBufferedCalibrationWriter.flush with buffer of %d, err %v", b.n, b.err)
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteCalibrationBatch(b.buf[0:b.n])
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

// Flush writes any buffered data to the underlying io.Writer.
func (b *TimeBufferedCalibrationWriter) Flush() error {
	log.Printf("User flush with timed buffer of [%d] calibrations", b.n)
	if err := b.flush(); err != nil {
		return err
	}

	return b.wr.Flush()
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *TimeBufferedCalibrationWriter) Buffered() int {
	return b.n
}
