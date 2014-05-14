package streaming

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"time"
)

type GlucoseReadStreamer struct {
	head    *GlucoseReadNode
	tailVal *model.GlucoseRead
	wr      glukitio.GlucoseReadBatchWriter
	d       time.Duration
}

type GlucoseReadNode struct {
	next  *GlucoseReadNode
	value model.GlucoseRead
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
		return nil, err
	}

	for _, c := range p {
		t := c.GetTime()

		if g.head == nil {
			g = newGlucoseStreamerDuration(NewGlucoseReadNode(nil, c), &c, b.wr, b.d)
		} else if t.Sub(g.tailVal.GetTime()) >= g.d {
			g, err = g.Flush()
			g = newGlucoseStreamerDuration(NewGlucoseReadNode(nil, c), &c, b.wr, b.d)
		} else {
			g = newGlucoseStreamerDuration(NewGlucoseReadNode(g.head, c), b.tailVal, b.wr, b.d)
		}
	}

	return g, err
}

func newGlucoseStreamerDuration(head *GlucoseReadNode, tailVal *model.GlucoseRead, wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.head = head
	w.tailVal = tailVal
	w.wr = wr
	w.d = bufferDuration

	return w
}

// NewGlucoseStreamerDuration returns a new GlucoseReadStreamer whose buffer has the specified size.
func NewGlucoseStreamerDuration(wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.head = nil
	w.wr = wr
	w.d = bufferDuration

	return w
}

func NewGlucoseReadNode(next *GlucoseReadNode, value model.GlucoseRead) *GlucoseReadNode {
	n := new(GlucoseReadNode)
	n.next = next
	n.value = value

	return n
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *GlucoseReadStreamer) Flush() (*GlucoseReadStreamer, error) {
	batch := ReverseList(b.head)
	log.Printf("Flushing batch %v", batch)
	if len(batch) > 0 {
		n, err := b.wr.WriteGlucoseReadBatch(batch)
		if n < 1 && err == nil {
			err = glukitio.ErrShortWrite
		} else if err != nil {
			return nil, err
		}
	}

	return newGlucoseStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ReverseList(head *GlucoseReadNode) []model.GlucoseRead {
	if head == nil {
		return []model.GlucoseRead{}
	}

	var reverseListHead *GlucoseReadNode = &GlucoseReadNode{nil, head.value}
	size := 0
	for cursor := head; cursor != nil; size, cursor = size+1, cursor.next {
		reverseListHead = NewGlucoseReadNode(reverseListHead, cursor.value)
	}

	r := make([]model.GlucoseRead, size)
	cursor := reverseListHead
	for i := 0; i < size; i++ {
		r[i] = cursor.value
		cursor = cursor.next
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *GlucoseReadStreamer) Close() error {
	_, err := b.Flush()
	if err != nil {
		return err
	}

	return b.wr.Flush()
}
