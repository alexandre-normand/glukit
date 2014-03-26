package bufio

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedCalibrationBatchWriter struct {
	buf []model.DayOfCalibrationReads
	n   int
	wr  glukitio.CalibrationBatchWriter
	err error
}

// WriteCalibration writes a single model.DayOfCalibrationReads
func (b *BufferedCalibrationBatchWriter) WriteCalibrationBatch(p []model.CalibrationRead) (nn int, err error) {
	dayOfCalibrationReads := make([]model.DayOfCalibrationReads, 1)
	dayOfCalibrationReads[0] = model.DayOfCalibrationReads{p}

	n, err := b.WriteCalibrationBatches(dayOfCalibrationReads)
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteCalibrationBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedCalibrationBatchWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (nn int, err error) {
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

// NewCalibrationWriterSize returns a new Writer whose buffer has the specified size.
func NewCalibrationWriterSize(wr glukitio.CalibrationBatchWriter, size int) *BufferedCalibrationBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedCalibrationBatchWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedCalibrationBatchWriter)
	w.buf = make([]model.DayOfCalibrationReads, size)
	w.wr = wr

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedCalibrationBatchWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteCalibrationBatches(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = glukitio.ErrShortWrite
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
func (b *BufferedCalibrationBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *BufferedCalibrationBatchWriter) Buffered() int {
	return b.n
}
