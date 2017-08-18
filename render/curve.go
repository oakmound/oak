package render

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/physics"
)

// DrawCircle draws a circle on the input rgba, of color c.
func DrawCircle(rgba *image.RGBA, c color.Color, radius, thickness float64, offsets ...float64) {
	DrawCurve(rgba, c, radius, thickness, 0, 1, offsets...)
}

// DrawCurve draws a curve inward on the input rgba, of color c.
func DrawCurve(rgba *image.RGBA, c color.Color, radius, thickness, initialAngle, circlePercentage float64, offsets ...float64) {
	offX := 0.0
	offY := 0.0
	if len(offsets) > 0 {
		offX = offsets[0]
		if len(offsets) > 1 {
			offY = offsets[1]
		}
	}
	rVec := physics.NewVector(radius+offX, radius+offY)
	delta := physics.AngleVector(initialAngle)
	circum := 2 * radius * math.Pi
	rotation := 180 / circum
	for j := 0.0; j < circum*2*circlePercentage; j++ {
		delta.Rotate(rotation)
		// We add rVec to move from -1->1 to 0->2 in terms of radius scale
		start := delta.Copy().Scale(radius).Add(rVec)
		for i := 0.0; i <= thickness; i++ {
			// this pixel is radius minus the delta, to move inward
			cur := start.Add(delta.Copy().Scale(-1))
			rgba.Set(alg.RoundF64(cur.X()), alg.RoundF64(cur.Y()), c)
		}
	}
}
