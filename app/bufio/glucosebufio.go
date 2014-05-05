package bufio

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedGlucoseReadBatchWriter struct {
	buf []model.DayOfGlucoseReads
	n   int
	wr  glukitio.GlucoseReadBatchWriter
	err error
}

// NewGlucoseReadWriterSize returns a new Writer whose Buffer has the specified size.
func NewGlucoseReadWriterSize(wr glukitio.GlucoseReadBatchWriter, size int) *BufferedGlucoseReadBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedGlucoseReadBatchWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedGlucoseReadBatchWriter)
	w.buf = make([]model.DayOfGlucoseReads, size)

	w.wr = wr

	return w
}

// WriteGlucose writes a single model.DayOfGlucoseReads
func (b *BufferedGlucoseReadBatchWriter) WriteGlucoseReadBatch(p []model.GlucoseRead) (nn int, err error) {
	readsCopy := make([]model.GlucoseRead, len(p))
	if copied := copy(readsCopy, p); copied != len(p) {
		return 0, errors.New(fmt.Sprintf("Failed to create copy of glucose read batch to buffer write, copied [%d] elements but expected [%d]", copied, len(p)))
	}

	dayOfGlucoseReads := []model.DayOfGlucoseReads{model.DayOfGlucoseReads{readsCopy}}

	n, err := b.WriteGlucoseReadBatches(dayOfGlucoseReads)
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteGlucoseReadBatches writes the contents of p into the Buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedGlucoseReadBatchWriter) WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (nn int, err error) {
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

// Flush writes any Buffered data to the underlying glukitio.Writer.
func (b *BufferedGlucoseReadBatchWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteGlucoseReadBatches(b.buf[:b.n])

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

// Available returns how many bytes are unused in the Buffer.
func (b *BufferedGlucoseReadBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current Buffer.
func (b *BufferedGlucoseReadBatchWriter) Buffered() int {
	return b.n
}
