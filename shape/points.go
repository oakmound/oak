package shape

import (
	"github.com/oakmound/oak/v2/alg/intgeom"
)

// Points is a shape defined by a set of points.
// It ignores input width and height given to it as it only cares about its points.
type Points map[intgeom.Point2]struct{}

// NewPoints creates a Points shape from any number of intgeom Points
func NewPoints(ps ...intgeom.Point2) Shape {
	points := make(map[intgeom.Point2]struct{}, len(ps))
	for _, p := range ps {
		points[p] = struct{}{}
	}
	return Points(points)
}

// In returns whether the input x and y are a point in the point map
func (p Points) In(x, y int, sizes ...int) bool {
	_, ok := p[intgeom.Point2{x, y}]
	return ok
}

// Outline returns the set of points along the point map's outline, if
// one exists
func (p Points) Outline(sizes ...int) ([]intgeom.Point2, error) {
	return ToOutline(p)(sizes...)
}

// Rect returns a double slice of booleans representing the output of the In function in that rectangle
func (p Points) Rect(sizes ...int) [][]bool {
	return InToRect(p.In)(sizes...)
}
