package physics

import (
	"math"
)

// A Vector is a two-dimensional point or vector used throughout oak
// to maintain functionality between packages.
type Vector struct {
	X, Y float64
}

var (
	// CUTOFF is used for rounding after floating point operations to
	// zero out vector values that are sufficiently close to zero
	CUTOFF = 0.001
)

// NewVector returns a vector with the given x and y components
func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

// Copy copies a Vector
func (v Vector) Copy() Vector {
	return NewVector(v.X, v.Y)
}

// Magnitude returns the magnitude of the combined components of a Vector
func (v Vector) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

// Normalize divides both components in a vector by the vector's magnitude
func (v Vector) Normalize() Vector {
	v = v.round()
	mgn := v.Magnitude()
	if mgn == 0 {
		return v
	}
	v.X /= mgn
	v.Y /= mgn
	return v
}

// Zero is shorthand for NewVector(0,0)
func (v Vector) Zero() Vector {
	v.X = 0
	v.Y = 0
	return v
}

// Add combines a set of vectors through addition
func (v Vector) Add(vs ...Vector) Vector {
	for _, v2 := range vs {
		v.X += v2.X
		v.Y += v2.Y
	}
	return v.round()
}

// Scale scales a vector by a set of floating points
// Scale(f1,f2,f3) is equivalent to Scale(f1*f2*f3)
func (v Vector) Scale(fs ...float64) Vector {
	f2 := 1.0
	for _, f := range fs {
		f2 *= f
	}
	v.X *= f2
	v.Y *= f2
	return v.round()
}

// Rotate takes in a set of angles and rotates v by their sum
// the input angles are assumed to be in degrees.
func (v Vector) Rotate(fs ...float64) Vector {
	angle := 0.0
	for _, f := range fs {
		angle += f
	}
	mgn := v.Magnitude()
	angle = math.Atan2(v.Y, v.X) + (angle * (math.Pi) / 180)
	v.X = math.Cos(angle) * mgn
	v.Y = math.Sin(angle) * mgn

	return v.round()
}

// Angle returns this vector as an angle in degrees
func (v Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

// Dot returns the dot product of the vectors
func (v Vector) Dot(v2 Vector) float64 {
	x := v.X * v2.X
	y := v.Y * v2.Y
	return x + y
}

// Distance on two vectors returns the euclidean distance
// from v to v2
func (v Vector) Distance(v2 Vector) float64 {
	return v.Add(v2.Scale(-1)).Magnitude()
}

func (v Vector) round() Vector {
	if math.Abs(v.X) < CUTOFF {
		v.X = 0
	}
	if math.Abs(v.Y) < CUTOFF {
		v.Y = 0
	}
	return v
}

// ShiftX is equivalent to v.X += x
func (v Vector) ShiftX(x float64) Vector {
	v.X += x
	return v
}

// ShiftY is equivalent to v.Y += y
func (v Vector) ShiftY(y float64) Vector {
	v.Y += y
	return v
}

// GetX returns v.X
func (v Vector) GetX() float64 {
	return v.X
}

// GetY returns v.Y
func (v Vector) GetY() float64 {
	return v.Y
}

// SetPos is equivalent to NewVector(x,y)
func (v Vector) SetPos(x, y float64) Vector {
	v.X = x
	v.Y = y
	return v
}

// GetPos returns both v.X and v.Y
func (v Vector) GetPos() (float64, float64) {
	return v.X, v.Y
}
