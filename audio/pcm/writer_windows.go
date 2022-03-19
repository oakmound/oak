//go:build windows

package pcm

import (
	"fmt"
	"io"
	"sync"

	"github.com/oakmound/oak/v3/audio/wininternal"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oov/directsound-go/dsound"
)

func initOS(driver Driver) error {
	switch driver {
	case DriverDefault:
		fallthrough
	case DriverDirectSound:
		// OK
	default:
		return oakerr.UnsupportedPlatform{
			Operation: "pcm.Init:" + driver.String(),
		}
	}
	cfg, err := wininternal.Init()
	if err != nil {
		return err
	}
	directSoundInterface = cfg.Interface
	return nil
}

var directSoundInterface *dsound.IDirectSound

func newWriter(f Format) (Writer, error) {
	if directSoundInterface == nil {
		return nil, oakerr.NotFound{
			InputName: "directSoundInterface",
		}
	}

	blockAlign := f.Channels * f.Bits / 8
	bufferSize := f.BytesPerSecond() * WriterBufferLengthInSeconds

	dsbuff, err := directSoundInterface.CreateSoundBuffer(&dsound.BufferDesc{
		// These flags cover everything we should ever want to do
		Flags: dsound.DSBCAPS_GLOBALFOCUS | dsound.DSBCAPS_GETCURRENTPOSITION2 | dsound.DSBCAPS_CTRLVOLUME | dsound.DSBCAPS_CTRLPAN | dsound.DSBCAPS_CTRLFREQUENCY | dsound.DSBCAPS_LOCDEFER,
		Format: &dsound.WaveFormatEx{
			FormatTag:      dsound.WAVE_FORMAT_PCM,
			Channels:       f.Channels,
			SamplesPerSec:  f.SampleRate,
			BitsPerSample:  f.Bits,
			BlockAlign:     blockAlign,
			AvgBytesPerSec: f.SampleRate * uint32(blockAlign),
			ExtSize:        0,
		},
		BufferBytes: bufferSize,
	})
	if err != nil {
		return nil, err
	}
	return &directSoundWriter{
		Format:     f,
		buff:       dsbuff,
		bufferSize: bufferSize,
	}, nil
}

type directSoundWriter struct {
	sync.Mutex
	Format
	buff         *dsound.IDirectSoundBuffer
	lockedOffset uint32
	bufferSize   uint32
	playing      bool
}

func (dsw *directSoundWriter) Close() error {
	dsw.Lock()
	defer dsw.Unlock()

	var err error
	if dsw.playing {
		dsw.capOffAudio()
		err = dsw.buff.Stop()
		dsw.playing = false
	}
	dsw.buff.Release()
	return err
}

// this attempts to reduce the amount of noise incurred by stopping a playing audio
// it is not completely effective, a blip can still be heard
func (dsw *directSoundWriter) capOffAudio() {
	// 10k zero bytes
	cap := make([]byte, 10000)
	a, b, err := dsw.buff.LockBytes(dsw.lockedOffset, 10000, 0)
	if err != nil {
		// should not happen, but if it does, we can't proceed
		return
	}
	copy(a, cap)
	if len(b) != 0 {
		copy(b, cap[len(a):])
	}
	dsw.buff.UnlockBytes(a, b)
	dsw.buff.SetCurrentPosition(dsw.lockedOffset)
}

func (dsw *directSoundWriter) Reset() error {
	dsw.Lock()
	defer dsw.Unlock()
	emptyBuff := make([]byte, dsw.bufferSize)
	a, b, err := dsw.buff.LockBytes(0, dsw.bufferSize, 0)
	if err != nil {
		return err
	}
	copy(a, emptyBuff)
	if len(b) != 0 {
		copy(b, emptyBuff)
	}
	err = dsw.buff.UnlockBytes(a, b)
	dsw.Seek(0, io.SeekStart)
	return err
}

func (dsw *directSoundWriter) Seek(offset int64, whence int) (position int64, err error) {
	switch whence {
	case io.SeekStart:
		position = offset
	case io.SeekCurrent:
		current, _, err := dsw.buff.GetCurrentPosition()
		if err != nil {
			return 0, fmt.Errorf("failed to calculate current playing position: %w", err)
		}
		position = int64(current) + offset
	case io.SeekEnd:
		position = int64(dsw.bufferSize) + offset
	default:
		return 0, fmt.Errorf("invalid whence: %v", whence)
	}
	if position < 0 {
		return 0, fmt.Errorf("negative position")
	}
	pos := uint32(position)
	pos %= dsw.bufferSize
	err = dsw.buff.SetCurrentPosition(pos)
	dsw.lockedOffset = pos
	return int64(pos), err
}

func (dsw *directSoundWriter) WritePCM(data []byte) (n int, err error) {
	dsw.Lock()
	defer dsw.Unlock()

	a, b, err := dsw.buff.LockBytes(dsw.lockedOffset, uint32(len(data)), 0)
	if err != nil {
		fmt.Println(dsw.lockedOffset, len(data))
		return 0, fmt.Errorf("failed to lock bytes: %w", err)
	}
	copy(a, data)
	if len(b) != 0 {
		copy(b, data[len(a):])
	}
	err = dsw.buff.UnlockBytes(a, b)
	if err != nil {
		return 0, fmt.Errorf("failed to unlock bytes: %w", err)
	}
	dsw.lockedOffset += uint32(len(a))
	dsw.lockedOffset += uint32(len(b))
	dsw.lockedOffset %= dsw.bufferSize
	if !dsw.playing {
		// Always loop-- these buffers are small, and are continually reused even for
		// larger audio sources
		err = dsw.buff.Play(0, dsound.DSBPLAY_LOOPING)
		if err != nil {
			return len(data), fmt.Errorf("failed to play: %w", err)
		}
	}
	return len(data), nil
}
