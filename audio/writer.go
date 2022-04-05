package audio

import (
	"github.com/oakmound/oak/v3/audio/pcm"
)

// NewWriter returns a writer which can accept audio streamed matching the given format
func NewWriter(f pcm.Format) (pcm.Writer, error) {
	return newWriter(f)
}
