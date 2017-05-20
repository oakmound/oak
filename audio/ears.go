//+build windows

package audio

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/physics"
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

// SetEars sets "Ears" on the given font
// Ears are per Font, so for audio effects to obey ear restrictions they need
// to be built with a Font with those restrictions
// (but multiple fonts can have the same ears)
func (f *Font) SetEars(x, y *float64, panWidth float64, silentRadius float64) {
	ears := new(Ears)
	ears.X = x
	ears.Y = y
	ears.PanWidth = panWidth
	ears.SilenceRadius = silentRadius
	f.Ears = ears
}

// CalculatePan converts PanWidth and two x positions into a left / right pan
// value.
func (e *Ears) CalculatePan(x2 float64) int32 {
	v := (x2 - *e.X) * (winaudio.RIGHT_PAN / e.PanWidth)
	if v < winaudio.LEFT_PAN {
		return winaudio.LEFT_PAN
	} else if v > winaudio.RIGHT_PAN {
		return winaudio.RIGHT_PAN
	}
	dlog.Verb("Pan", *e.X, x2, v)
	return int32(v)
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
