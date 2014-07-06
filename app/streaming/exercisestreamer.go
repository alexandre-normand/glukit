package streaming

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"time"
)

type ExerciseStreamer struct {
	head      *container.ImmutableList
	startTime *time.Time
	wr        glukitio.ExerciseBatchWriter
	d         time.Duration
}

// NewExerciseStreamerDuration returns a new ExerciseStreamer whose buffer has the specified size.
func NewExerciseStreamerDuration(wr glukitio.ExerciseBatchWriter, bufferDuration time.Duration) *ExerciseStreamer {
	return newExerciseStreamerDuration(nil, nil, wr, bufferDuration)
}

func newExerciseStreamerDuration(head *container.ImmutableList, startTime *time.Time, wr glukitio.ExerciseBatchWriter, bufferDuration time.Duration) *ExerciseStreamer {
	w := new(ExerciseStreamer)
	w.head = head
	w.startTime = startTime
	w.wr = wr
	w.d = bufferDuration

	return w
}

// WriteExercise writes a single Exercise into the buffer.
func (b *ExerciseStreamer) WriteExercise(c apimodel.Exercise) (s *ExerciseStreamer, err error) {
	return b.WriteExercises([]apimodel.Exercise{c})
}

// WriteExercises writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *ExerciseStreamer) WriteExercises(p []apimodel.Exercise) (s *ExerciseStreamer, err error) {
	s = newExerciseStreamerDuration(b.head, b.startTime, b.wr, b.d)
	if err != nil {
		return s, err
	}

	for i := range p {
		c := p[i]
		t := c.GetTime()
		truncatedTime := t.Truncate(s.d)

		if s.head == nil {
			s = newExerciseStreamerDuration(container.NewImmutableList(nil, c), &truncatedTime, s.wr, s.d)
		} else if t.Sub(*s.startTime) >= s.d {
			s, err = s.Flush()
			if err != nil {
				return s, err
			}
			s = newExerciseStreamerDuration(container.NewImmutableList(nil, c), &truncatedTime, s.wr, s.d)
		} else {
			s = newExerciseStreamerDuration(container.NewImmutableList(s.head, c), s.startTime, s.wr, s.d)
		}
	}

	return s, err
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *ExerciseStreamer) Flush() (s *ExerciseStreamer, err error) {
	r, size := b.head.ReverseList()
	batch := ListToArrayOfExerciseReads(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteExerciseBatch(batch)
		if err != nil {
			return nil, err
		} else {
			return newExerciseStreamerDuration(nil, nil, innerWriter, b.d), nil
		}
	}

	return newExerciseStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ListToArrayOfExerciseReads(head *container.ImmutableList, size int) []apimodel.Exercise {
	r := make([]apimodel.Exercise, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.Exercise)
		cursor = cursor.Next()
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *ExerciseStreamer) Close() (s *ExerciseStreamer, err error) {
	g, err := b.Flush()
	if err != nil {
		return g, err
	}

	innerWriter, err := g.wr.Flush()
	if err != nil {
		return newExerciseStreamerDuration(g.head, g.startTime, innerWriter, b.d), err
	}

	return newExerciseStreamerDuration(nil, nil, innerWriter, g.d), nil
}
