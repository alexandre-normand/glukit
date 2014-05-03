package bufio

import (
	"github.com/alexandre-normand/glukit/app/glukitio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
)

type BufferedGlucoseReadBatchWriter struct {
	buf []model.DayOfGlucoseReads
	n   int
	wr  glukitio.GlucoseReadBatchWriter
	err error
}

// NewGlucoseReadWriterSize returns a new Writer whose Buffer has the specified size.
func NewGlucoseReadWriterSize(wr glukitio.GlucoseReadBatchWriter, size int) *BufferedGlucoseReadBatchWriter {
	// Is it already a Writer?
	b, ok := wr.(*BufferedGlucoseReadBatchWriter)
	if ok && len(b.buf) >= size {
		log.Printf("Already a writer")
		return b
	}

	if size <= 0 {
		size = defaultBufSize
	}
	w := new(BufferedGlucoseReadBatchWriter)
	w.buf = make([]model.DayOfGlucoseReads, size)

	log.Printf("Creating empty Buffer [%v] at address [%p]", w.buf[:w.n], &w.buf)
	w.wr = wr

	return w
}

// WriteGlucose writes a single model.DayOfGlucoseReads
func (b *BufferedGlucoseReadBatchWriter) WriteGlucoseReadBatch(p []model.GlucoseRead) (nn int, err error) {
	log.Printf("Before building of batch, current buffer is:\n%v\n", b.buf[:b.n])
	dayOfGlucoseReads := []model.DayOfGlucoseReads{model.DayOfGlucoseReads{p}}

	log.Printf("Creating batch [%p] of reads with %d elements for data: [%v]", &dayOfGlucoseReads, len(p), dayOfGlucoseReads)
	n, err := b.WriteGlucoseReadBatches(dayOfGlucoseReads)
	if err != nil {
		return n * len(p), err
	}

	return len(p), nil
}

// WriteGlucoseReadBatches writes the contents of p into the Buffer.
// It returns the number of batches written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufferedGlucoseReadBatchWriter) WriteGlucoseReadBatches(p []model.DayOfGlucoseReads) (nn int, err error) {
	log.Printf("Before write, current Buffer is:\n%v\n", b.buf[:b.n])
	log.Printf("Buffer has %d space available with this batch being of length %d", b.Available(), len(p))
	for len(p) > b.Available() && b.err == nil {
		var n int
		log.Printf("Inner copying data at position %d: %v into\n\n-------\n%v", b.n, p, b.buf[:b.n])
		n = copy(b.buf[b.n:], p)
		log.Printf("Inner copied %d items", n)
		b.n += n
		log.Printf("Buffer is [%v]", b.buf[:b.n])
		b.Flush()
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	//log.Printf("Copying data at position %d: %v into\n\n-------\n%v", b.n, p, b.Buf[:b.n])
	log.Printf("Doing the copy starting at index %d, Buffer is [%v]", b.n, b.buf[:b.n])
	n := copy(b.buf[b.n:], p)
	log.Printf("Buffer is [%v]", b.buf[:b.n])
	b.n += n
	nn += n

	log.Printf("Copied %d items from source (%p) to Buffer (%p), current Buffer is:\n%v\n", n, &p, &b.buf, b.buf[:b.n])
	log.Printf("Buffer is [%v]", b.buf[:b.n])

	return nn, nil
}

// Flush writes any Buffered data to the underlying glukitio.Writer.
func (b *BufferedGlucoseReadBatchWriter) Flush() error {
	log.Printf("In flush.. with n=%d", b.n)
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}

	log.Printf("Calling physical persistence with b=%d and data: %v", b.n, b.buf[:b.n])
	n, err := b.wr.WriteGlucoseReadBatches(b.buf[:b.n])
	log.Printf("Buffer is [%v]", b.buf[:b.n])
	if n < b.n && err == nil {
		err = glukitio.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
			log.Printf("Buffer is [%v]", b.buf[:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

// Available returns how many bytes are unused in the Buffer.
func (b *BufferedGlucoseReadBatchWriter) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current Buffer.
func (b *BufferedGlucoseReadBatchWriter) Buffered() int {
	return b.n
}
