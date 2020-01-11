package intgeom

import (
	"math"

	"github.com/oakmound/oak/v2/alg"
)

// Point2 represents a 2D point in space.
type Point2 [2]int

// Point3 represents a 3D point in space.
type Point3 [3]int

// Dim returns the value of p in the ith dimension.
// Panics if i > 1. No check is made for efficiency's sake, pending benchmarks,
// but adding an error here would significantly worsen the API.
func (p Point2) Dim(i int) int {
	return p[i]
}

// Dim returns the value of p in the ith dimension.
// Panics if i > 2. No check is made for efficiency's sake, pending benchmarks,
// but adding an error here would significantly worsen the API.
func (p Point3) Dim(i int) int {
	return p[i]
}

// X returns p's value on the X axis.
func (p Point2) X() int {
	return p.Dim(0)
}

// Y returns p's value on the Y axis.
func (p Point2) Y() int {
	return p.Dim(1)
}

// X returns p's value on the X axis.
func (p Point3) X() int {
	return p.Dim(0)
}

// Y returns p's value on the Y axis.
func (p Point3) Y() int {
	return p.Dim(1)
}

// Z returns p's value on the Z axis.
func (p Point3) Z() int {
	return p.Dim(2)
}

// Distance calculates the distance between this Point2 and another.
func (p Point2) Distance(p2 Point2) float64 {
	return Distance2(p.X(), p.Y(), p2.X(), p2.Y())
}

// Distance calculates the distance between this Point3 and another.
func (p Point3) Distance(p2 Point3) float64 {
	return Distance3(p.X(), p.Y(), p.Z(), p2.X(), p2.Y(), p2.Z())
}

// Distance2 calculates the euclidean distance between two points, as two (x,y) pairs
func Distance2(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(
		math.Pow(float64(x1-x2), 2) +
			math.Pow(float64(y1-y2), 2))
}

// Distance3 calculates the euclidean distance between two points, as two (x,y,z) triplets
func Distance3(x1, y1, z1, x2, y2, z2 int) float64 {
	return math.Sqrt(
		math.Pow(float64(x1-x2), 2) +
			math.Pow(float64(y1-y2), 2) +
			math.Pow(float64(z1-z2), 2))
}

// LesserOf returns the lowest values on each axis of the input points as a point.
func (p Point2) LesserOf(ps ...Point2) Point2 {
	for _, p2 := range ps {
		if p2[0] < p[0] {
			p[0] = p2[0]
		}
		if p2[1] < p[1] {
			p[1] = p2[1]
		}
	}
	return p
}

// LesserOf returns the lowest values on each axis of the input points as a point.
func (p Point3) LesserOf(ps ...Point3) Point3 {
	for _, p2 := range ps {
		if p2[0] < p[0] {
			p[0] = p2[0]
		}
		if p2[1] < p[1] {
			p[1] = p2[1]
		}
		if p2[2] < p[2] {
			p[2] = p2[2]
		}
	}
	return p
}

// GreaterOf returns the highest values on each axis of the input points as a point.
func (p Point2) GreaterOf(ps ...Point2) Point2 {
	for _, p2 := range ps {
		if p2[0] > p[0] {
			p[0] = p2[0]
		}
		if p2[1] > p[1] {
			p[1] = p2[1]
		}
	}
	return p
}

// GreaterOf returns the highest values on each axis of the input points as a point.
func (p Point3) GreaterOf(ps ...Point3) Point3 {
	for _, p2 := range ps {
		if p2[0] > p[0] {
			p[0] = p2[0]
		}
		if p2[1] > p[1] {
			p[1] = p2[1]
		}
		if p2[2] > p[2] {
			p[2] = p2[2]
		}
	}
	return p
}

// Add combines the input points via addition.
func (p Point2) Add(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] += p2[0]
		p[1] += p2[1]
	}
	return p
}

// Sub combines the input points via subtraction.
func (p Point2) Sub(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] -= p2[0]
		p[1] -= p2[1]
	}
	return p
}

// Mul combines in the input points via multiplication.
func (p Point2) Mul(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] *= p2[0]
		p[1] *= p2[1]
	}
	return p
}

// MulConst multiplies all elements of a point by the input floats
func (p Point2) MulConst(fs ...int) Point2 {
	for _, f := range fs {
		p[0] *= f
		p[1] *= f
	}
	return p
}

