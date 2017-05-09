//+build windows

package audio

import (
	"fmt"

	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

type ScaleType int

const (
	LINEAR ScaleType = iota
)

type Ears struct {
	X             *float64
	Y             *float64
	PanWidth      float64
	SilenceRadius float64
	VolumeScale   ScaleType
	PanScale      ScaleType
}

// For Pan and volume calculation
func (f *Font) SetEars(x, y *float64, panWidth float64, silentRadius float64) {
	ears := new(Ears)
	ears.X = x
	ears.Y = y
	ears.PanWidth = panWidth
	ears.SilenceRadius = silentRadius
	f.Ears = ears
}

func (e *Ears) CalculatePan(x2 float64) int32 {
	v := (x2 - *e.X) * (winaudio.RIGHT_PAN / e.PanWidth)
	if v < winaudio.LEFT_PAN {
		return winaudio.LEFT_PAN
	} else if v > winaudio.RIGHT_PAN {
		return winaudio.RIGHT_PAN
	}
	fmt.Println("Pan", *e.X, x2, v)
	return int32(v)
}

func (e *Ears) CalculateVolume(v physics.Vector) float64 {
	v2 := physics.NewVector(*e.X, *e.Y)
	dist := v2.Distance(v)

	fmt.Println("Vector Distance:", dist, v, v2)

	// Ignore scaling variable
	lin := (e.SilenceRadius - dist) / e.SilenceRadius
	if lin < 0 {
		lin = 0
	}

	fmt.Println("Silence scale", lin, e.SilenceRadius, dist)

	return lin
}
