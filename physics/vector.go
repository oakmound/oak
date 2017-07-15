package physics

import (
	"fmt"
	"math"
	"runtime"
)

// A Vector is a two-dimensional point or vector used throughout oak
// to maintain functionality between packages.
type Vector struct {
	x, y       *float64
	offX, offY float64
}

const (
	// CUTOFF is used for rounding after floating point operations to
	// zero out vector values that are sufficiently close to zero
	CUTOFF = 0.001
)

// NewVector returns a vector with the given x and y components
func NewVector(x, y float64) Vector {
	x2 := new(float64)
	y2 := new(float64)
	*x2 = x
	*y2 = y
	return Vector{x2, y2, 0, 0}
}

// MaxVector returns whichever vector has a greater magnitude
func MaxVector(a, b Vector) Vector {
	if a.Magnitude() > b.Magnitude() {
		return a
	}
	return b
}

// Copy copies a Vector
func (v Vector) Copy() Vector {
	if v.x == nil || v.y == nil {
		_, f, line, _ := runtime.Caller(2)
		fmt.Println("This vector was bad ", v, f, line)
		return v.Zero()
	}
	return NewVector(*v.x, *v.y)
}

// Magnitude returns the magnitude of the combined components of a Vector
func (v Vector) Magnitude() float64 {
	return math.Sqrt((*v.x * *v.x) + (*v.y * *v.y))
}

// Normalize divides both components in a vector by the vector's magnitude
func (v Vector) Normalize() Vector {
	v = v.round()
	mgn := v.Magnitude()
	if mgn == 0 {
		return v
	}
	*v.x /= mgn
	*v.y /= mgn
	return v
}

// Zero is shorthand for NewVector(0,0)
func (v Vector) Zero() Vector {
	return v.SetPos(0, 0)
}

// Add combines a set of vectors through addition
func (v Vector) Add(vs ...Vector) Vector {
	for _, v2 := range vs {
		*v.x += *v2.x
		*v.y += *v2.y
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
	*v.x *= f2
	*v.y *= f2
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
	angle = math.Atan2(*v.y, *v.x) + (angle * (math.Pi) / 180)

	return v.SetPos(math.Cos(angle)*mgn, math.Sin(angle)*mgn)
}

// Angle returns this vector as an angle in degrees
func (v Vector) Angle() float64 {
	return math.Atan2(*v.y, *v.x) * 180 / math.Pi
}

// Dot returns the dot product of the vectors
func (v Vector) Dot(v2 Vector) float64 {
	x := *v.x * *v2.x
	y := *v.y * *v2.y
	return x + y
}

// Distance on two vectors returns the euclidean distance
// from v to v2
func (v Vector) Distance(v2 Vector) float64 {
	v3 := v.Copy()
	v4 := v2.Copy()
	return v3.Add(v4.Scale(-1)).Magnitude()
}

func (v Vector) round() Vector {
	if math.Abs(*v.x) < CUTOFF {
		*v.x = 0
	}
	if math.Abs(*v.y) < CUTOFF {
		*v.y = 0
	}
	return v
}

// ShiftX is equivalent to v.X() += x
func (v Vector) ShiftX(x float64) Vector {
	*v.x += x
	return v
}

// ShiftY is equivalent to v.Y() += y
func (v Vector) ShiftY(y float64) Vector {
	*v.y += y
	return v
}

// X returns this vector's x component
func (v Vector) X() float64 {
	return *v.x + v.offX
}

// GetX returns this vector's x component todo: consider not having this
func (v Vector) GetX() float64 {
	return v.X()
}

// Y returns this vector's x component
func (v Vector) Y() float64 {
	return *v.y + v.offY
}

// GetY returns this vector's x component todo: consider not having this
func (v Vector) GetY() float64 {
	return v.Y()
}

// SetX returns a vector with its x component set to x
func (v Vector) SetX(x float64) Vector {
	*v.x = x
	return v.round()
}

// SetY returns a vector with its y component set to y
func (v Vector) SetY(y float64) Vector {
	*v.y = y
	return v.round()
}

// Xp returns the real pointer behind this vector's x component
func (v Vector) Xp() *float64 {
	return v.x
}

// Yp returns the real pointer behind  this vector's y component
func (v Vector) Yp() *float64 {
	return v.y
}

// SetPos is equivalent to NewVector(x,y)
func (v Vector) SetPos(x, y float64) Vector {
	*v.x = x
	*v.y = y
	return v
}

// GetPos returns both v.X() and v.Y()
func (v Vector) GetPos() (float64, float64) {
	return *v.x, *v.y
}

// AngleVector creates a unit vector by the cosine and sine of the given angle
func AngleVector(angle float64) Vector {
	angle *= math.Pi / 180
	return NewVector(math.Cos(angle), math.Sin(angle))
}
