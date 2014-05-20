package streaming

import (
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type CarbStreamer struct {
	head    *container.ImmutableList
	tailVal *model.Carb
	wr      glukitio.CarbBatchWriter
	d       time.Duration
}

// NewCarbStreamerDuration returns a new CarbStreamer whose buffer has the specified size.
func NewCarbStreamerDuration(wr glukitio.CarbBatchWriter, bufferDuration time.Duration) *CarbStreamer {
	return newCarbStreamerDuration(nil, nil, wr, bufferDuration)
}

func newCarbStreamerDuration(head *container.ImmutableList, tailVal *model.Carb, wr glukitio.CarbBatchWriter, bufferDuration time.Duration) *CarbStreamer {
	w := new(CarbStreamer)
	w.head = head
	w.tailVal = tailVal
	w.wr = wr
	w.d = bufferDuration

	return w
}

// WriteCarb writes a single Carb into the buffer.
func (b *CarbStreamer) WriteCarb(c model.Carb) (s *CarbStreamer, err error) {
	return b.WriteCarbs([]model.Carb{c})
}

// WriteCarbs writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *CarbStreamer) WriteCarbs(p []model.Carb) (s *CarbStreamer, err error) {
	s = newCarbStreamerDuration(b.head, b.tailVal, b.wr, b.d)
	if err != nil {
		return s, err
	}

	for _, c := range p {
		t := c.GetTime()

		if s.head == nil {
			s = newCarbStreamerDuration(container.NewImmutableList(nil, c), &c, s.wr, s.d)
		} else if t.Sub(s.tailVal.GetTime()) >= s.d {
			s, err = s.Flush()
			if err != nil {
				return s, err
			}
			s = newCarbStreamerDuration(container.NewImmutableList(nil, c), &c, s.wr, s.d)
		} else {
			s = newCarbStreamerDuration(container.NewImmutableList(s.head, c), s.tailVal, s.wr, s.d)
		}
	}

	return s, err
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *CarbStreamer) Flush() (s *CarbStreamer, err error) {
	r, size := b.head.ReverseList()
	batch := ListToArrayOfCarbReads(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteCarbBatch(batch)
		if err != nil {
			return nil, err
		} else {
			return newCarbStreamerDuration(nil, nil, innerWriter, b.d), nil
		}
	}

	return newCarbStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ListToArrayOfCarbReads(head *container.ImmutableList, size int) []model.Carb {
	r := make([]model.Carb, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(model.Carb)
		cursor = cursor.Next()
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *CarbStreamer) Close() (s *CarbStreamer, err error) {
	g, err := b.Flush()
	if err != nil {
		return g, err
	}

	innerWriter, err := g.wr.Flush()
	if err != nil {
		return newCarbStreamerDuration(g.head, g.tailVal, innerWriter, b.d), err
	}

	return g, nil
}
