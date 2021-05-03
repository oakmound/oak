package floatrange

import (
	"math/rand"

	"github.com/oakmound/oak/v2/alg/range/internal/random"
)

// NewSpread returns a linear range from base-spread to base+spread
func NewSpread(base, spread float64) Range {
	if spread == 0 {
		return constant(base)
	}
	return linear{
		Min: base - spread,
		Max: base + spread,
		rng: random.Rand(),
	}
}

// NewLinear returns a linear range from min to max
func NewLinear(min, max float64) Range {
	if max == min {
		return constant(min)
	}
	flipped := false
	if max < min {
		max, min = min, max
		flipped = true
	}
	return linear{
		Min:     min,
		Max:     max,
		rng:     random.Rand(),
		flipped: flipped,
	}
}

// linear is a range from min to max
type linear struct {
	Max, Min float64
	rng      *rand.Rand
	flipped  bool
}

// Poll on a linear float range returns a float at uniform
// distribution in lfr's range
func (lfr linear) Poll() float64 {
	return ((lfr.Max - lfr.Min) * lfr.rng.Float64()) + lfr.Min
}

// Mult scales a Linear by f
func (lfr linear) Mult(f float64) Range {
	lfr.Max *= f
	lfr.Min *= f
	return lfr
}

// EnforceRange returns f, if is within the range, or the closest value
// in the range to f.
func (lfr linear) EnforceRange(f float64) float64 {
	if f < lfr.Min {
		return lfr.Min
	} else if f > lfr.Max {
		return lfr.Max
	}
	return f
}

// Percentile returns the fth percentile value along this range
func (lfr linear) Percentile(f float64) float64 {
	return ((lfr.Max - lfr.Min) * f) + lfr.Min
}
