package alg

import (
	"errors"
	"math"
	"math/rand"
)

// IntRange represents the ability
// to poll a struct and return an integer,
// distributed over some range dependant
// on the implementing struct.
type IntRange interface {
	Poll() int
	Mult(int) IntRange
}

func NewLinearIntRange(min, max int) (IntRange, error) {
	if max <= min {
		return Constant(min), errors.New("Max cannot exceed or equal Min, returning constant(Min)")
	}
	return linearIntRange{min, max}, nil
}

func NewSpreadIntRange(base, spread int) (IntRange, error) {
	if spread < 0 {
		return Constant(base), errors.New("Spread cannot be < 0. Returning constant(Base)")
	}
	return spreadIntRange{base, spread}, nil
}

// NewConstant and NewInfinite are not necessary.
// They only exist to match the NewLIR and NewSIR functions.
func NewConstant(i int) IntRange {
	return Constant(i)
}

func NewInfinite() IntRange {
	return Infinite{}
}

// LinearIntRange polls on a linear scale
// between a minimum and a maximum
type linearIntRange struct {
	Min, Max int
}

func (lir linearIntRange) Poll() int {
	return rand.Intn(lir.Max-lir.Min) + lir.Min
}

func (lir linearIntRange) Mult(i int) IntRange {
	lir.Max *= i
	lir.Min *= i
	return lir
}

type spreadIntRange struct {
	Base, Spread int
}

func (sir spreadIntRange) Poll() int {
	return sir.Base + rand.Intn((sir.Spread*2)+1) - sir.Spread
}

func (sir spreadIntRange) Mult(i int) IntRange {
	sir.Base *= i
	sir.Spread *= i
	return sir
}

// Constant implements IntRange as a poll
// which always returns the same integer.
type Constant int

// Poll returns c cast to an int
func (c Constant) Poll() int {
	return int(c)
}

func (c Constant) Mult(i int) IntRange {
	return Constant(int(c) * i)
}

type Infinite struct{}

func (inf Infinite) Poll() int {
	return math.MaxInt32
}

func (inf Infinite) Mult(i int) IntRange {
	return inf
}
