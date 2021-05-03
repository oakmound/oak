package intrange

import (
	"math/rand"

	"github.com/oakmound/oak/v2/alg/range/internal/random"
)

// NewLinear returns a linear range between min and max
func NewLinear(min, max int) Range {
	if max == min {
		return constant(min)
	}
	flipped := false
	if max < min {
		max, min = min, max
		flipped = true
	}
	return linear{
		min:     min,
		max:     max,
		rng:     random.Rand(),
		flipped: flipped,
	}
}

// NewSpread returns a linear range from base - s to base + s
func NewSpread(base, spread int) Range {
	if spread == 0 {
		return constant(base)
	}
	if spread < 0 {
		spread *= -1
	}
	return linear{base - spread, base + spread, random.Rand(), false}
}

// linear polls on a linear scale between a minimum and a maximum
type linear struct {
	min, max int
	rng      *rand.Rand
	flipped  bool
}

func (lir linear) Poll() int {
	return int(float64(lir.max-lir.min)*lir.rng.Float64()) + lir.min
}

func (lir linear) Mult(i float64) Range {
	lir.max = int(float64(lir.max) * i)
	lir.min = int(float64(lir.min) * i)
	return lir
}

func (lir linear) EnforceRange(i int) int {
	if i < lir.min {
		return lir.min
	} else if i > lir.max {
		return lir.max
	}
	return i
}

func (lir linear) Percentile(f float64) int {
	diff := float64(lir.max-lir.min) * f // 0 - 255 * .1 = -25 + 255 = 230 // 255 - 0 * .1 = 25
	if lir.flipped {
		return lir.max - int(diff)
	}
	return lir.min + int(diff)
}
