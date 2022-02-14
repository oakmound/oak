// Package wav provides functionality to handle .wav files and .wav encoded data.
package wav

import (
	"errors"
	"io"

	"encoding/binary"

	audio "github.com/oakmound/oak/v3/audio/klang"
)

// Load loads wav data from the incoming reader as an audio
func Load(r io.Reader) (audio.Audio, error) {
	wav, err := Read(r)
	if err != nil {
		return nil, err
	}
	return audio.EncodeBytes(
		audio.Encoding{
			Data: wav.Data,
			Format: audio.Format{
				SampleRate: wav.SampleRate,
				Channels:   wav.NumChannels,
				Bits:       wav.BitsPerSample,
			},
		})
}

// Save will eventually save an audio encoded as a wav to the given writer
func Save(r io.ReadWriter, a audio.Audio) error {
	return errors.New("Unsupported Functionality")
}

// Read reads a WAV header from the provided reader, returning the PCM format and
// leaving the PCM data in the reader available for consumption.
func ReadFormat(r io.Reader) (audio.Format, error) {
	data, err := ReadData(r)
	if err != nil {
		return audio.Format{}, err
	}

	return audio.Format{
		SampleRate: data.SampleRate,
		Channels:   data.NumChannels,
		Bits:       data.BitsPerSample,
	}, nil
}

// The following is a "fork" of verdverm's go-wav library

// Data stores the raw information contained in a wav file
type Data struct {
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

func ReadData(r io.Reader) (Data, error) {
	data := Data{}

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

// Read returns raw wav data from an input reader
func Read(r io.Reader) (Data, error) {
	data, err := ReadData(r)
	if err != nil {
		return data, err
	}

	data.Data = make([]byte, data.Subchunk2Size)
	err = binary.Read(r, binary.LittleEndian, &data.Data)

	return data, err
}
