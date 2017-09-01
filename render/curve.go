package render

import (
	"errors"
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

// BezierCurve will form a Bezier on the given coordinates, expected in (x,y)
// pairs. If the inputs have an odd length, an error noting so is returned, and
// the Bezier returned is nil.
func BezierCurve(coords ...float64) (Bezier, error) {
	if len(coords)%2 != 0 {
		return nil, errors.New("Invalid number of inputs, len must be divisible by 2")
	}
	pts := make([]Bezier, len(coords)/2)
	for i := 0; i < len(coords); i += 2 {
		pts[i/2] = bezierPoint{coords[i], coords[i+1]}
	}
	for len(pts) > 1 {
		for i := 0; i < len(pts)-1; i++ {
			pts[i] = bezierNode{pts[i], pts[i+1]}
		}
		pts = pts[:len(pts)-1]
	}
	return pts[0], nil
}

// A Bezier has a function indicating how far along a curve something is given
// some float64 progress between 0 and 1. This allows points, lines, and limitlessly complex
// bezier curves to be represented under this interface.
//
// Beziers will not necessarily break if given an input outside of 0-1, but the results
// shouldn't be relied upon.
type Bezier interface {
	Pos(progress float64) (x, y float64)
}

type bezierNode struct {
	left, right Bezier
}

func (bn bezierNode) Pos(progress float64) (x, y float64) {
	x1, y1 := bn.left.Pos(progress)
	x2, y2 := bn.right.Pos(progress)
	return x1 + ((x2 - x1) * progress), y1 + ((y2 - y1) * progress)
}

// bezierPoints cover cases where only 1 point is supplied, and serve as roots.
type bezierPoint struct {
	x, y float64
}

func (bp bezierPoint) Pos(progress float64) (x, y float64) {
	return bp.x, bp.y
}
