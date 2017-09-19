package render

import (
	"fmt"
	"image"
	"image/color"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/shape"
)

// BezierLine converts a bezier into a line sprite.
func BezierLine(b shape.Bezier, c color.Color) *Sprite {
	low := 0.0
	high := 1.0
	pts := make([]intgeom.Point, 2)
	pts[0] = roundToIntPoint(b.Pos(low))
	pts[1] = roundToIntPoint(b.Pos(high))
	bezierDraw(b, &pts, low, high)

	// Obtain min and max x y values
	min := pts[0]
	max := pts[0]
	for _, p := range pts {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, max.X-min.X, max.Y-min.Y))

	for _, p := range pts {
		rgba.Set(p.X-min.X, p.Y-min.Y, c)
	}

	return NewSprite(float64(min.X), float64(min.Y), rgba)

}

func roundToIntPoint(x, y float64) intgeom.Point {
	return intgeom.NewPoint(alg.RoundF64(x), alg.RoundF64(y))
}

func bezierDraw(b shape.Bezier, pts *[]intgeom.Point, low, high float64) {
	p1 := roundToIntPoint(b.Pos(low))
	p2 := roundToIntPoint(b.Pos(high))
	mid := (low + high) / 2
	p3 := roundToIntPoint(b.Pos(mid))
	// If we've already added this point at this low or high value
	// I.E. if the low and high values are sufficiently close together
	// (high-low < 0.01) &&
	if (p3.X == p2.X && p3.Y == p2.Y) ||
		(p3.X == p1.X && p3.Y == p1.Y) {
		fmt.Println(p3, p2, p1)
		return
	}
	*pts = append(*pts, p3)
	bezierDraw(b, pts, low, mid)
	bezierDraw(b, pts, mid, high)
}
