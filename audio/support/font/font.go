// Package font provides utilities to package together audio manipulations as
// a 'font'
package font

import "github.com/oakmound/oak/v3/audio/support/audio"

// Font represents some group of settings which modify how an Audio
// should be played. The name is derived from the concept of a SoundFont
type Font struct {
	Filters []audio.Filter
}

// New returns a *Font.
// It is recommended for future API changes to avoid &Font{} and use NewFont instead
func New() *Font {
	return &Font{}
}

// Filter on a font is applied to all audios as they are played.
// Each call of Filter will completely reset a Font's filters
func (f *Font) Filter(fs ...audio.Filter) *Font {
	f.Filters = fs
	return f
}

// Play on a font is equivalent to Audio.Copy().Filter(Font.GetFilters()).Play()
func (f *Font) Play(a audio.Audio) <-chan error {
	a2, err := a.Copy()
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	_, err = a2.Filter(f.Filters...)
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	return a2.Play()
}
