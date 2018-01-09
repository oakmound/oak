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

// Distance is the euclidean distance function
func (p Point) Distance(p2 Point) float64 {
	return Distance(p.X, p.Y, p2.X, p2.Y)
}

// Distance is the euclidean distance function
// from two implicit int pairs
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

// LesserOf returns a point of the lowest X and Y component in the inputs.
func (p Point) LesserOf(ps ...Point) Point {
	for _, p2 := range ps {
		if p2.X < p.X {
			p.X = p2.X
		}
		if p2.Y < p.Y {
			p.Y = p2.Y
		}
	}
	return p
}

// GreaterOf returns a point of the highest X and Y component in the inputs.
func (p Point) GreaterOf(ps ...Point) Point {
	for _, p2 := range ps {
		if p2.X > p.X {
			p.X = p2.X
		}
		if p2.Y > p.Y {
			p.Y = p2.Y
		}
	}
	return p
}

// PointsBetween returns a line of points connecting p and p2
func (p Point) PointsBetween(p2 Point) []Point {

	out := make([]Point, 0)

	x1 := p.X
	y1 := p.Y
	x2 := p2.X
	y2 := p2.Y

	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	xSlope := -1
	if x2 < x1 {
		xSlope = 1
	}
	ySlope := -1
	if y2 < y1 {
		ySlope = 1
	}

	err := xDelta - yDelta
	var err2 float64
	for i := 0; true; i++ {
		out = append(out, Point{x2, y2})
		if x2 == x1 && y2 == y1 {
			break
		}
		err2 = 2 * err
		if err2 > -1*yDelta {
			err -= yDelta
			x2 += xSlope
		}
		if err2 < xDelta {
			err += xDelta
			y2 += ySlope
		}
	}
	return out
}
