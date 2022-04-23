package audio

import (
	"github.com/oakmound/oak/v4/audio/pcm"
)

// NewWriter returns a writer which can accept audio streamed matching the given format
func NewWriter(f pcm.Format) (pcm.Writer, error) {
	return newWriter(f)
}

// MustNewWriter calls NewWriter and panics if an error is returned.
func MustNewWriter(f pcm.Format) pcm.Writer {
	w, err := NewWriter(f)
	if err != nil {
		panic(err)
	}
	return w
}
