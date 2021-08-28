// Package floatrange provides distributions that accept and return float64s.
package floatrange

// Range represents a range of floating point numbers
type Range interface {
	Poll() float64
	Mult(f float64) Range
	EnforceRange(f float64) float64
	Percentile(f float64) float64
}
