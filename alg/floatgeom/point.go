package floatgeom

import "math"

// Point2 represents a 2D point in space.
type Point2 [2]float64

// Point3 represents a 3D point in space.
type Point3 [3]float64

// AnglePoint creates a unit vector from the given angle in degrees as a Point2.
func AnglePoint(angle float64) Point2 {
	return RadianPoint(angle * math.Pi / 180)
}

// RadianPoint creates a unit vector from the given angle in radians as a Point2.
func RadianPoint(radians float64) Point2 {
	return Point2{math.Cos(radians), math.Sin(radians)}
}

// Dim returns the value of p in the ith dimension.
// Panics if i > 1. No check is made for efficiency's sake, pending benchmarks,
// but adding an error here would significantly worsen the API.
func (p Point2) Dim(i int) float64 {
	return p[i]
}

// Dim returns the value of p in the ith dimension.
// Panics if i > 2. No check is made for efficiency's sake, pending benchmarks,
// but adding an error here would significantly worsen the API.
func (p Point3) Dim(i int) float64 {
	return p[i]
}

// X returns p's value on the X axis.
func (p Point2) X() float64 {
	return p.Dim(0)
}

// Y returns p's value on the Y axis.
func (p Point2) Y() float64 {
	return p.Dim(1)
}

// X returns p's value on the X axis.
func (p Point3) X() float64 {
	return p.Dim(0)
}

// Y returns p's value on the Y axis.
func (p Point3) Y() float64 {
	return p.Dim(1)
}

// Z returns p's value on the Z axis.
func (p Point3) Z() float64 {
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
func Distance2(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(
		math.Pow(x1-x2, 2) +
			math.Pow(y1-y2, 2))
}

// Distance3 calculates the euclidean distance between two points, as two (x,y,z) triplets
func Distance3(x1, y1, z1, x2, y2, z2 float64) float64 {
	return math.Sqrt(
		math.Pow(x1-x2, 2) +
			math.Pow(y1-y2, 2) +
			math.Pow(z1-z2, 2))
}

// LesserOf returns the lowest values on each axis of the input points as a point.
func (p Point2) LesserOf(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] = math.Min(p[0], p2[0])
		p[1] = math.Min(p[1], p2[1])
	}
	return p
}

// LesserOf returns the lowest values on each axis of the input points as a point.
func (p Point3) LesserOf(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] = math.Min(p[0], p2[0])
		p[1] = math.Min(p[1], p2[1])
		p[2] = math.Min(p[2], p2[2])
	}
	return p
}

// GreaterOf returns the highest values on each axis of the input points as a point.
func (p Point2) GreaterOf(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] = math.Max(p[0], p2[0])
		p[1] = math.Max(p[1], p2[1])
	}
	return p
}

// GreaterOf returns the highest values on each axis of the input points as a point.
func (p Point3) GreaterOf(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] = math.Max(p[0], p2[0])
		p[1] = math.Max(p[1], p2[1])
		p[2] = math.Max(p[2], p2[2])
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
func (p Point2) MulConst(fs ...float64) Point2 {
	for _, f := range fs {
		p[0] *= f
		p[1] *= f
	}
	return p
}

// Div combines the input points via division.
func (p Point2) Div(ps ...Point2) Point2 {
	for _, p2 := range ps {
		p[0] /= p2[0]
		p[1] /= p2[1]
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
func (p Point3) MulConst(fs ...float64) Point3 {
	for _, f := range fs {
		p[0] *= f
		p[1] *= f
		p[2] *= f
	}
	return p
}

// Div combines the input points via division.
func (p Point3) Div(ps ...Point3) Point3 {
	for _, p2 := range ps {
		p[0] /= p2[0]
		p[1] /= p2[1]
		p[2] /= p2[2]
	}
	return p
}

// ToRect converts this point into a rectangle spanning span distance
// in each axis.
func (p Point2) ToRect(span float64) Rect2 {
	return NewRect2WH(p[0], p[1], span, span)
}

// ToRect converts this point into a rectangle spanning span distance
// in each axis.
func (p Point3) ToRect(span float64) Rect3 {
	return NewRect3WH(p[0], p[1], p[2], span, span, span)
}

// ToAngle returns this point as an angle in degrees.
func (p Point2) ToAngle() float64 {
	return p.ToRadians() * 180 / math.Pi
}

// ToRadians returns this point as an angle in radians.
func (p Point2) ToRadians() float64 {
	return math.Atan2(p[1], p[0])
}

// AngleTo returns the angle from p to p2 in degrees.
func (p Point2) AngleTo(p2 Point2) float64 {
	return p.Sub(p2).ToAngle()
}

// RadiansTo returns the angle from p to p2 in radians.
func (p Point2) RadiansTo(p2 Point2) float64 {
	return p.Sub(p2).ToRadians()
}
