package streaming

import (
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type CalibrationReadStreamer struct {
	head    *container.ImmutableList
	tailVal *model.CalibrationRead
	wr      glukitio.CalibrationBatchWriter
	d       time.Duration
}

// NewCalibrationReadStreamerDuration returns a new CalibrationReadStreamer whose buffer has the specified size.
func NewCalibrationReadStreamerDuration(wr glukitio.CalibrationBatchWriter, bufferDuration time.Duration) *CalibrationReadStreamer {
	return newCalibrationStreamerDuration(nil, nil, wr, bufferDuration)
}

func newCalibrationStreamerDuration(head *container.ImmutableList, tailVal *model.CalibrationRead, wr glukitio.CalibrationBatchWriter, bufferDuration time.Duration) *CalibrationReadStreamer {
	w := new(CalibrationReadStreamer)
	w.head = head
	w.tailVal = tailVal
	w.wr = wr
	w.d = bufferDuration

	return w
}

// WriteCalibration writes a single CalibrationRead into the buffer.
func (b *CalibrationReadStreamer) WriteCalibration(c model.CalibrationRead) (s *CalibrationReadStreamer, err error) {
	return b.WriteCalibrations([]model.CalibrationRead{c})
}

// WriteCalibrations writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *CalibrationReadStreamer) WriteCalibrations(p []model.CalibrationRead) (s *CalibrationReadStreamer, err error) {
	s = newCalibrationStreamerDuration(b.head, b.tailVal, b.wr, b.d)
	if err != nil {
		return s, err
	}

	for _, c := range p {
		t := c.GetTime()

		if s.head == nil {
			s = newCalibrationStreamerDuration(container.NewImmutableList(nil, c), &c, s.wr, s.d)
		} else if t.Sub(s.tailVal.GetTime()) >= s.d {
			s, err = s.Flush()
			if err != nil {
				return s, err
			}
			s = newCalibrationStreamerDuration(container.NewImmutableList(nil, c), &c, s.wr, s.d)
		} else {
			s = newCalibrationStreamerDuration(container.NewImmutableList(s.head, c), s.tailVal, s.wr, s.d)
		}
	}

	return s, err
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *CalibrationReadStreamer) Flush() (s *CalibrationReadStreamer, err error) {
	r, size := b.head.ReverseList()
	batch := ListToArrayOfCalibrationReads(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteCalibrationBatch(batch)
		if err != nil {
			return nil, err
		} else {
			return newCalibrationStreamerDuration(nil, nil, innerWriter, b.d), nil
		}
	}

	return newCalibrationStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ListToArrayOfCalibrationReads(head *container.ImmutableList, size int) []model.CalibrationRead {
	r := make([]model.CalibrationRead, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(model.CalibrationRead)
		cursor = cursor.Next()
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *CalibrationReadStreamer) Close() (s *CalibrationReadStreamer, err error) {
	g, err := b.Flush()
	if err != nil {
		return g, err
	}

	innerWriter, err := g.wr.Flush()
	if err != nil {
		return newCalibrationStreamerDuration(g.head, g.tailVal, innerWriter, b.d), err
	}

	return g, nil
}
