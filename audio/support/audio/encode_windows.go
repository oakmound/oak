//+build windows

package audio

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/oov/directsound-go/dsound"
)

var (
	user32           = syscall.NewLazyDLL("user32")
	getDesktopWindow = user32.NewProc("GetDesktopWindow")
	ds               *dsound.IDirectSound
	err              error
)

func init() {
	hasDefaultDevice := false
	dsound.DirectSoundEnumerate(func(guid *dsound.GUID, description string, module string) bool {
		if guid == nil {
			hasDefaultDevice = true
			return false
		}
		return true
	})
	if !hasDefaultDevice {
		ds = nil
		err = errors.New("No default device available to play audio off of")
		return
	}

	ds, err = dsound.DirectSoundCreate(nil)
	if err != nil {
		return
	}
	// We don't check this error because Call() can return
	// "The operation was completed successfully" as an error!
	// Todo: type switch? Do we know the type of "success errors"?
	desktopWindow, _, err := getDesktopWindow.Call()
	if !strings.Contains(err.Error(), "success") {
		fmt.Println("Dsound initialization result:", err)
	}
	err = ds.SetCooperativeLevel(syscall.Handle(desktopWindow), dsound.DSSCL_PRIORITY)
	if err != nil {
		ds = nil
	}
}

// EncodeBytes converts an encoding to Audio
func EncodeBytes(enc Encoding) (Audio, error) {
	// An error here would be an error from init()
	if err != nil {
		return nil, err
	}

	// Create the object which stores the wav data in a playable format
	blockAlign := enc.Channels * enc.Bits / 8
	dsbuff, err := ds.CreateSoundBuffer(&dsound.BufferDesc{
		// These flags cover everything we should ever want to do
		Flags: dsound.DSBCAPS_GLOBALFOCUS | dsound.DSBCAPS_GETCURRENTPOSITION2 | dsound.DSBCAPS_CTRLVOLUME | dsound.DSBCAPS_CTRLPAN | dsound.DSBCAPS_CTRLFREQUENCY | dsound.DSBCAPS_LOCDEFER,
		Format: &dsound.WaveFormatEx{
			FormatTag:      dsound.WAVE_FORMAT_PCM,
			Channels:       enc.Channels,
			SamplesPerSec:  enc.SampleRate,
			BitsPerSample:  enc.Bits,
			BlockAlign:     blockAlign,
			AvgBytesPerSec: enc.SampleRate * uint32(blockAlign),
			ExtSize:        0,
		},
		BufferBytes: uint32(len(enc.Data)),
	})
	if err != nil {
		return nil, err
	}

	// Reserve some space in the sound buffer object to write to.
	// The Lock function (and by extension LockBytes) actually
	// reserves two spaces, but we ignore the second.
	by1, by2, err := dsbuff.LockBytes(0, uint32(len(enc.Data)), 0)
	if err != nil {
		return nil, err
	}

	// Write to the pointer we were given.
	copy(by1, enc.Data)

	// Update the buffer object with the new data.
	err = dsbuff.UnlockBytes(by1, by2)
	if err != nil {
		return nil, err
	}
	return &dsAudio{
		Encoding:           &enc,
		IDirectSoundBuffer: dsbuff,
	}, nil
}
