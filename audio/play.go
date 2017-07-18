package audio

import (
	"github.com/200sc/klangsynthese/font"
	"github.com/oakmound/oak/dlog"
)

// Play is shorthand for Get followed by Play.
func Play(f *font.Font, filename string) error {
	ad, err := Get(filename)
	if err == nil {
		a := New(f, ad)
		a.Play()
	} else {
		dlog.Error(err)
	}
	return err
}

// DefPlay acts like play when given DefFont
func DefPlay(filename string) error {
	return Play(DefFont, filename)
}
