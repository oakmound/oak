// Package wav provides functionality to handle .wav files and .wav encoded data
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

// Read returns raw wav data from an input reader
func Read(r io.Reader) (Data, error) {
	wav := Data{}

	err := binary.Read(r, binary.BigEndian, &wav.bChunkID)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.ChunkSize)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.BigEndian, &wav.bFormat)
	if err != nil {
		return wav, err
	}

	err = binary.Read(r, binary.BigEndian, &wav.bSubchunk1ID)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.Subchunk1Size)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.AudioFormat)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.NumChannels)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.SampleRate)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.ByteRate)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.BlockAlign)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.BitsPerSample)
	if err != nil {
		return wav, err
	}

	err = binary.Read(r, binary.BigEndian, &wav.bSubchunk2ID)
	if err != nil {
		return wav, err
	}
	err = binary.Read(r, binary.LittleEndian, &wav.Subchunk2Size)
	if err != nil {
		return wav, err
	}

	wav.Data = make([]byte, wav.Subchunk2Size)
	err = binary.Read(r, binary.LittleEndian, &wav.Data)

	return wav, err
}
