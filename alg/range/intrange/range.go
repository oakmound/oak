// Package intrange holds distributions that return ints
package intrange

// Range represents a range of integer numbers
type Range interface {
	Poll() int
	Mult(float64) Range
	EnforceRange(int) int
	Percentile(float64) int
}
