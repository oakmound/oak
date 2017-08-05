package audio

import (
	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"
	"github.com/200sc/klangsynthese/audio/filter/supports"
	"github.com/oakmound/oak/physics"
)

// SupportsPos is a type used by filters to check that the audio they are given
// has a position.
type SupportsPos interface {
	supports.Encoding
	GetX() *float64
	GetY() *float64
}

var (
	_ audio.Filter = Pos(func(SupportsPos) {})
)

// Pos functions are filters that require a SupportsPos interface
type Pos func(SupportsPos)

// Apply is a function allowing Pos to satisfy the audio.Filter interface.
// Pos applies itself to any audio it is given that supports it.
func (xp Pos) Apply(a audio.Audio) (audio.Audio, error) {
	if sxp, ok := a.(SupportsPos); ok {
		xp(sxp)
		return a, nil
	}
	return a, supports.NewUnsupported([]string{"Pos"})
}

// PosFilter is the only Pos generating function right now. It takes in ears
// to listen from and changes incoming audio to be quiter and panned based
// on positional relation to those ears.
func PosFilter(e *Ears) Pos {
	return func(sp SupportsPos) {
		filter.AssertStereo()(sp)
		x, y := sp.GetX(), sp.GetY()
		if x != nil {
			p := e.CalculatePan(*x)
			filter.Pan(p)(sp)
			if y != nil {
				v := e.CalculateVolume(physics.NewVector(*x, *y))
				filter.Volume(v)(sp)
			}
		}
	}
}
