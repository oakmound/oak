//+build windows

package klang

import (
	"errors"

	"github.com/oov/directsound-go/dsound"
)

type dsAudio struct {
	*Encoding
	*dsound.IDirectSoundBuffer
	flags dsound.BufferPlayFlag
}

func (ds *dsAudio) Play() <-chan error {
	ch := make(chan error)
	if ds.Loop {
		ds.flags = dsound.DSBPLAY_LOOPING
	}
	go func(dsbuff *dsound.IDirectSoundBuffer, flags dsound.BufferPlayFlag, ch chan error) {
		err := dsbuff.SetCurrentPosition(0)
		if err != nil {
			select {
			case ch <- err:
			default:
			}
		} else {
			err = dsbuff.Play(0, flags)
			if err != nil {
				select {
				case ch <- err:
				default:
				}
			} else {
				select {
				case ch <- nil:
				default:
				}
			}
		}
	}(ds.IDirectSoundBuffer, ds.flags, ch)
	return ch
}

func (ds *dsAudio) Stop() error {
	err := ds.IDirectSoundBuffer.Stop()
	if err != nil {
		return err
	}
	return ds.IDirectSoundBuffer.SetCurrentPosition(0)
}

// SetVolume uses an underlying directsound command to set
// the volume of the audio. Applies multiplicatively with volume
// filters. Accepts int32s from -10000 to 0, 0 being the max and
// default volume.
func (ds *dsAudio) SetVolume(vol int32) error {
	return ds.IDirectSoundBuffer.SetVolume(vol)
}

func (ds *dsAudio) Filter(fs ...Filter) (Audio, error) {
	var a Audio = ds
	var err, consErr error
	for _, f := range fs {
		a, err = f.Apply(a)
		if err != nil {
			if consErr == nil {
				consErr = err
			} else {
				consErr = errors.New(err.Error() + ":" + consErr.Error())
			}
		}
	}
	// Consider: this is a significant amount
	// of work to do just to make this an in-place filter.
	// would it be worth it to offer both in place and non-inplace
	// filter functions?
	a2, err2 := EncodeBytes(*ds.Encoding)
	if err2 != nil {
		return nil, err2
	}
	// reassign the contents of ds to be that of the
	// new audio, so that this filters in place
	*ds = *a2.(*dsAudio)
	return ds, consErr
}

// MustFilter acts like Filter, but ignores errors (it does not panic,
// as filter errors are expected to be non-fatal)
func (ds *dsAudio) MustFilter(fs ...Filter) Audio {
	a, _ := ds.Filter(fs...)
	return a
}
