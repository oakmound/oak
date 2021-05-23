package floatrange

import "math"

// Infinite is an immutable range that will always return math.MaxFloat64
type Infinite struct{}

// NewInfinite returns an infinite.
func NewInfinite() Range {
	return Infinite{}
}

// Poll returns MaxFloat64 on an infinite
func (i Infinite) Poll() float64 {
	return math.MaxFloat64
}

// Mult returns an infinite from an infinite.
func (i Infinite) Mult(f float64) Range {
	return i
}

// EnforceRange returns math.MaxFloat64
func (i Infinite) EnforceRange(f float64) float64 {
	return math.MaxFloat64
}

// Percentile returns the float behind the constant
func (i Infinite) Percentile(float64) float64 {
	return math.MaxFloat64
}
