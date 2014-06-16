package bufio

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
)

type BufferedExerciseBatchWriter struct {
	head      *container.ImmutableList
	size      int
	flushSize int
	wr        glukitio.ExerciseBatchWriter
}

// NewExerciseWriterSize returns a new Writer whose buffer has the specified size.
func NewExerciseWriterSize(wr glukitio.ExerciseBatchWriter, flushSize int) *BufferedExerciseBatchWriter {
	return newExerciseWriterSize(wr, nil, 0, flushSize)
}

func newExerciseWriterSize(wr glukitio.ExerciseBatchWriter, head *container.ImmutableList, size int, flushSize int) *BufferedExerciseBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedExerciseBatchWriter)
	if ok && b.flushSize >= flushSize {
		return b
	}

	w := new(BufferedExerciseBatchWriter)
	w.size = size
	w.flushSize = flushSize
	w.wr = wr
	w.head = head

	return w
}

// WriteExercise writes a single apimodel.DayOfExercises
func (b *BufferedExerciseBatchWriter) WriteExerciseBatch(p []apimodel.Exercise) (glukitio.ExerciseBatchWriter, error) {
	return b.WriteExerciseBatches([]apimodel.DayOfExercises{apimodel.DayOfExercises{p}})
}

// WriteExerciseBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedExerciseBatchWriter) WriteExerciseBatches(p []apimodel.DayOfExercises) (glukitio.ExerciseBatchWriter, error) {
	w := b
	for _, batch := range p {
		if w.size >= w.flushSize {
			fw, err := w.Flush()
			if err != nil {
				return fw, err
			}
			w = fw.(*BufferedExerciseBatchWriter)
		}

		w = newExerciseWriterSize(w.wr, container.NewImmutableList(w.head, batch), w.size+1, w.flushSize)
	}

	return w, nil
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedExerciseBatchWriter) Flush() (glukitio.ExerciseBatchWriter, error) {
	if b.size == 0 {
		return newExerciseWriterSize(b.wr, nil, 0, b.flushSize), nil
	}
	r, size := b.head.ReverseList()
	batch := ListToArrayOfExerciseBatch(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteExerciseBatches(batch)
		if err != nil {
			return nil, err
		}

		return newExerciseWriterSize(innerWriter, nil, 0, b.flushSize), nil
	}

	return newExerciseWriterSize(b.wr, nil, 0, b.flushSize), nil
}

func ListToArrayOfExerciseBatch(head *container.ImmutableList, size int) []apimodel.DayOfExercises {
	r := make([]apimodel.DayOfExercises, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.DayOfExercises)
		cursor = cursor.Next()
	}

	return r
}
