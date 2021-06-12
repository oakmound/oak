// Package mp3 provides functionality to handle .mp3 files and .mp3 encoded data
package mp3

import (
	"bytes"
	"errors"
	"io"

	audio "github.com/oakmound/oak/v3/audio/klang"

	haj "github.com/hajimehoshi/go-mp3"
)

// Load loads an mp3-encoded reader into an audio
func Load(r io.ReadCloser) (audio.Audio, error) {
	d, err := haj.NewDecoder(r)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, d.Length()))
	_, err = io.Copy(buf, d)
	if err != nil {
		return nil, err
	}
	mformat := audio.Format{
		SampleRate: uint32(d.SampleRate()),
		Bits:       16,
		Channels:   2,
	}
	return audio.EncodeBytes(
		audio.Encoding{
			Data:   buf.Bytes(),
			Format: mformat,
		})
}

// Save will eventually save an audio encoded as an MP3 to r
func Save(r io.ReadWriter, a audio.Audio) error {
	return errors.New("Unsupported Functionality")
}
