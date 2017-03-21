package alg

import (
	"errors"
	"math/rand"
	"testing"
)

// This introduces a range type which
// does not require defining a struct for
// each range type.

// This benchmark appears to present that, as
// the scale function can not get in-lined,
// this is slower by ~15 percent over
// using separate structs for each
// int range type.
type intRange struct {
	x, y  int
	scale func(int, int) int
}

func (ir intRange) Poll() int {
	return ir.scale(ir.x, ir.y)
}

func (ir intRange) Mult(i int) IntRange {
	ir.x *= i
	ir.y *= i
	return ir
}

func linearInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func baseSpreadInt(base, spread int) int {
	return base + rand.Intn((spread*2)+1) - spread
}

func benchNewLinearIntRange(min, max int) (IntRange, error) {
	if max <= min {
		return Constant(min), errors.New("Max cannot exceed or equal Min, returning constant(Min)")
	}
	return intRange{min, max, linearInt}, nil
}

func benchNewBaseSpreadIntRange(base, spread int) (IntRange, error) {
	if spread <= 0 {
		return Constant(base), errors.New("Spread cannot be <= 0. Returning constant(Base)")
	}
	return intRange{base, spread, baseSpreadInt}, nil
}

func BenchmarkIntRanges(b *testing.B) {
	lin, _ := NewLinearIntRange(0, 100)
	bs, _ := NewSpreadIntRange(0, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lin.Poll()
		bs.Poll()
	}
}

func BenchmarkSameStructIntRanges(b *testing.B) {
	lin, _ := benchNewLinearIntRange(0, 100)
	bs, _ := benchNewBaseSpreadIntRange(0, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lin.Poll()
		bs.Poll()
	}
}
