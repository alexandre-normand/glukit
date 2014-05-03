package streaming

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"time"
)

type GlucoseReadStreamer struct {
	buf []model.GlucoseRead
	n   int
	wr  glukitio.GlucoseReadBatchWriter
	t   time.Time
	d   time.Duration
	err error
}

const (
	BUFFER_SIZE = 86400
)

// WriteGlucoseRead writes a single GlucoseRead into the buffer.
func (b *GlucoseReadStreamer) WriteGlucoseRead(c model.GlucoseRead) (nn int, err error) {
	if b.err != nil {
		return 0, b.err
	}

	// Special case, we don't have a recorded value yet so we start our
	// buffer with the date of the first element
	if b.n == 0 {
		b.resetFirstReadOfBatch(c)
	}

	t := c.GetTime()
	log.Printf("Streaming value [%v]", c)
	if t.Sub(b.t) >= b.d {
		b.Flush()

		// If the flush resulted in an error, abort the write
		if b.err != nil {
			return 0, b.err
		}

		b.resetFirstReadOfBatch(c)
	}

	b.buf[b.n] = c
	b.n += 1
	//log.Printf("Streamer just added one element and current buffer of writer is %v", b.wr.Buf[:b.wr.N])
	return 1, nil
}

func (b *GlucoseReadStreamer) resetFirstReadOfBatch(r model.GlucoseRead) {
	b.t = r.GetTime().Truncate(b.d)
	log.Printf("First read of batch reset to [%v]", b.t)
}

// NewGlucoseStreamerDuration returns a new GlucoseReadStreamer whose buffer has the specified size.
func NewGlucoseStreamerDuration(wr glukitio.GlucoseReadBatchWriter, bufferLength time.Duration) *GlucoseReadStreamer {
	w := new(GlucoseReadStreamer)
	w.buf = make([]model.GlucoseRead, BUFFER_SIZE)
	log.Printf("streamer buffer is %p", w.buf)
	w.wr = wr
	w.d = bufferLength

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *GlucoseReadStreamer) Flush() error {
	log.Printf("Flushing day of reads with %d reads: %v", b.n, b.buf[:b.n])
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteGlucoseReadBatch(b.buf[:b.n])
	if n < 1 && err == nil {
		err = glukitio.ErrShortWrite
	}
	b.n = 0
	return err
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *GlucoseReadStreamer) Close() error {
	err := b.Flush()
	if err != nil {
		return err
	}

	return b.wr.Flush()
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *GlucoseReadStreamer) Buffered() int {
	return b.n
}
