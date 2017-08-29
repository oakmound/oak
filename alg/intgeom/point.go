package intgeom

import "math"

// Point is a basic integer pair
type Point struct {
	X, Y int
}

// NewPoint returns an (X,Y) point structure
func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

// Distance is the euclidean distance function applied to the two points
// Consider: offer non-euclidean distance functions
func (p Point) Distance(p2 Point) float64 {
	return Distance(p.X, p.Y, p2.X, p2.Y)
}

// Distance calculates the euclidean distance between two points, as two (x,y) pairs
func Distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(
		math.Pow((float64(x1)-float64(x2)), 2) +
			math.Pow((float64(y1)-float64(y2)), 2))
}

// Add sums the location of a second point onto the Point
func (p Point) Add(p2 Point) Point {
	p.X += p2.X
	p.Y += p2.Y
	return p
}
