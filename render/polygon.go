package render

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/span"
)

// A Polygon is a renderable that is represented by a set of in order points
// on a plane.
type Polygon struct {
	*Sprite
	floatgeom.Polygon2
}

// NewPointsPolygon is a helper function for `NewPolygon(floatgeom.NewPolygon2(p1, p2, p3, pn...))`
func NewPointsPolygon(p1, p2, p3 floatgeom.Point2, pn ...floatgeom.Point2) *Polygon {
	return NewPolygon(floatgeom.NewPolygon2(p1, p2, p3, pn...))
}

// NewPolygon constructs a renderable polygon. It will display nothing until
// Fill or FillInverse is called on it.
func NewPolygon(poly floatgeom.Polygon2) *Polygon {
	return &Polygon{
		Sprite: NewSprite(poly.Bounding.Min.X(), poly.Bounding.Min.Y(),
			image.NewRGBA(image.Rect(0, 0, int(poly.Bounding.W()), int(poly.Bounding.H())))),
		Polygon2: poly,
	}
}

// GetOutline returns a set of lines of the given color along this polygon's outline
func (pg *Polygon) GetOutline(c color.Color) *CompositeM {
	return pg.GetColoredOutline(IdentityColorer(c), 0)
}

// GetThickOutline returns a set of lines of the given color along this polygon's outline,
// at the given thickness
func (pg *Polygon) GetThickOutline(c color.Color, thickness int) *CompositeM {
	return pg.GetColoredOutline(IdentityColorer(c), thickness)
}

// GetGradientOutline returns a set of lines of the given color along this polygon's outline,
// at the given thickness, ranging from c1 to c2 in color
func (pg *Polygon) GetGradientOutline(c1, c2 color.Color, thickness int) *CompositeM {
	return pg.GetColoredOutline(span.NewLinearColor(c1, c2).Percentile, thickness)
}

// GetColoredOutline returns a set of lines of the given color along this polygon's outline
func (pg *Polygon) GetColoredOutline(colorer Colorer, thickness int) *CompositeM {
	sl := NewCompositeM()
	j := len(pg.Points) - 1
	for i, p2 := range pg.Points {
		p1 := pg.Points[j]
		MinX := math.Min(p1.X(), p2.X())
		MinY := math.Min(p1.Y(), p2.Y())
		sl.AppendOffset(
			NewColoredLine(p1.X(), p1.Y(), p2.X(), p2.Y(), colorer, thickness),
			floatgeom.Point2{MinX, MinY})
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

// Fill fills the inside of this polygon with the input color
func (pg *Polygon) Fill(c color.Color) {
	// Reset the rgba of the polygon
	bounds := pg.r.Bounds()
	rect := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	rgba := image.NewRGBA(rect)
	minx := pg.Bounding.Min.X()
	miny := pg.Bounding.Min.Y()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			if pg.Contains(float64(x)+minx, float64(y)+miny) {
				rgba.Set(x, y, c)
			}
		}
	}
	pg.r = rgba
}
