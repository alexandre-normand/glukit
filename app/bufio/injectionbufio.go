package bufio

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
)

type BufferedInjectionBatchWriter struct {
	head      *container.ImmutableList
	size      int
	flushSize int
	wr        glukitio.InjectionBatchWriter
}

// NewInjectionWriterSize returns a new Writer whose buffer has the specified size.
func NewInjectionWriterSize(wr glukitio.InjectionBatchWriter, flushSize int) *BufferedInjectionBatchWriter {
	return newInjectionWriterSize(wr, nil, 0, flushSize)
}

func newInjectionWriterSize(wr glukitio.InjectionBatchWriter, head *container.ImmutableList, size int, flushSize int) *BufferedInjectionBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedInjectionBatchWriter)
	if ok && b.flushSize >= flushSize {
		return b
	}

	w := new(BufferedInjectionBatchWriter)
	w.size = size
	w.flushSize = flushSize
	w.wr = wr
	w.head = head

	return w
}

// WriteInjection writes a single apimodel.DayOfInjections
func (b *BufferedInjectionBatchWriter) WriteInjectionBatch(p []apimodel.Injection) (glukitio.InjectionBatchWriter, error) {
	return b.WriteInjectionBatches([]apimodel.DayOfInjections{apimodel.NewDayOfInjections(p)})
}

// WriteInjectionBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedInjectionBatchWriter) WriteInjectionBatches(p []apimodel.DayOfInjections) (glukitio.InjectionBatchWriter, error) {
	w := b
	for _, batch := range p {
		if w.size >= w.flushSize {
			fw, err := w.Flush()
			if err != nil {
				return fw, err
			}
			w = fw.(*BufferedInjectionBatchWriter)
		}

		w = newInjectionWriterSize(w.wr, container.NewImmutableList(w.head, batch), w.size+1, w.flushSize)
	}

	return w, nil
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedInjectionBatchWriter) Flush() (glukitio.InjectionBatchWriter, error) {
	if b.size == 0 {
		return newInjectionWriterSize(b.wr, nil, 0, b.flushSize), nil
	}
	r, size := b.head.ReverseList()
	batch := ListToArrayOfInjectionBatch(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteInjectionBatches(batch)
		if err != nil {
			return nil, err
		}

		return newInjectionWriterSize(innerWriter, nil, 0, b.flushSize), nil
	}

	return newInjectionWriterSize(b.wr, nil, 0, b.flushSize), nil
}

func ListToArrayOfInjectionBatch(head *container.ImmutableList, size int) []apimodel.DayOfInjections {
	r := make([]apimodel.DayOfInjections, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.DayOfInjections)
		cursor = cursor.Next()
	}

	return r
}
