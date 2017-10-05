package collision

import "github.com/oakmound/oak/alg/floatgeom"

// A Point is a specific point where
// collision occurred and a zone to identify
// what was collided with.
type Point struct {
	floatgeom.Point3
	Zone *Space
}

// NewPoint creates a new point
func NewPoint(s *Space, x, y float64) Point {
	return Point{floatgeom.Point3{x, y, 0}, s}
}

// IsNil returns whether the underlying zone of a Point is nil
func (cp Point) IsNil() bool {
	return cp.Zone == nil
}
