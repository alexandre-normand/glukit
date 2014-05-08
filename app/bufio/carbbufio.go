package bufio

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedCarbBatchWriter struct {
	buf []model.DayOfCarbs
	n   int
	wr  glukitio.CarbBatchWriter
	err error
}

// WriteCarb writes a single model.DayOfCarbs
func (b *BufferedCarbBatchWriter) WriteCarbBatch(p []model.Carb) (nn int, err error) {
	c := make([]model.Carb, len(p))
	if copied := copy(c, p); copied != len(p) {
		return 0, errors.New(fmt.Sprintf("Failed to create copy of carb batch to buffer write, copied [%d] elements but expected [%d]", copied, len(p)))
	}

	n, err := b.WriteCarbBatches([]model.DayOfCarbs{model.DayOfCarbs{c}})
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteCarbBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedCarbBatchWriter) WriteCarbBatches(p []model.DayOfCarbs) (nn int, err error) {
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

// NewCarbWriterSize returns a new Writer whose buffer has the specified size.
func NewCarbWriterSize(wr glukitio.CarbBatchWriter, size int) *BufferedCarbBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedCarbBatchWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedCarbBatchWriter)
	w.buf = make([]model.DayOfCarbs, size)
	w.wr = wr

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedCarbBatchWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteCarbBatches(b.buf[0:b.n])
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
func (b *BufferedCarbBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *BufferedCarbBatchWriter) Buffered() int {
	return b.n
}
