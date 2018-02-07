package alg

import (
	"math"
)

// RoundF64 rounds a float64 to an int
func RoundF64(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	}
	return int(math.Floor(a + 0.5))
}

const (
	ε = 1.0e-7
)

// F64eq equates two float64s within a small epsilon.
func F64eq(f1, f2 float64) bool {
	return F64eqEps(f1, f2, ε)
}

// F64eqEps equates two float64s within a provided epsilon.
func F64eqEps(f1, f2, epsilon float64) bool {
	return math.Abs(f1-f2) <= epsilon
}
