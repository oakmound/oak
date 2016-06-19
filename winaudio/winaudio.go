package winaudio

import (
	"bufio"
	bin "encoding/binary"
	"fmt"
	"github.com/oov/directsound-go/dsound"
	//"math/rand"
	"os"
	"syscall"
	"time"
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
	user32           = syscall.MustLoadDLL("user32.dll")
	GetDesktopWindow = user32.MustFindProc("GetDesktopWindow")
)

func Initialize() error {
	// Initialize direct sound and the primary sound buffer.
	ds := InitializeDirectSound()
	//defer ds.Release()

	// Load a wave audio file onto a secondary buffer.
	dsbuff, err := LoadWaveFile("../assets/audio/Glass2.wav", ds)
	if err != nil {
		panic(err)
	}

	// Play the wave file now that it has been loaded.
	go func(dsbuff *dsound.IDirectSoundBuffer) {
		time.Sleep(time.Second)
		PlayWav(dsbuff)
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

func LoadWaveFile(filename string, ds *dsound.IDirectSound) (*dsound.IDirectSoundBuffer, error) {
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
		Flags:       dsound.DSBCAPS_GLOBALFOCUS | dsound.DSBCAPS_GETCURRENTPOSITION2 | dsound.DSBCAPS_CTRLVOLUME | dsound.DSBCAPS_CTRLPAN | dsound.DSBCAPS_CTRLFREQUENCY | dsound.DSBCAPS_LOCDEFER,
		Format:      &wf,
		BufferBytes: uint32(len(w.data)),
	}

	primaryBuf, err := ds.CreateSoundBuffer(&dsound.BufferDesc{
		Flags:       dsound.DSBCAPS_PRIMARYBUFFER,
		BufferBytes: 0,
		Format:      nil,
	})
	if err != nil {
		panic(err)
	}

	err = primaryBuf.SetFormatWaveFormatEx(&dsound.WaveFormatEx{
		FormatTag:      dsound.WAVE_FORMAT_PCM,
		Channels:       Channels,
		SamplesPerSec:  SampleRate,
		BitsPerSample:  Bits,
		BlockAlign:     BlockAlign,
		AvgBytesPerSec: BytesPerSec,
		ExtSize:        0,
	})
	if err != nil {
		panic(err)
	}

	primaryBuf.Release()

	// Create a temporary sound buffer with the specific buffer settings.
	dsbuff, err := ds.CreateSoundBuffer(&buffdsc)
	if err != nil {
		panic(err)
	}

	// Test the buffer format against the direct sound 8 interface
	// unkn, err := dsbuff.QueryInterface(&dsound.GUID{})
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(unkn)

	// Now that the secondary buffer is ready we can load in the wave data from the audio file.
	// First I load it into a memory buffer so I can check and modify the data if I need to.
	// Once the data is in memory you then lock the secondary buffer, copy the data to it using a memcpy, and then unlock it.
	// This secondary buffer is now ready for use.
	// Note that locking the secondary buffer can actually take in two pointers and two positions to write to.
	// This is because it is a circular buffer and if you start by writing to the middle of it you will need the size of the buffer from that point so that you dont write outside the bounds of it.
	// This is useful for streaming audio and such. In this tutorial we create a buffer that is the same size as the audio file and write from the beginning to make things simple.

	by1, by2, err := dsbuff.LockBytes(0, uint32(len(w.data)), 0)
	if err != nil {
		panic(err)
	}

	copy(by1, w.data)

	err = dsbuff.UnlockBytes(by1, by2)
	if err != nil {
		panic(err)
	}

	// is1, is2, err := dsbuff.LockInt16s(0, BytesPerSec, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("LockInt16s Buf1Len:", len(is1), "Buf2Len:", len(is2))

	// // noise fade-in
	// p, ld4 := 0.0, float64(len(is1))
	// for i := range is1 {
	// 	is1[i] = int16((rand.Float64()*10000 - 5000) * (p / ld4))
	// 	p += 1
	// }
	// err = dsbuff.UnlockInt16s(is1, is2)
	// if err != nil {
	// 	panic(err)
	// }

	return dsbuff, nil
}

// The PlayWaveFile function will play the audio file stored in the secondary buffer.
// The moment you use the Play function it will automatically mix the audio into the primary buffer and start it playing if it wasnt already.
// Also note that we set the position to start playing at the beginning of the secondary sound buffer otherwise it will continue from where it last stopped playing.
// And since we set the capabilities of the buffer to allow us to control the sound we set the volume to maximum here.

func PlayWav(dsbuff *dsound.IDirectSoundBuffer) error {
	// err := dsbuff.SetCurrentPosition(0)
	// if err != nil {
	// 	panic(err)
	// 	return err
	// }

	// err = dsbuff.SetVolume(0)
	// if err != nil {
	// 	panic(err)
	// 	return err
	// }

	// Play the contents of the secondary sound buffer.
	err := dsbuff.Play(0, dsound.DSBPLAY_LOOPING)
	if err != nil {
		panic(err)
		return err
	}

	// status, err := dsbuff.GetStatus()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Status:", status)

	// time.Sleep(time.Second)

	return nil
}

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
		fmt.Printf("Error opening\n")
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

	/*
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "ChunkID*: %s\n", ChunkID )
	 *   fmt.Printf( "ChunkSize: %d\n", ChunkSize )
	 *   fmt.Printf( "Format: %s\n", Format )
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "Subchunk1ID: %s\n", Subchunk1ID )
	 *   fmt.Printf( "Subchunk1Size: %d\n", Subchunk1Size )
	 *   fmt.Printf( "AudioFormat: %d\n", AudioFormat )
	 *   fmt.Printf( "NumChannels: %d\n", NumChannels )
	 *   fmt.Printf( "SampleRate: %d\n", SampleRate )
	 *   fmt.Printf( "ByteRate: %d\n", ByteRate )
	 *   fmt.Printf( "BlockAlign: %d\n", BlockAlign )
	 *   fmt.Printf( "BitsPerSample: %d\n", BitsPerSample )
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "Subchunk2ID: %s\n", Subchunk2ID )
	 *   fmt.Printf( "Subchunk2Size: %d\n", Subchunk2Size )
	 *   fmt.Printf( "NumSamples: %d\n", Subchunk2Size / uint32(NumChannels) / uint32(BitsPerSample/8) )
	 *   fmt.Printf( "\ndata: %v\n", len(data) )
	 *   fmt.Printf( "\n\n" )
	 */
	return
}

const (
	mid16 uint16 = 1 >> 2
	big16 uint16 = 1 >> 1
	big32 uint32 = 65535
)

func btou(b []byte) (u []uint16) {
	u = make([]uint16, len(b)/2)
	for i := range u {
		val := uint16(b[i*2])
		val += uint16(b[i*2+1]) << 8
		u[i] = val
	}
	return
}

func btoi16(b []byte) (u []int16) {
	u = make([]int16, len(b)/2)
	for i := range u {
		val := int16(b[i*2])
		val += int16(b[i*2+1]) << 8
		u[i] = val
	}
	return
}

func btof32(b []byte) (f []float32) {
	u := btoi16(b)
	f = make([]float32, len(u))
	for i, v := range u {
		f[i] = float32(v) / float32(32768)
	}
	return
}

func utob(u []uint16) (b []byte) {
	b = make([]byte, len(u)*2)
	for i, val := range u {
		lo := byte(val)
		hi := byte(val >> 8)
		b[i*2] = lo
		b[i*2+1] = hi
	}
	return
}
