package alg

import (
	"math"
)

// RoundF64 rounds a float to an int
func RoundF64(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	}
	return int(math.Floor(a + 0.5))
}

var (
	ε = 1.0e-7
)

// F64eq uses epsilon equality to compare two float64s
func F64eq(f1, f2 float64) bool {
	return math.Abs(f1-f2) <= ε
}
