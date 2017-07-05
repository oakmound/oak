package shape

import "math"

// Eq represents a basic equation-- a mapping of x values to
// y values. Specifically, this equation is expected to be significant
// to represent some part or all of a shape from -1 to 1. This range
// is chosen because it's often easier to write shape equations around
// the center of a graph.
type Eq func(x float64) (y float64)

var (
	// Top half of heart
	hf1 Eq = func(x float64) (y float64) {
		return -2.2*math.Pow(.4+x, 2) + 1
	}
	hf2 Eq = func(x float64) (y float64) {
		return hf1(-x)
	}
	// Bottom half of heart
	hf3 Eq = func(x float64) (y float64) {
		return -math.Sqrt((x + 1)) + .2
	}
	hf4 Eq = func(x float64) (y float64) {
		return hf3(-x)
	}
)

// Below returns an In which reports true for all x,y coordinates below
// the equation curve.
func (eq Eq) Below() In {
	return func(x, y int, sizes ...int) bool {
		w := sizes[0]
		h := sizes[0]
		if len(sizes) > 1 {
			h = sizes[1]
		}
		// shift from 0 to size to -1 to 1
		x2 := float64(x-w/2) / float64(w/2)
		y2 := (float64(y-h/2) / float64(h/2)) * -1
		return eq(x2) > y2
	}
}

// Above returns an In which reports true for all x,y coordinates above
// the equation curve.
func (eq Eq) Above() In {
	return func(x, y int, sizes ...int) bool {
		w := sizes[0]
		h := sizes[0]
		if len(sizes) > 1 {
			h = sizes[1]
		}
		x2 := float64(x-w/2) / float64(w/2)
		y2 := (float64(y-h/2) / float64(h/2)) * -1
		return eq(x2) < y2
	}
}

var (
	// Heart
	// . . t . t . .
	// . t t t t t .
	// t t t t t t t
	// t t t t t t t
	// . t t t t t .
	// . . t t t . .
	// . . . . . . .
	Heart = JustIn(OrIn(
		AndIn(
			XRange(0, 0.5), hf1.Below(), hf3.Above()),
		AndIn(
			XRange(0.5, 1), hf2.Below(), hf4.Above()),
	))
)
