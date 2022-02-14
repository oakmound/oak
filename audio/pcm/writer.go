package pcm

import "io"

// NewWriter returns a writer which can accept audio streamed matching the given format
func NewWriter(f Format) (Writer, error) {
	return newWriter(f)
}

// A Writer can have PCM formatted audio data written to it. It mimics io.Writer.
type Writer interface {
	io.Closer
	Formatted
	// WritePCM expects PCM bytes matching the format this speaker was initialized with.
	// WritePCM will block until all of the bytes are consumed.
	WritePCM([]byte) (n int, err error)
	// Reset must clear out any written data from buffers, without stopping playback
	// TODO: do we need this?
	Reset() error
}
