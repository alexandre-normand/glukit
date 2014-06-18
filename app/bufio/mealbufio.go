package bufio

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
)

type BufferedMealBatchWriter struct {
	head      *container.ImmutableList
	size      int
	flushSize int
	wr        glukitio.MealBatchWriter
}

// NewMealWriterSize returns a new Writer whose buffer has the specified size.
func NewMealWriterSize(wr glukitio.MealBatchWriter, flushSize int) *BufferedMealBatchWriter {
	return newMealWriterSize(wr, nil, 0, flushSize)
}

func newMealWriterSize(wr glukitio.MealBatchWriter, head *container.ImmutableList, size int, flushSize int) *BufferedMealBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedMealBatchWriter)
	if ok && b.flushSize >= flushSize {
		return b
	}

	w := new(BufferedMealBatchWriter)
	w.size = size
	w.flushSize = flushSize
	w.wr = wr
	w.head = head

	return w
}

// WriteMeal writes a single apimodel.DayOfMeals
func (b *BufferedMealBatchWriter) WriteMealBatch(p []apimodel.Meal) (glukitio.MealBatchWriter, error) {
	return b.WriteMealBatches([]apimodel.DayOfMeals{apimodel.NewDayOfMeals(p)})
}

// WriteMealBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedMealBatchWriter) WriteMealBatches(p []apimodel.DayOfMeals) (glukitio.MealBatchWriter, error) {
	w := b
	for _, batch := range p {
		if w.size >= w.flushSize {
			fw, err := w.Flush()
			if err != nil {
				return fw, err
			}
			w = fw.(*BufferedMealBatchWriter)
		}

		w = newMealWriterSize(w.wr, container.NewImmutableList(w.head, batch), w.size+1, w.flushSize)
	}

	return w, nil
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedMealBatchWriter) Flush() (glukitio.MealBatchWriter, error) {
	if b.size == 0 {
		return newMealWriterSize(b.wr, nil, 0, b.flushSize), nil
	}
	r, size := b.head.ReverseList()
	batch := ListToArrayOfMealBatch(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteMealBatches(batch)
		if err != nil {
			return nil, err
		}

		return newMealWriterSize(innerWriter, nil, 0, b.flushSize), nil
	}

	return newMealWriterSize(b.wr, nil, 0, b.flushSize), nil
}

func ListToArrayOfMealBatch(head *container.ImmutableList, size int) []apimodel.DayOfMeals {
	r := make([]apimodel.DayOfMeals, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.DayOfMeals)
		cursor = cursor.Next()
	}

	return r
}
