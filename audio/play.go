package audio

import (
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/dlog"
)

// Play is shorthand for Get followed by Play on the DefaultCache.
func Play(f *font.Font, filename string) error {
	ad, err := DefaultCache.Get(filename)
	if err == nil {
		a := New(f, ad)
		a.Play()
	} else {
		dlog.Error(err)
	}
	return err
}

// DefaultPlay acts like play when given DefaultFont
func DefaultPlay(filename string) error {
	return Play(DefaultFont, filename)
}
