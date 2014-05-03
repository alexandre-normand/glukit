package streaming

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type ExerciseStreamer struct {
	buf []model.Exercise
	n   int
	wr  glukitio.ExerciseBatchWriter
	t   time.Time
	d   time.Duration
	err error
}

// WriteExercise writes a single Exercise into the buffer.
func (b *ExerciseStreamer) WriteExercise(c model.Exercise) (nn int, err error) {
	if b.err != nil {
		return 0, b.err
	}

	p := make([]model.Exercise, 1)
	p[0] = c

	return b.WriteExercises(p)
}

// WriteExercises writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *ExerciseStreamer) WriteExercises(p []model.Exercise) (nn int, err error) {
	// Special case, we don't have a recorded value yet so we start our
	// buffer with the date of the first element
	if b.n == 0 {
		b.resetFirstReadOfBatch(p[0])
	}

	i := 0
	for j, c := range p {
		t := c.GetTime()
		if t.Sub(b.t) >= b.d {
			var n int
			n = copy(b.buf[b.n:], p[i:j])
			b.n += n
			b.Flush()

			nn += n

			// If the flush resulted in an error, abort the write
			if b.err != nil {
				return nn, b.err
			}

			// Move beginning of next batch
			i = j
			b.resetFirstReadOfBatch(p[i])
		}
	}

	n := copy(b.buf[b.n:], p[i:])
	b.n += n
	nn += n

	return nn, nil
}

func (b *ExerciseStreamer) resetFirstReadOfBatch(r model.Exercise) {
	b.t = r.GetTime().Truncate(b.d)
}

// NewExerciseStreamerDuration returns a new ExerciseStreamer whose buffer has the specified size.
func NewExerciseStreamerDuration(wr glukitio.ExerciseBatchWriter, bufferLength time.Duration) *ExerciseStreamer {
	w := new(ExerciseStreamer)
	w.buf = make([]model.Exercise, BUFFER_SIZE)
	w.wr = wr
	w.d = bufferLength

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *ExerciseStreamer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteExerciseBatch(b.buf[0:b.n])
	if n < 1 && err == nil {
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

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *ExerciseStreamer) Close() error {
	err := b.Flush()
	if err != nil {
		return err
	}

	return b.wr.Flush()
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *ExerciseStreamer) Buffered() int {
	return b.n
}
