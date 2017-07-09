package audio

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"github.com/200sc/klangsynthese/font"
)

// PlayWav is shorthand for GetWav followed by Play.
func PlayWav(f *font.Font, filename string) error {
	ad, err := GetWav(filename)
	if err == nil {
		a := New(f, ad)
		a.Play()
	} else {
		dlog.Error(err)
	}
	return err
}

func DefPlayWav(filename string) error {
	return PlayWav(DefFont, filename)
}
