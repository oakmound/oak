package alg

import "math"

const (
	// DegToRad is the constant value something in
	// degrees should be multiplied by to obtain
	// something in radians.
	DegToRad = math.Pi / 180
	// RadToDeg is the constant value something in
	// radians should be multiplied by to obtain
	// something in degrees.
	RadToDeg = 180 / math.Pi
)

// We might not want these types
// It might be too much of a hassle to deal with them

// A Radian value is a float that specifies it should be in radians.
type Radian float64

// Degrees converts a Radian to Degrees.
func (r Radian) Degrees() Degree {
	return Degree(r * RadToDeg)
}

// A Degree value is a float that specifies it should be in degrees.
type Degree float64

// Radians converts a Degree to Radians.
func (d Degree) Radians() Radian {
	return Radian(d * DegToRad)
}
