package winaudio

// This file was put together through a combination of
// oov's dsound library (imported),
// this direct sound tutorial page http://www.rastertek.com/dx11tut14.html (author unknown)
// and verdverm's go-wav library (which is copied here as we needed access to a private field)

import (
	"bufio"
	bin "encoding/binary"
	"fmt"
	"github.com/oov/directsound-go/dsound"
	"os"
	"syscall"
)

const (
	SampleRate  = 48000
	Bits        = 16
	Channels    = 2
	BlockAlign  = Channels * Bits / 8
	BytesPerSec = SampleRate * BlockAlign
	NumBlock    = 8
	BlockSize   = (SampleRate / NumBlock) * BlockAlign
)

var (
	user32           *syscall.DLL
	GetDesktopWindow *syscall.Proc
	ds               *dsound.IDirectSound
)

func InitAudio() {
	user32 = syscall.MustLoadDLL("user32.dll")
	GetDesktopWindow = user32.MustFindProc("GetDesktopWindow")
	ds = InitializeDirectSound()
}

func PlayWav(filename string) error {
	// Load a wave audio file onto a secondary buffer.
	dsbuff, err := LoadWaveFile(filename)
	if err != nil {
		panic(err)
	}

	// Play the wave file now that it has been loaded.
	go func(dsbuff *dsound.IDirectSoundBuffer) {
		dsbuff.SetCurrentPosition(0)
		// Play the contents of the secondary sound buffer.
		err := dsbuff.Play(0, 0)
		if err != nil {
			panic(err)
		}
	}(dsbuff)

	return nil
}

func InitializeDirectSound() *dsound.IDirectSound {
	hasDefaultDevice := false
	dsound.DirectSoundEnumerate(func(guid *dsound.GUID, description string, module string) bool {
		if guid == nil {
			hasDefaultDevice = true
			return false
		}
		return true
	})
	if !hasDefaultDevice {
		return nil
	}

	ds, err := dsound.DirectSoundCreate(nil)
	if err != nil {
		panic(err)
	}

	desktopWindow, _, err := GetDesktopWindow.Call()
	err = ds.SetCooperativeLevel(syscall.Handle(desktopWindow), dsound.DSSCL_PRIORITY)
	if err != nil {
		panic(err)
	}

	return ds
}

// The LoadWaveFile function is what handles loading in a .wav audio file and then copies the data onto a new secondary buffer.
// If you are looking to do different formats you would replace this function or write a similar one.

func LoadWaveFile(filename string) (*dsound.IDirectSoundBuffer, error) {

	w := ReadWavData(filename)

	// Set the wave format of secondary buffer that this wave file will be loaded onto.
	wf := dsound.WaveFormatEx{
		FormatTag:      dsound.WAVE_FORMAT_PCM,
		Channels:       Channels,
		SamplesPerSec:  SampleRate,
		BitsPerSample:  Bits,
		BlockAlign:     Channels * Bits / 8,
		AvgBytesPerSec: BytesPerSec,
		ExtSize:        0,
	}

	buffdsc := dsound.BufferDesc{
		// These flags cover everything
		Flags:       dsound.DSBCAPS_GLOBALFOCUS | dsound.DSBCAPS_GETCURRENTPOSITION2 | dsound.DSBCAPS_CTRLVOLUME | dsound.DSBCAPS_CTRLPAN | dsound.DSBCAPS_CTRLFREQUENCY | dsound.DSBCAPS_LOCDEFER,
		Format:      &wf,
		BufferBytes: w.Subchunk2Size,
	}

	// Create the object which stores the wav data in a playable format
	dsbuff, err := ds.CreateSoundBuffer(&buffdsc)
	if err != nil {
		panic(err)
	}

	// Reserve some space in the sound buffer object to write to.
	// The Lock function (and by extension LockBytes) actually
	// reserves two spaces, but we ignore the second.
	by1, by2, err := dsbuff.LockBytes(0, w.Subchunk2Size, 0)
	if err != nil {
		panic(err)
	}

	// Write to the pointer we were given.
	copy(by1, w.data)

	// Update the buffer object with the new data.
	err = dsbuff.UnlockBytes(by1, by2)
	if err != nil {
		panic(err)
	}

	return dsbuff, nil
}

// This is an abbreviated form of the go-wav library, which was copied
// as we required access to the unexported data field

type WavData struct {
	bChunkID  [4]byte // B
	ChunkSize uint32  // L
	bFormat   [4]byte // B

	bSubchunk1ID  [4]byte // B
	Subchunk1Size uint32  // L
	AudioFormat   uint16  // L
	NumChannels   uint16  // L
	SampleRate    uint32  // L
	ByteRate      uint32  // L
	BlockAlign    uint16  // L
	BitsPerSample uint16  // L

	bSubchunk2ID  [4]byte // B
	Subchunk2Size uint32  // L
	data          []byte  // L
}

func ReadWavData(fn string) (wav WavData) {
	ftotal, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("Error opening file", fn)
	}
	file := bufio.NewReader(ftotal)

	bin.Read(file, bin.BigEndian, &wav.bChunkID)
	bin.Read(file, bin.LittleEndian, &wav.ChunkSize)
	bin.Read(file, bin.BigEndian, &wav.bFormat)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk1ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk1Size)
	bin.Read(file, bin.LittleEndian, &wav.AudioFormat)
	bin.Read(file, bin.LittleEndian, &wav.NumChannels)
	bin.Read(file, bin.LittleEndian, &wav.SampleRate)
	bin.Read(file, bin.LittleEndian, &wav.ByteRate)
	bin.Read(file, bin.LittleEndian, &wav.BlockAlign)
	bin.Read(file, bin.LittleEndian, &wav.BitsPerSample)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk2ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk2Size)

	wav.data = make([]byte, wav.Subchunk2Size)
	bin.Read(file, bin.LittleEndian, &wav.data)

	return
}
