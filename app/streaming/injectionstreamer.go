package streaming

import (
	"github.com/alexandre-normand/glukit/app/apimodel"
	"github.com/alexandre-normand/glukit/app/container"
	"github.com/alexandre-normand/glukit/app/glukitio"
	"time"
)

type InjectionStreamer struct {
	head      *container.ImmutableList
	startTime *time.Time
	wr        glukitio.InjectionBatchWriter
	d         time.Duration
}

// NewInjectionStreamerDuration returns a new InjectionStreamer whose buffer has the specified size.
func NewInjectionStreamerDuration(wr glukitio.InjectionBatchWriter, bufferDuration time.Duration) *InjectionStreamer {
	return newInjectionStreamerDuration(nil, nil, wr, bufferDuration)
}

func newInjectionStreamerDuration(head *container.ImmutableList, startTime *time.Time, wr glukitio.InjectionBatchWriter, bufferDuration time.Duration) *InjectionStreamer {
	w := new(InjectionStreamer)
	w.head = head
	w.startTime = startTime
	w.wr = wr
	w.d = bufferDuration

	return w
}

// WriteInjection writes a single Injection into the buffer.
func (b *InjectionStreamer) WriteInjection(c apimodel.Injection) (s *InjectionStreamer, err error) {
	return b.WriteInjections([]apimodel.Injection{c})
}

// WriteInjections writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short. p must be sorted by time (oldest to most recent).
func (b *InjectionStreamer) WriteInjections(p []apimodel.Injection) (s *InjectionStreamer, err error) {
	s = newInjectionStreamerDuration(b.head, b.startTime, b.wr, b.d)
	if err != nil {
		return s, err
	}

	for i := range p {
		c := p[i]
		t := c.GetTime()
		truncatedTime := t.Truncate(s.d)

		if s.head == nil {
			s = newInjectionStreamerDuration(container.NewImmutableList(nil, c), &truncatedTime, s.wr, s.d)
		} else if t.Sub(*s.startTime) >= s.d {
			s, err = s.Flush()
			if err != nil {
				return s, err
			}
			s = newInjectionStreamerDuration(container.NewImmutableList(nil, c), &truncatedTime, s.wr, s.d)
		} else {
			s = newInjectionStreamerDuration(container.NewImmutableList(s.head, c), s.startTime, s.wr, s.d)
		}
	}

	return s, err
}

// Flush writes any buffered data to the underlying glukitio.Writer as a batch.
func (b *InjectionStreamer) Flush() (s *InjectionStreamer, err error) {
	r, size := b.head.ReverseList()
	batch := ListToArrayOfInjectionReads(r, size)

	if len(batch) > 0 {
		innerWriter, err := b.wr.WriteInjectionBatch(batch)
		if err != nil {
			return nil, err
		} else {
			return newInjectionStreamerDuration(nil, nil, innerWriter, b.d), nil
		}
	}

	return newInjectionStreamerDuration(nil, nil, b.wr, b.d), nil
}

func ListToArrayOfInjectionReads(head *container.ImmutableList, size int) []apimodel.Injection {
	r := make([]apimodel.Injection, size)
	cursor := head
	for i := 0; i < size; i++ {
		r[i] = cursor.Value().(apimodel.Injection)
		cursor = cursor.Next()
	}

	return r
}

// Close flushes the buffer and the inner writer to effectively ensure nothing is left
// unwritten
func (b *InjectionStreamer) Close() (s *InjectionStreamer, err error) {
	g, err := b.Flush()
	if err != nil {
		return g, err
	}

	innerWriter, err := g.wr.Flush()
	if err != nil {
		return newInjectionStreamerDuration(g.head, g.startTime, innerWriter, b.d), err
	}

	return g, nil
}
