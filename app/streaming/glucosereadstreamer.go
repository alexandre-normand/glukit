package streaming

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"time"
)

type GlucoseReadStreamer struct {
	head *GlucoseReadNode
	wr   glukitio.GlucoseReadBatchWriter
	d    time.Duration
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
	g = newGlucoseStreamerDuration(b.head, b.wr, b.d)
	if err != nil {
		return nil, err
	}

	for _, c := range p {
		t := c.GetTime()

		if g.head == nil {
			g.head = &GlucoseReadNode{nil, c}
		} else if t.Sub(g.head.value.GetTime()) >= g.d {
			g, err = g.Flush()
		} else {
			g.head = &GlucoseReadNode{g.head, c}
		}
	}

	log.Printf("Done Writing reads, state is %p: %v", g.head, g.head)
	return g, err
}

func newGlucoseStreamerDuration(head *GlucoseReadNode, wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.head = head
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

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *GlucoseReadStreamer) Flush() (*GlucoseReadStreamer, error) {
	batch := reverseList(b.head)

	if len(batch) > 0 {
		n, err := b.wr.WriteGlucoseReadBatch(batch)
		if n < 1 && err == nil {
			err = glukitio.ErrShortWrite
		} else if err != nil {
			return nil, err
		}
	}

	return newGlucoseStreamerDuration(nil, b.wr, b.d), nil
}

func reverseList(head *GlucoseReadNode) []model.GlucoseRead {
	if head == nil {
		return []model.GlucoseRead{}
	}

	var reverseListHead *GlucoseReadNode = &GlucoseReadNode{nil, head.value}
	size := 0
	for cursor := head; cursor != nil; size, cursor = size+1, cursor.next {
		log.Printf("Current item is %p: %v, head is %p", cursor, cursor, reverseListHead)
		reverseListHead = &GlucoseReadNode{reverseListHead, cursor.value}
		log.Printf("Current head is %p", reverseListHead)
	}

	r := make([]model.GlucoseRead, size)
	cursor := reverseListHead
	for i := 0; i < size; i++ {
		log.Printf("Building array with %p: %v", cursor, cursor)
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
