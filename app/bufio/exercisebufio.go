package bufio

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedExerciseBatchWriter struct {
	buf []model.DayOfExercises
	n   int
	wr  glukitio.ExerciseBatchWriter
	err error
}

// WriteExercise writes a single model.DayOfExercises
func (b *BufferedExerciseBatchWriter) WriteExerciseBatch(p []model.Exercise) (nn int, err error) {
	DayOfExercises := make([]model.DayOfExercises, 1)
	DayOfExercises[0] = model.DayOfExercises{p}

	n, err := b.WriteExerciseBatches(DayOfExercises)
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteExerciseBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedExerciseBatchWriter) WriteExerciseBatches(p []model.DayOfExercises) (nn int, err error) {
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

// NewExerciseWriterSize returns a new Writer whose buffer has the specified size.
func NewExerciseWriterSize(wr glukitio.ExerciseBatchWriter, size int) *BufferedExerciseBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedExerciseBatchWriter)
	if ok && len(b.buf) >= size {
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedExerciseBatchWriter)
	w.buf = make([]model.DayOfExercises, size)
	w.wr = wr

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedExerciseBatchWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteExerciseBatches(b.buf[0:b.n])
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
func (b *BufferedExerciseBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *BufferedExerciseBatchWriter) Buffered() int {
	return b.n
}