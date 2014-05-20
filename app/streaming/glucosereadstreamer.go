package streaming

import (
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type GlucoseReadStreamer struct {
	head    *container.ImmutableList
	tailVal *model.GlucoseRead
	wr      glukitio.GlucoseReadBatchWriter
	d       time.Duration
}

const (
	BUFFER_SIZE = 86400
)

// WriteGlucoseRead writes a single GlucoseRead into the buffer.
func (b *GlucoseReadStreamer) WriteGlucoseRead(c model.GlucoseRead) (g *GlucoseReadStreamer, err error) {
	return b.WriteGlucoseReads([]model.GlucoseRead{c})
}

// WriteGlucoseReads writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *GlucoseReadStreamer) WriteGlucoseReads(p []model.GlucoseRead) (g *GlucoseReadStreamer, err error) {
	g = newGlucoseStreamerDuration(b.head, b.tailVal, b.wr, b.d)
	if err != nil {
		return g, err
	}

	for _, c := range p {
		t := c.GetTime()

		if g.head == nil {
			g = newGlucoseStreamerDuration(container.NewImmutableList(nil, c), &c, g.wr, g.d)
		} else if t.Sub(g.tailVal.GetTime()) >= g.d {
			g, err = g.Flush()
			if err != nil {
				return g, err
			}
			g = newGlucoseStreamerDuration(container.NewImmutableList(nil, c), &c, g.wr, g.d)
		} else {
			g = newGlucoseStreamerDuration(container.NewImmutableList(g.head, c), g.tailVal, g.wr, g.d)
		}
	}

	return g, err
}

func newGlucoseStreamerDuration(head *container.ImmutableList, tailVal *model.GlucoseRead, wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.head = head
	w.tailVal = tailVal
	w.wr = wr
	w.d = bufferDuration

	return w
}

// NewGlucoseStreamerDuration returns a new GlucoseReadStreamer whose buffer has the specified size.
func NewGlucoseStreamerDuration(wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	return newGlucoseStreamerDuration(nil, nil, wr, bufferDuration)
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *GlucoseReadStreamer) Flush() (*GlucoseReadStreamer, error) {
	r, size := b.head.ReverseList()
	batch := ListToArrayOfGlucoseReads(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteGlucoseReadBatch(batch)
		if err != nil {
			return nil, err
		} else {
			return newGlucoseStreamerDuration(nil, nil, innerWriter, b.d), nil
		}
	}

	return newGlucoseStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ListToArrayOfGlucoseReads(head *container.ImmutableList, size int) []model.GlucoseRead {
	r := make([]model.GlucoseRead, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(model.GlucoseRead)
		cursor = cursor.Next()
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *GlucoseReadStreamer) Close() (*GlucoseReadStreamer, error) {
	g, err := b.Flush()
	if err != nil {
		return g, err
	}

	innerWriter, err := g.wr.Flush()
	if err != nil {
		return newGlucoseStreamerDuration(g.head, g.tailVal, innerWriter, b.d), err
	}

	return g, nil
}
