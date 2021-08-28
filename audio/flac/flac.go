// Package flac provides functionality to handle .flac files and .flac encoded data.
package flac

import (
	"fmt"
	"io"

	"github.com/eaburns/flac"
	audio "github.com/oakmound/oak/v3/audio/klang"
)

// def flac format
var format = audio.Format{
	SampleRate: 44100,
	Bits:       16,
	Channels:   2,
}

// Load loads flac data from the incoming reader as an audio
func Load(r io.Reader) (audio.Audio, error) {
	data, meta, err := flac.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to load flac: %w", err)
	}

	fformat := audio.Format{
		SampleRate: uint32(meta.SampleRate),
		Channels:   uint16(meta.NChannels),
		Bits:       uint16(meta.BitsPerSample),
	}
	return audio.EncodeBytes(
		audio.Encoding{
			Data:   data,
			Format: fformat,
		})
}

// Save will eventually save an audio encoded as flac to the given writer
func Save(r io.ReadWriter, a audio.Audio) error {
	return fmt.Errorf("unsupported Functionality")
}

// Format returns the default flac formatting
func Format() audio.Format {
	return format
}
