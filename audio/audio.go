package audio

import (
	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/font"
)

// Audio is a struct of some audio data and the variables
// required to filter it through a sound font.
type Audio struct {
	*font.Audio
	X, Y *float64
}

func New(f *font.Font, aa audio.Audio, coords ...*float64) *Audio {
	a := new(Audio)
	a.Audio = font.NewAudio(f, aa)
	if len(coords) > 0 {
		a.X = coords[0]
		if len(coords) > 1 {
			a.Y = coords[1]
		}
	}
	return a
}

func (a *Audio) GetX() *float64 {
	return a.X
}

func (a *Audio) GetY() *float64 {
	return a.Y
}
