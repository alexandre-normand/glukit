package streaming

import (
	"errors"
	"fmt"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type GlucoseReadStreamer struct {
	buf []model.GlucoseRead
	wr  glukitio.GlucoseReadBatchWriter
	d   time.Duration
}

const (
	BUFFER_SIZE = 86400
)

// WriteGlucoseRead writes a single GlucoseRead into the buffer.
func (b *GlucoseReadStreamer) WriteGlucoseRead(c model.GlucoseRead) (nn int, g *GlucoseReadStreamer, err error) {
	return b.WriteGlucoseReads([]model.GlucoseRead{c})
}

// WriteGlucoseReads writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *GlucoseReadStreamer) WriteGlucoseReads(p []model.GlucoseRead) (nn int, g *GlucoseReadStreamer, err error) {
	// Special case, we don't have a recorded value yet so we start our
	// buffer with the date of the first element
	var i int
	for j, c := range p {
		t := c.GetTime()

		if len(b.buf) > 0 && t.Sub(b.buf[0].GetTime()) >= b.d {
			var n int
			for _, read := range p[i:j] {
				b.buf = append(b.buf, read)
				n++
			}
			b, err = b.Flush()

			nn += n

			// If the flush resulted in an error, abort the write
			if err != nil {
				return nn, nil, err
			}

			// Move beginning of next batch
			i = j
		}
	}

	var n int
	for _, read := range p[i:] {
		b.buf = append(b.buf, read)
		n++
	}
	nn += n

	g, err = newGlucoseStreamerDuration(b.buf, b.wr, b.d)
	return nn, g, err
}

func newGlucoseStreamerDuration(buf []model.GlucoseRead, wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) (*GlucoseReadStreamer, error) {
	w := new(GlucoseReadStreamer)
	c := make([]model.GlucoseRead, len(buf))
	if copied := copy(c, buf); copied != len(buf) {
		return nil, errors.New(fmt.Sprintf("Failed to create copy of glucose read batch to create new glucose streamer, copied [%d] elements but expected [%d]", copied, len(buf)))
	}
	w.buf = c

	w.wr = wr
	w.d = bufferDuration

	return w, nil
}

// NewGlucoseStreamerDuration returns a new GlucoseReadStreamer whose buffer has the specified size.
func NewGlucoseStreamerDuration(wr glukitio.GlucoseReadBatchWriter, bufferDuration time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.buf = []model.GlucoseRead{}

	w.wr = wr
	w.d = bufferDuration

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *GlucoseReadStreamer) Flush() (*GlucoseReadStreamer, error) {
	if len(b.buf) > 0 {
		n, err := b.wr.WriteGlucoseReadBatch(b.buf)
		if n < 1 && err == nil {
			err = glukitio.ErrShortWrite
		} else if err != nil {
			return nil, err
		}
	}

	g, err := newGlucoseStreamerDuration([]model.GlucoseRead{}, b.wr, b.d)

	return g, err
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

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *GlucoseReadStreamer) Buffered() int {
	return len(b.buf)
}
