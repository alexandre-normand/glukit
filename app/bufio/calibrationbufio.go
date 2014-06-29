package bufio

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
)

type BufferedCalibrationBatchWriter struct {
	head      *container.ImmutableList
	size      int
	flushSize int
	wr        glukitio.CalibrationBatchWriter
}

// NewCalibrationWriterSize returns a new Writer whose Buffer has the specified size.
func NewCalibrationWriterSize(wr glukitio.CalibrationBatchWriter, flushSize int) *BufferedCalibrationBatchWriter {
	return newCalibrationWriterSize(wr, nil, 0, flushSize)
}

func newCalibrationWriterSize(wr glukitio.CalibrationBatchWriter, head *container.ImmutableList, size int, flushSize int) *BufferedCalibrationBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedCalibrationBatchWriter)
	if ok && b.flushSize >= flushSize {
		return b
	}

	w := new(BufferedCalibrationBatchWriter)
	w.size = size
	w.flushSize = flushSize
	w.wr = wr
	w.head = head

	return w
}

// WriteCalibration writes a single apimodel.DayOfCalibrationReads
func (b *BufferedCalibrationBatchWriter) WriteCalibrationBatch(p []apimodel.CalibrationRead) (glukitio.CalibrationBatchWriter, error) {
	return b.WriteCalibrationBatches([]apimodel.DayOfCalibrationReads{apimodel.NewDayOfCalibrationReads(p)})
}

// WriteCalibrationBatches writes the contents of p into the buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedCalibrationBatchWriter) WriteCalibrationBatches(p []apimodel.DayOfCalibrationReads) (glukitio.CalibrationBatchWriter, error) {
	w := b
	for i := range p {
		batch := p[i]
		if w.size >= w.flushSize {
			fw, err := w.Flush()
			if err != nil {
				return fw, err
			}
			w = fw.(*BufferedCalibrationBatchWriter)
		}

		w = newCalibrationWriterSize(w.wr, container.NewImmutableList(w.head, batch), w.size+1, w.flushSize)
	}

	return w, nil
}

// Flush writes any buffered data to the underlying glukitio.Writer.
func (b *BufferedCalibrationBatchWriter) Flush() (w glukitio.CalibrationBatchWriter, err error) {
	if b.size == 0 {
		return newCalibrationWriterSize(b.wr, nil, 0, b.flushSize), nil
	}
	r, size := b.head.ReverseList()
	batch := ListToArrayOfCalibrationReadBatch(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteCalibrationBatches(batch)
		if err != nil {
			return nil, err
		}

		return newCalibrationWriterSize(innerWriter, nil, 0, b.flushSize), nil
	}

	return newCalibrationWriterSize(b.wr, nil, 0, b.flushSize), nil
}

func ListToArrayOfCalibrationReadBatch(head *container.ImmutableList, size int) []apimodel.DayOfCalibrationReads {
	r := make([]apimodel.DayOfCalibrationReads, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.DayOfCalibrationReads)
		cursor = cursor.Next()
	}

	return r
}
