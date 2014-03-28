package bufio

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"time"
)

type InjectionStreamer struct {
	buf []model.Injection
	n   int
	wr  glukitio.InjectionBatchWriter
	t   time.Time
	d   time.Duration
	err error
}

// WriteInjection writes a single Injection into the buffer.
func (b *InjectionStreamer) WriteInjection(c model.Injection) (nn int, err error) {
	if b.err != nil {
		return 0, b.err
	}

	p := make([]model.Injection, 1)
	p[0] = c

	return b.WriteInjections(p)
}

// WriteInjections writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *InjectionStreamer) WriteInjections(p []model.Injection) (nn int, err error) {
	// Special case, we don't have a recorded value yet so we start our
	// buffer with the date of the first element
	if b.n == 0 {
		b.resetFirstReadOfBatch(p[0])
	}

	i := 0
	for j, c := range p {
		t := c.GetTime()
		if t.Sub(b.t) >= b.d {
			var n int
			n = copy(b.buf[b.n:], p[i:j])
			b.n += n
			b.Flush()

			nn += n

			// If the flush resulted in an error, abort the write
			if b.err != nil {
				return nn, b.err
			}

			// Move beginning of next batch
			i = j
			b.resetFirstReadOfBatch(p[i])
		}
	}

	n := copy(b.buf[b.n:], p[i:])
	b.n += n
	nn += n

	return nn, nil
}

func (b *InjectionStreamer) resetFirstReadOfBatch(r model.Injection) {
	b.t = r.GetTime().Truncate(b.d)
}

// NewInjectionStreamerDuration returns a new InjectionStreamer whose buffer has the specified size.
func NewInjectionStreamerDuration(wr glukitio.InjectionBatchWriter, bufferLength time.Duration) *InjectionStreamer {
	w := new(InjectionStreamer)
	w.buf = make([]model.Injection, bufferSize)
	w.wr = wr
	w.d = bufferLength

	return w
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *InjectionStreamer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	n, err := b.wr.WriteInjectionBatch(b.buf[0:b.n])
	if n < 1 && err == nil {
		err = glukitio.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *InjectionStreamer) Buffered() int {
	return b.n
}
