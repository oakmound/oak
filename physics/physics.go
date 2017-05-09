package physics

import (
	"math"
)

type Vector struct {
	X, Y float64
}

var (
	CUTOFF = 0.001
)

func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

func (v Vector) Copy() Vector {
	return NewVector(v.X, v.Y)
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

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

func (v Vector) Zero() Vector {
	v.X = 0
	v.Y = 0
	return v
}

func (v Vector) Add(vs ...Vector) Vector {
	for _, v2 := range vs {
		v.X += v2.X
		v.Y += v2.Y
	}
	return v.round()
}

func (v Vector) Scale(fs ...float64) Vector {
	f2 := 1.0
	for _, f := range fs {
		f2 *= f
	}
	v.X *= f2
	v.Y *= f2
	return v.round()
}

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

// Defaults to Degrees
func (v Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

func (v Vector) Dot(v2 Vector) float64 {
	x := v.X * v2.X
	y := v.Y * v2.Y
	return x + y
}

// Distance is Euclidean
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

func (v Vector) ShiftX(x float64) Vector {
	v.X += x
	return v
}
func (v Vector) ShiftY(y float64) Vector {
	v.Y += y
	return v
}
func (v Vector) GetX() float64 {
	return v.X
}
func (v Vector) GetY() float64 {
	return v.Y
}

func (v Vector) SetPos(x, y float64) Vector {
	v.X = x
	v.Y = y
	return v
}

func (v Vector) GetPos() (float64, float64) {
	return v.X, v.Y
}
