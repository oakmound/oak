// Package flac provides functionality to handle .flac files and .flac encoded data
package flac

import (
	"errors"
	"io"

	"github.com/eaburns/flac"
	"github.com/oakmound/oak/v3/audio/internal/audio"
)

// def wav format
var format = audio.Format{
	SampleRate: 44100,
	Bits:       16,
	Channels:   2,
}

// Load loads wav data from the incoming reader as an audio
func Load(r io.Reader) (audio.Audio, error) {
	data, meta, err := flac.Decode(r)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load flac")
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

// Save will eventually save an audio encoded as a wav to the given writer
func Save(r io.ReadWriter, a audio.Audio) error {
	return errors.New("Unsupported Functionality")
}

// Format returns the default wav formatting
func Format() audio.Format {
	return format
}
