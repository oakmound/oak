package audio

import (
	"errors"

	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/font"
)

// Audio is a struct of some audio data and the variables
// required to filter it through a sound font.
type Audio struct {
	*font.Audio
	toStop audio.Audio
	X, Y   *float64
}

// New returns an audio from a font, some audio data, and optional
// positional coordinates
func New(f *font.Font, d Data, coords ...*float64) *Audio {
	a := new(Audio)
	a.Audio = font.NewAudio(f, d)
	if len(coords) > 0 {
		a.X = coords[0]
		if len(coords) > 1 {
			a.Y = coords[1]
		}
	}
	return a
}

// Play begin's an audio's playback
func (a *Audio) Play() <-chan error {
	a2, err := a.Copy()
	if err != nil {
		return errChannel(err)
	}
	a3, err := a2.Filter(a.Font.Filters...)
	if err != nil {
		return errChannel(err)
	}
	a4, err := a3.(*Audio).FullAudio.Copy()
	if err != nil {
		return errChannel(err)
	}
	a.toStop = a4
	return a4.Play()
}

func errChannel(err error) <-chan error {
	ch := make(chan error)
	go func() {
		ch <- err
	}()
	return ch
}

// Stop stops an audio's playback
func (a *Audio) Stop() error {
	if a == nil || a.toStop == nil {
		return errors.New("Nil audio stopped")
	}
	return a.toStop.Stop()
}

// Copy returns a copy of the audio
func (a *Audio) Copy() (audio.Audio, error) {
	a2, err := a.Audio.Copy()
	if err != nil {
		return nil, err
	}
	return New(a.Audio.Font, a2.(audio.FullAudio), a.X, a.Y), nil
}

// MustCopy acts like Copy, but panics on an error.
func (a *Audio) MustCopy() audio.Audio {
	return New(a.Audio.Font, a.Audio.MustCopy().(audio.FullAudio), a.X, a.Y)
}

// Filter returns the audio with some set of filters applied to it.
func (a *Audio) Filter(fs ...audio.Filter) (audio.Audio, error) {
	var ad audio.Audio = a
	var err, consErr error
	for _, f := range fs {
		ad, err = f.Apply(ad)
		if err != nil {
			if consErr == nil {
				consErr = err
			} else {
				consErr = errors.New(err.Error() + ":" + consErr.Error())
			}
		}
	}
	return ad, consErr
}

// MustFilter acts like Filter but ignores errors.
func (a *Audio) MustFilter(fs ...audio.Filter) audio.Audio {
	ad, _ := a.Filter(fs...)
	return ad
}

// GetX returns the X value of where this audio is coming from
func (a *Audio) GetX() *float64 {
	return a.X
}

// GetY returns the Y value of where this audio is coming from
func (a *Audio) GetY() *float64 {
	return a.Y
}

var (
	// Guarantee that Audio can have positional filters applied to it
	_ SupportsPos = &Audio{}
)
