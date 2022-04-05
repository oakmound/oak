// Package flac provides functionality to handle .flac files and .flac encoded data.
package flac

import (
	"fmt"
	"io"

	"github.com/eaburns/flac"
	"github.com/oakmound/oak/v3/audio/pcm"
)

// Load loads flac data from the incoming reader as an audio
func Load(r io.Reader) (pcm.Reader, error) {
	d, err := flac.NewDecoder(r)
	if err != nil {
		return nil, fmt.Errorf("failed to load flac: %w", err)
	}

	return &pcm.IOReader{
		Format: pcm.Format{
			SampleRate: uint32(d.SampleRate),
			Channels:   uint16(d.NChannels),
			Bits:       uint16(d.BitsPerSample),
		},
		Reader: &reader{d: d},
	}, nil
}

type reader struct {
	d         *flac.Decoder
	readAhead []byte
}

func (r *reader) Read(data []byte) (int, error) {
	if len(r.readAhead) == 0 {
		read, err := r.d.Next()
		if err != nil {
			return 0, err
		}
		r.readAhead = read
	}
	copy(data, r.readAhead)
	if len(r.readAhead) < len(data) {
		n := len(r.readAhead)
		r.readAhead = []byte{}
		return n, nil
	}
	r.readAhead = r.readAhead[len(data):]
	return len(data), nil
}
