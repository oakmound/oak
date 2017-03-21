package alg

import (
	"math/rand"
)

type FloatRange interface {
	Poll() float64
	Mult(f float64) FloatRange
}

func NewSpreadFloatRange(base, spread float64) FloatRange {
	return spreadFloatRange{base, spread}
}

// SpreadFloatRange is private because SpreadIntRange is private.
// using a negative value for spread is technically legal, because
// it won't cause any crashes, unlike for SpreadIntRange.
type spreadFloatRange struct {
	Base, Spread float64
}

func (sfr spreadFloatRange) Poll() float64 {
	return sfr.Base + (sfr.Spread * 2 * rand.Float64()) - sfr.Spread
}

func (sfr spreadFloatRange) Mult(f float64) FloatRange {
	sfr.Base = sfr.Base * f
	sfr.Spread = sfr.Spread * f
	return sfr
}

type Constantf float64

func (c Constantf) Poll() float64 {
	return float64(c)
}

func (c Constantf) Mult(f float64) FloatRange {
	c = Constantf(float64(c) * f)
	return c
}
