package collision

import "bitbucket.org/oakmoundstudio/oak/physics"

// A Point is a specific point where
// collision occured and a zone to identify
// what was collided with.
type Point struct {
	physics.Vector
	Zone *Space
}

// NilPoint returns a Point representing no collision
func NilPoint() Point {
	return Point{physics.NewVector(0, 0), nil}
}

// NewPoint creates a new point
func NewPoint(s *Space, x, y float64) Point {
	return Point{physics.NewVector(x, y), s}
}

// IsNil returns whether the underlying zone of a Point is nil
func (cp Point) IsNil() bool {
	return cp.Zone == nil
}
