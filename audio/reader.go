package audio

import (
	"errors"
	"io"

	"github.com/oakmound/oak/v3/audio/pcm"
)

var _ pcm.Reader = &LoopingReader{}
var _ pcm.Reader = &BytesReader{}

// LoopReader will cache read bytes as they are read and resend them after the reader returns EOF.
func LoopReader(r pcm.Reader) pcm.Reader {
	return &LoopingReader{Reader: r}
}

// A LoopingReader will read from Reader continually, even after it has been fully consumed. The data read
// from the reader will be cached after read within the LoopingReader structure, potentially inflating memory
// if provided a large stream.
type LoopingReader struct {
	pcm.Reader
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
	pcm.Format
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

func (b *BytesReader) Copy() *BytesReader {
	copyBuff := make([]byte, len(b.Buffer))
	copy(b.Buffer, copyBuff)
	return &BytesReader{
		Format: b.Format,
		Buffer: b.Buffer,
		Offset: b.Offset,
	}
}

// ReadAll will read all of the content within a reader and convert it into a BytesReader. Use carefully; use on
// a LoopingReader or reader which generates its data (e.g. synth types) will likely read until OOM.
func ReadAll(r pcm.Reader) *BytesReader {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.ReadPCM(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return &BytesReader{
		Format: r.PCMFormat(),
		Buffer: b,
	}
}

// ReadFull acts like io.ReadFull with a pcm Reader. It will read until the provided buffer
// is competely populated by the reader.
func ReadFull(r pcm.Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

// ReadAtLeast acts like io.ReadAtLeast with a pcm Reader. It will read until at least min
// bytes have been read into the provided buffer.
func ReadAtLeast(r pcm.Reader, buf []byte, min int) (n int, err error) {
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

