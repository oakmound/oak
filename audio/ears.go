package audio

import (
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/physics"
)

// ScaleType should be moved to a different package that handles global
// scale varieties
type ScaleType int

const (
	// LINEAR is the only ScaleType right now.
	LINEAR ScaleType = iota
)

// Ears are assisting variables and some position in the game world where
// audio should be 'heard' from, like the player character. Passing in that
// position's x and y as pointers then will allow for sounds further away from
// that point to be quieter and sounds to the left / right of that point to
// be panned left and right.
type Ears struct {
	X             *float64
	Y             *float64
	PanWidth      float64
	SilenceRadius float64
	// VolumeScale and PanScale are currently ignored because there is only
	// one scale type
	VolumeScale ScaleType
	PanScale    ScaleType
}

// NewEars returns a new set of ears to hear pan/volume modified audio from
func NewEars(x, y *float64, panWidth float64, silentRadius float64) *Ears {
	ears := new(Ears)
	ears.X = x
	ears.Y = y
	ears.PanWidth = panWidth
	ears.SilenceRadius = silentRadius
	return ears
}

// CalculatePan converts PanWidth and two x positions into a left / right pan
// value.
func (e *Ears) CalculatePan(x2 float64) float64 {
	v := (x2 - *e.X) / e.PanWidth
	if v < -1 {
		return -1
	} else if v > 1 {
		return 1
	}
	return v
}

// CalculateVolume converts two vector positions and SilenceRadius into a
// volume scale
func (e *Ears) CalculateVolume(v physics.Vector) float64 {
	v2 := physics.NewVector(*e.X, *e.Y)
	dist := v2.Distance(v)

	dlog.Verb("Vector Distance:", dist, v, v2)

	// Ignore scaling variable
	lin := (e.SilenceRadius - dist) / e.SilenceRadius
	if lin < 0 {
		lin = 0
	}

	dlog.Verb("Silence scale", lin, e.SilenceRadius, dist)

	return lin
}
