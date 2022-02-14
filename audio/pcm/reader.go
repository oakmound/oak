package pcm

import (
	"errors"
	"io"
)

var _ Reader = &LoopingReader{}
var _ Reader = &BytesReader{}
var _ Reader = &IOReader{}

// A Reader mimics io.Reader for pcm data.
type Reader interface {
	Formatted
	ReadPCM(b []byte) (n int, err error)
}

// LoopReader will cache read bytes as they are read and resend them after the reader returns EOF.
func LoopReader(r Reader) Reader {
	return &LoopingReader{Reader: r}
}

// A LoopingReader will read from Reader continually, even after it has been fully consumed. The data read
// from the reader will be cached after read within the LoopingReader structure, potentially inflating memory
// if provided a large stream.
type LoopingReader struct {
	Reader
	buffer     []byte
	bufferPos  int
	eofReached bool
}

func (l *LoopingReader) ReadPCM(p []byte) (n int, err error) {
	if l.eofReached {
		// Note a quirk of this implementation: read calls in succession
		// will not return buffers of a similar size, instead of fully populating
		// the requested p we assert the caller will recall Read resuming from
		// the front of our buffer
		copy(p, l.buffer[l.bufferPos:])
		if len(p) >= len(l.buffer[l.bufferPos:]) {
			n = len(l.buffer[l.bufferPos:])
			l.bufferPos = 0
			return
		}
		n = len(p)
		l.bufferPos += n
		return
	}
	n, err = l.Reader.ReadPCM(p)
	if err != nil {
		if errors.Is(err, io.EOF) {
			l.eofReached = true
			err = nil
		}
	}
	l.buffer = append(l.buffer, p...)
	return
}

// A BytesReader acts like a bytes.Buffer for converting raw []bytes into pcm Readers.
type BytesReader struct {
	Format
	Buffer []byte
	Offset int
}

func (b *BytesReader) ReadPCM(p []byte) (n int, err error) {
	copy(p, b.Buffer)
	if len(p) >= len(b.Buffer[b.Offset:]) {
		b.Offset = len(b.Buffer)
		return len(p), io.EOF
	}
	b.Offset += len(p)
	return len(p), nil
}

// ReadFull acts like io.ReadFull with a pcm Reader. It will read until the provided buffer
// is competely populated by the reader.
func ReadFull(r Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

// ReadAtLeast acts like io.ReadAtLeast with a pcm Reader. It will read until at least min
// bytes have been read into the provided buffer.
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.ReadPCM(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

// An IOReader converts an io.Reader into a pcm.Reader
type IOReader struct {
	Format
	io.Reader
}

func (ior *IOReader) ReadPCM(p []byte) (n int, err error) {
	return ior.Read(p)
}
