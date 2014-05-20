package bufio

import (
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
)

type BufferedCarbBatchWriter struct {
	head      *container.ImmutableList
	size      int
	flushSize int
	wr        glukitio.CarbBatchWriter
}

// NewCarbWriterSize returns a new Writer whose buffer has the specified size.
func NewCarbWriterSize(wr glukitio.CarbBatchWriter, flushSize int) *BufferedCarbBatchWriter {
	return newCarbWriterSize(wr, nil, 0, flushSize)
}

func newCarbWriterSize(wr glukitio.CarbBatchWriter, head *container.ImmutableList, size int, flushSize int) *BufferedCarbBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedCarbBatchWriter)
	if ok && b.flushSize >= flushSize {
		return b
	}

	w := new(BufferedCarbBatchWriter)
	w.size = size
	w.flushSize = flushSize
	w.wr = wr
	w.head = head

	return w
}

// WriteCarb writes a single model.DayOfCarbs
func (b *BufferedCarbBatchWriter) WriteCarbBatch(p []model.Carb) (glukitio.CarbBatchWriter, error) {
	return b.WriteCarbBatches([]model.DayOfCarbs{model.DayOfCarbs{p}})
}

// WriteCarbBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedCarbBatchWriter) WriteCarbBatches(p []model.DayOfCarbs) (glukitio.CarbBatchWriter, error) {
	w := b
	for _, batch := range p {
		if w.size >= w.flushSize {
			fw, err := w.Flush()
			if err != nil {
				return fw, err
			}
			w = fw.(*BufferedCarbBatchWriter)
		}

		w = newCarbWriterSize(w.wr, container.NewImmutableList(w.head, batch), w.size+1, w.flushSize)
	}

	return w, nil
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedCarbBatchWriter) Flush() (glukitio.CarbBatchWriter, error) {
	if b.size == 0 {
		return newCarbWriterSize(b.wr, nil, 0, b.flushSize), nil
	}
	r, size := b.head.ReverseList()
	batch := ListToArrayOfCarbBatch(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteCarbBatches(batch)
		if err != nil {
			return nil, err
		}

		return newCarbWriterSize(innerWriter, nil, 0, b.flushSize), nil
	}

	return newCarbWriterSize(b.wr, nil, 0, b.flushSize), nil
}

func ListToArrayOfCarbBatch(head *container.ImmutableList, size int) []model.DayOfCarbs {
	r := make([]model.DayOfCarbs, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(model.DayOfCarbs)
		cursor = cursor.Next()
	}

	return r
}
