package audio

import (
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/dlog"
)

// DefaultFont is the font used for default functions. It can be publicly
// modified to apply a default font to generated audios through def
// methods. If it is not modified, it is a font of zero filters.
var DefaultFont = font.New()

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

// DefaultPlay is shorthand for Play(DefaultFont, filename)
func DefaultPlay(filename string) error {
	return Play(DefaultFont, filename)
}
