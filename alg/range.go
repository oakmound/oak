package alg

import (
	"math"
	"math/rand"
)

// IntRange represents the ability
// to poll a struct and return an integer,
// distributed over some range dependant
// on the implementing struct.
type IntRange interface {
	Poll() int
}

// LinearIntRange polls on a linear scale
// between a minimum and a maximum
type LinearIntRange struct {
	Min, Max int
}

// Poll returns an integer distributed
// between lir.Min and lir.Max
func (lir LinearIntRange) Poll() int {
	return rand.Intn(lir.Max-lir.Min) + lir.Min
}

type BaseSpreadRangei struct {
	Base, Spread int
}

func (b BaseSpreadRangei) Poll() int {
	return b.Base + rand.Intn((b.Spread*2)+1) - b.Spread
}

// Constant implements IntRange as a poll
// which always returns the same integer.
type Constant int

// Poll returns c cast to an int
func (c Constant) Poll() int {
	return int(c)
}

type Infinite struct{}

func (inf Infinite) Poll() int {
	return math.MaxInt32
}

type FloatRange interface {
	Poll() float64
	Mult(f float64) FloatRange
}

type BaseSpreadRange struct {
	Base, Spread float64
}

func (b BaseSpreadRange) Poll() float64 {
	return b.Base + (b.Spread * 2 * rand.Float64()) - b.Spread
}

func (b BaseSpreadRange) Mult(f float64) FloatRange {
	b.Base = b.Base * f
	b.Spread = b.Spread * f
	return b
}

type Constantf float64

func (c Constantf) Poll() float64 {
	return float64(c)
}

func (c Constantf) Mult(f float64) FloatRange {
	c = Constantf(float64(c) * f)
	return c
}
