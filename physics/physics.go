package physics

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{x, y}
}

func (v *Vector) Copy() *Vector {
	return NewVector(v.X, v.Y)
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v *Vector) Normalize() *Vector {
	mgn := v.Magnitude()
	if mgn == 0 {
		return v
	}
	v.X /= mgn
	v.Y /= mgn
	return v
}

func (v *Vector) Zero() *Vector {
	v.X = 0
	v.Y = 0
	return v
}

func (v *Vector) Add(vs ...*Vector) *Vector {
	for _, v2 := range vs {
		v.X += v2.X
		v.Y += v2.Y
	}
	return v
}

func (v *Vector) Scale(fs ...float64) *Vector {
	f2 := 1.0
	for _, f := range fs {
		f2 *= f
	}
	v.X *= f2
	v.Y *= f2
	return v
}

func (v *Vector) Dot(v2 *Vector) float64 {
	x := v.X * v2.X
	y := v.Y * v2.Y
	return x + y
}
