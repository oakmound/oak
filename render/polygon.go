package render

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/oakerr"
)

// A Polygon is a renderable that is represented by a set of in order points
// on a plane.
type Polygon struct {
	*Sprite
	Rect2  floatgeom.Rect2
	points []floatgeom.Point2
}

// NewStrictPolygon will draw a polygon of points within a given rectangle,
// and if the input points lie outside of that rectangle the polygon will clip
// into and not be drawn outside of that border.
func NewStrictPolygon(bounds floatgeom.Rect2, points ...floatgeom.Point2) (*Polygon, error) {
	if len(points) < 3 {
		return nil, oakerr.InsufficientInputs{AtLeast: 3, InputName: "points"}
	}
	return &Polygon{
		Sprite: NewSprite(bounds.Min.X(), bounds.Min.Y(),
			image.NewRGBA(image.Rect(0, 0, int(bounds.W()), int(bounds.H())))),
		Rect2:  bounds,
		points: points,
	}, nil
}

// NewPolygon takes in a set of points and returns a polygon. At least three points
// must be provided.
func NewPolygon(points ...floatgeom.Point2) (*Polygon, error) {

	if len(points) < 3 {
		return nil, oakerr.InsufficientInputs{AtLeast: 3, InputName: "points"}
	}

	// Calculate the bounding rectangle of the polygon by
	// finding the maximum and minimum x and y values of the given points
	return NewStrictPolygon(floatgeom.NewBoundingRect2(points...), points...)
}

// UpdatePoints resets the points of this polygon to be the passed in points
// Todo 2.0: Take in a variadic set instead of a slice
func (pg *Polygon) UpdatePoints(points ...floatgeom.Point2) error {
	if len(points) < 3 {
		return oakerr.InsufficientInputs{AtLeast: 3, InputName: "points"}
	}
	pg.points = points
	pg.Rect2 = floatgeom.NewBoundingRect2(points...)
	return nil
}

// Fill fills the inside of this polygon with the input color
func (pg *Polygon) Fill(c color.Color) {
	// Reset the rgba of the polygon
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	minx := pg.Rect2.Min.X()
	miny := pg.Rect2.Min.Y()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if pg.Contains(float64(x)+minx, float64(y)+miny) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

// GetOutline returns a set of lines of the given color along this polygon's outline
func (pg *Polygon) GetOutline(c color.Color) *Composite {
	sl := NewComposite()
	j := len(pg.points) - 1
	for i, p2 := range pg.points {
		p1 := pg.points[j]
		MinX := math.Min(p1.X(), p2.X())
		MinY := math.Min(p1.Y(), p2.Y())
		sl.AppendOffset(NewLine(p1.X(), p1.Y(), p2.X(), p2.Y(), c), floatgeom.Point2{MinX, MinY})
		j = i
	}
	return sl
}

// FillInverse colors this polygon's exterior the given color
func (pg *Polygon) FillInverse(c color.Color) {
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if !pg.ConvexContains(float64(x), float64(y)) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}

// Todo: almost all of this junk below should be in alg, under floatgeom or something.

// Contains returns whether or not the current Polygon contains the passed in Point.
// It is the default containment function, versus wrapping and convex.
func (pg *Polygon) Contains(x, y float64) (contains bool) {

	if !pg.Rect2.Contains(floatgeom.Point2{x, y}) {
		return
	}

	j := len(pg.points) - 1
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[j]
		if (tp1.Y() > y) != (tp2.Y() > y) { // Three comparisons
			if x < (tp2.X()-tp1.X())*(y-tp1.Y())/(tp2.Y()-tp1.Y())+tp1.X() { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		j = i
	}
	return
}

// ConvexContains returns whether the given point is contained by the input polygon.
// It assumes the polygon is convex. It outperforms the alternatives.
func (pg *Polygon) ConvexContains(x, y float64) bool {

	p := floatgeom.Point2{x, y}

	if !pg.Rect2.Contains(p) {
		return false
	}

	prev := 0
	for i := 0; i < len(pg.points); i++ {
		tp1 := pg.points[i]
		tp2 := pg.points[(i+1)%len(pg.points)]
		tp3 := tp2.Sub(tp1)
		tp4 := p.Sub(tp1)
		cur := getSide(tp3, tp4)
		if cur == 0 {
			return false
		} else if prev == 0 {
		} else if prev != cur {
			return false
		}
		prev = cur
	}
	return true
}

func getSide(a, b floatgeom.Point2) int {
	x := a.X()*b.Y() - a.Y()*b.X()
	if x == 0 {
		return 0
	} else if x < 1 {
		return -1
	} else {
		return 1
	}
}
