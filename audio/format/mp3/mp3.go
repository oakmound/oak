// Package mp3 provides functionality to handle .mp3 files and .mp3 encoded data.
package mp3

import (
	"io"

	"github.com/oakmound/oak/v3/audio/format"
	"github.com/oakmound/oak/v3/audio/pcm"

	"github.com/hajimehoshi/go-mp3"
)

func init() {
	format.Register(".mp3", Load)
}

// Load loads an mp3-encoded reader into an audio
func Load(r io.Reader) (pcm.Reader, error) {
	d, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, err
	}
	return &pcm.IOReader{
		Format: pcm.Format{
			SampleRate: uint32(d.SampleRate()),
			Bits:       16,
			Channels:   2,
		},
		Reader: d,
	}, nil
}
