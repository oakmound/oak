// Package wav provides functionality to handle .wav files and .wav encoded data.
//
// This package may be imported solely to register wavs as a parseable file type within oak:
//
//     import (
//         _ "github.com/oakmound/oak/v4/audio/format/wav"
//     )
//
package wav

import (
	"io"

	"encoding/binary"

	"github.com/oakmound/oak/v4/audio/format"
	"github.com/oakmound/oak/v4/audio/pcm"
)

func init() {
	format.Register(".wav", Load)
}

// Load reads a WAV header from a reader, parsing it's PCM format and returning
// a pcm Reader for the data following the header. It will error if the reader
// does not contain enough data to fill a WAV header. It does not validate that the
// WAV header makes sense.
func Load(r io.Reader) (pcm.Reader, error) {
	data, err := readData(r)
	if err != nil {
		return nil, err
	}

	return &pcm.IOReader{
		Format: pcm.Format{
			SampleRate: data.SampleRate,
			Channels:   data.NumChannels,
			Bits:       data.BitsPerSample,
		},
		Reader: r,
	}, nil
}

// The following is a fork of verdverm's go-wav library

// data stores the raw information contained in a wav file
type data struct {
	bChunkID  [4]byte // B
	ChunkSize uint32  // L
	bFormat   [4]byte // B

	bSubchunk1ID  [4]byte // B
	Subchunk1Size uint32  // L

	AudioFormat   uint16 // L
	NumChannels   uint16 // L
	SampleRate    uint32 // L
	ByteRate      uint32 // L
	BlockAlign    uint16 // L
	BitsPerSample uint16 // L

	bSubchunk2ID  [4]byte // B
	Subchunk2Size uint32  // L
	Data          []byte  // L
}

func readData(r io.Reader) (data, error) {
	data := data{}

	err := binary.Read(r, binary.BigEndian, &data.bChunkID)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.ChunkSize)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.BigEndian, &data.bFormat)
	if err != nil {
		return data, err
	}

	err = binary.Read(r, binary.BigEndian, &data.bSubchunk1ID)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.Subchunk1Size)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.AudioFormat)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.NumChannels)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.SampleRate)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.ByteRate)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.BlockAlign)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.BitsPerSample)
	if err != nil {
		return data, err
	}

	err = binary.Read(r, binary.BigEndian, &data.bSubchunk2ID)
	if err != nil {
		return data, err
	}
	err = binary.Read(r, binary.LittleEndian, &data.Subchunk2Size)
	if err != nil {
		return data, err
	}
	return data, nil
}
