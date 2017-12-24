package render

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/shape"
)

// BezierLine converts a bezier into a line sprite.
func BezierLine(b shape.Bezier, c color.Color) *Sprite {
	return BezierThickLine(b, c, 0)
}

// BezierThickLine draws a BezierLine wrapping each colored pixel in
// a square of width and height = thickness
func BezierThickLine(b shape.Bezier, c color.Color, thickness int) *Sprite {
	low := 0.0
	high := 1.0
	pts := make([]intgeom.Point, 2)
	pts[0] = roundToIntPoint(b.Pos(low))
	pts[1] = roundToIntPoint(b.Pos(high))
	bezierDraw(b, &pts, low, high, pts[0], pts[1])

	min := pts[0].LesserOf(pts...)
	max := pts[0].GreaterOf(pts...)

	rgba := image.NewRGBA(image.Rect(0, 0, 1+(max.X-min.X), 1+(max.Y-min.Y)))

	for _, p := range pts {
		x := p.X - min.X
		y := p.Y - min.Y
		for i := x - thickness; i <= x+thickness; i++ {
			for j := y - thickness; j <= y+thickness; j++ {
				rgba.Set(i, j, c)
			}
		}
	}

	return NewSprite(float64(min.X), float64(min.Y), rgba)
}

func roundToIntPoint(x, y float64) intgeom.Point {
	return intgeom.NewPoint(alg.RoundF64(x), alg.RoundF64(y))
}

func bezierDraw(b shape.Bezier, pts *[]intgeom.Point, low, high float64, lowPt, highPt intgeom.Point) {
	mid := (low + high) / 2
	p := roundToIntPoint(b.Pos(mid))
	// If we haven't yet added this point at this low or high value
	// I.E. if the low and high values are sufficiently close together
	if p != lowPt && p != highPt {
		*pts = append(*pts, p)

		bezierDraw(b, pts, low, mid, lowPt, p)
		bezierDraw(b, pts, mid, high, p, highPt)
	}
}
