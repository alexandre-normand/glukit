package bufio

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedInjectionBatchWriter struct {
	buf []model.DayOfInjections
	n   int
	wr  glukitio.InjectionBatchWriter
	err error
}

// WriteInjection writes a single model.DayOfInjections
func (b *BufferedInjectionBatchWriter) WriteInjectionBatch(p []model.Injection) (nn int, err error) {
	DayOfInjections := make([]model.DayOfInjections, 1)
	DayOfInjections[0] = model.DayOfInjections{p}

	n, err := b.WriteInjectionBatches(DayOfInjections)
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteInjectionBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedInjectionBatchWriter) WriteInjectionBatches(p []model.DayOfInjections) (nn int, err error) {
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

// NewInjectionWriterSize returns a new Writer whose buffer has the specified size.
func NewInjectionWriterSize(wr glukitio.InjectionBatchWriter, size int) *BufferedInjectionBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedInjectionBatchWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedInjectionBatchWriter)
	w.buf = make([]model.DayOfInjections, size)
	w.wr = wr

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedInjectionBatchWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteInjectionBatches(b.buf[0:b.n])
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
func (b *BufferedInjectionBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *BufferedInjectionBatchWriter) Buffered() int {
	return b.n
}
