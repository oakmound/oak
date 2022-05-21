package span

import (
	"math/rand"

	"github.com/oakmound/oak/v4/alg/span/internal/random"
	"golang.org/x/exp/constraints"
)

// A Spanable must be usable in basic arithmetic-- addition, subtraction, and multiplication.
type Spanable interface {
	constraints.Float | constraints.Integer
}

// NewConstant returns a span where the minimum and maximum are both i. Poll, Percentile, and Clamp will always return i.
func NewConstant[T Spanable](i T) Span[T] {
	return constant[T]{i}
}

type constant[T Spanable] struct {
	val T
}

func (c constant[T]) Poll() T {
	return c.val
}

func (c constant[T]) MulSpan(i float64) Span[T] {
	return constant[T]{T(float64(c.val) * i)}
}

func (c constant[T]) Clamp(T) T {
	return c.val
}

func (c constant[T]) Percentile(float64) T {
	return c.val
}

// NewLinear returns a linear span between min and max. The linearity implies that no point in the span is preferred,
// and Percentile will scale in a constant fashion from min to max.
func NewLinear[T Spanable](min, max T) Span[T] {
	if max == min {
		return constant[T]{min}
	}
	flipped := false
	if max < min {
		max, min = min, max
		flipped = true
	}
	return linear[T]{
		min:     min,
		max:     max,
		rng:     random.Rand(),
		flipped: flipped,
	}
}

// NewSpread returns a linear span from base-spread to base+spread.
func NewSpread[T Spanable](base, spread T) Span[T] {
	if spread < 0 {
		return NewLinear(base+spread, base-spread)
	}
	return NewLinear(base-spread, base+spread)
}

type linear[T Spanable] struct {
	min, max T
	rng      *rand.Rand
	flipped  bool
}

func (lir linear[T]) Poll() T {
	return T(float64(lir.max-lir.min)*lir.rng.Float64()) + lir.min
}

func (lir linear[T]) MulSpan(i float64) Span[T] {
	lir.max = T(float64(lir.max) * i)
	lir.min = T(float64(lir.min) * i)
	return lir
}

func (lir linear[T]) Clamp(i T) T {
	if i < lir.min {
		return lir.min
	} else if i > lir.max {
		return lir.max
	}
	return i
}

func (lir linear[T]) Percentile(f float64) T {
	diff := float64(lir.max-lir.min) * f // 0 - 255 * .1 = -25 + 255 = 230 // 255 - 0 * .1 = 25
	if lir.flipped {
		return lir.max - T(diff)
	}
	return lir.min + T(diff)
}
