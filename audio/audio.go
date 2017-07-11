package audio

import (
	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter/supports"
	"github.com/200sc/klangsynthese/font"
)

// Audio is a struct of some audio data and the variables
// required to filter it through a sound font.
type Audio struct {
	*font.Audio
	toStop audio.Audio
	X, Y   *float64
}

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

func (a *Audio) Play() <-chan error {
	a2, err := a.Copy()
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	a3, err := a2.Filter(a.Font.Filters...)
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	a4, err := a3.(*Audio).FullAudio.Copy()
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	a.toStop = a4
	return a4.Play()
}

func (a *Audio) Stop() error {
	return a.toStop.Stop()
}

func (a *Audio) Copy() (audio.Audio, error) {
	a2, err := a.Audio.Copy()
	return New(a.Audio.Font, a2.(audio.FullAudio), a.X, a.Y), err
}

func (a *Audio) MustCopy() audio.Audio {
	return New(a.Audio.Font, a.Audio.MustCopy().(audio.FullAudio), a.X, a.Y)
}

func (a *Audio) Filter(fs ...audio.Filter) (audio.Audio, error) {
	var ad audio.Audio = a
	var err error
	var consError supports.ConsError
	for _, f := range fs {
		ad, err = f.Apply(ad)
		if err != nil {
			if consError == nil {
				consError = err.(supports.ConsError)
			} else {
				consError = consError.Cons(err)
			}
		}
	}
	return ad, consError
}

func (a *Audio) MustFilter(fs ...audio.Filter) audio.Audio {
	ad, _ := a.Filter(fs...)
	return ad
}

func (a *Audio) GetX() *float64 {
	return a.X
}

func (a *Audio) GetY() *float64 {
	return a.Y
}

var (
	_ SupportsPos = &Audio{}
)