// Div combines the input points via division.
// Div does not check that the inputs are non zero before operating,
// and can panic if that is not true.
func (p Point2) Div(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] /= p2[0]
		p[1] /= p2[1]
	}
	return p
}

// DivConst divides all elements of a point by the input floats
// DivConst does not check that the inputs are non zero before operating,
// and can panic if that is not true.
func (p Point2) DivConst(fs ...int) Point2 {
	for _, f := range fs {
		p[0] /= f
		p[1] /= f
	}
	return p
}

// Add combines the input points via addition.
func (p Point3) Add(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] += p2[0]
		p[1] += p2[1]
		p[2] += p2[2]
	}
	return p
}

// Sub combines the input points via subtraction.
func (p Point3) Sub(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] -= p2[0]
		p[1] -= p2[1]
		p[2] -= p2[2]
	}
	return p
}

// Mul combines in the input points via multiplication.
func (p Point3) Mul(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] *= p2[0]
		p[1] *= p2[1]
		p[2] *= p2[2]
	}
	return p
}

// MulConst multiplies all elements of a point by the input floats
func (p Point3) MulConst(fs ...int) Point3 {
	for _, f := range fs {
		p[0] *= f
		p[1] *= f
		p[2] *= f
	}
	return p
}

// Div combines the input points via division.
// Div does not check that the inputs are non zero before operating,
// and can panic if that is not true.
func (p Point3) Div(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] /= p2[0]
		p[1] /= p2[1]
		p[2] /= p2[2]
	}
	return p
}

// DivConst divides all elements of a point by the input floats
// DivConst does not check that the inputs are non zero before operating,
// and can panic if that is not true.
func (p Point3) DivConst(fs ...int) Point3 {
	for _, f := range fs {
		p[0] /= f
		p[1] /= f
		p[2] /= f
	}
	return p
}

// Dot returns the dot product of the input points
func (p Point2) Dot(p2 Point2) int {
	return p[0]*p2[0] + p[1]*p2[1]
}

// Dot returns the dot product of the input points
func (p Point3) Dot(p2 Point3) int {
	return p[0]*p2[0] + p[1]*p2[1] + p[2]*p2[2]
}

// Magnitude returns the magnitude of the combined components of a Point
func (p Point2) Magnitude() float64 {
	return math.Sqrt(float64(p.Dot(p)))
}

// Magnitude returns the magnitude of the combined components of a Point
func (p Point3) Magnitude() float64 {
	return math.Sqrt(float64(p.Dot(p)))
}

// ToRect converts this point into a rectangle spanning span distance
// in each axis.
func (p Point2) ToRect(span int) Rect2 {
	return NewRect2WH(p[0], p[1], span, span)
}

// ToRect converts this point into a rectangle spanning span distance
// in each axis.
func (p Point3) ToRect(span int) Rect3 {
	return NewRect3WH(p[0], p[1], p[2], span, span, span)
}

// ProjectX projects the Point3 onto the x axis, removing it's
// x component and returning a Point2
// todo: I'm not sure about this (these) function name
func (p Point3) ProjectX() Point2 {
	return Point2{p[1], p[2]}
}

// ProjectY projects the Point3 onto the y axis, removing it's
// y component and returning a Point2
func (p Point3) ProjectY() Point2 {
	return Point2{p[0], p[2]}
}

// ProjectZ projects the Point3 onto the z axis, removing it's
// z component and returning a Point2
func (p Point3) ProjectZ() Point2 {
	return Point2{p[0], p[1]}
}

// ToAngle returns this point as an angle in degrees.
func (p Point2) ToAngle() float64 {
	return p.ToRadians() * alg.RadToDeg
}

// ToRadians returns this point as an angle in radians.
func (p Point2) ToRadians() float64 {
	return math.Atan2(float64(p[1]), float64(p[0]))
}

// AngleTo returns the angle from p to p2 in degrees.
func (p Point2) AngleTo(p2 Point2) float64 {
	return p.Sub(p2).ToAngle()
}

// RadiansTo returns the angle from p to p2 in radians.
func (p Point2) RadiansTo(p2 Point2) float64 {
	return p.Sub(p2).ToRadians()
}
