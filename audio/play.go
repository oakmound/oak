package audio

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"github.com/200sc/klangsynthese/font"
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

func DefPlay(filename string) error {
	return Play(DefFont, filename)
}
