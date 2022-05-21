// Package mp3 provides functionality to handle .mp3 files and .mp3 encoded data.
//
// This package may be imported solely to register mp3s as a parseable file type within oak:
//
//     import (
//         _ "github.com/oakmound/oak/v4/audio/format/mp3"
//     )
//
package mp3

import (
	"io"

	"github.com/oakmound/oak/v4/audio/format"
	"github.com/oakmound/oak/v4/audio/pcm"

	"github.com/hajimehoshi/go-mp3"
)

func init() {
	format.Register(".mp3", Load)
}

// Load reads MP3 data from a reader, parsing it's PCM format and returning
// a pcm Reader for the data contained within. It will error if the reader
// does not contain enough data to fill a file header. The resulting format
// will always be 16 bits and 2 channels.
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
