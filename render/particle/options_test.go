package particle

import (
	"math"
	"testing"
)

func TestInfiniteLifeSpan(t *testing.T) {
	g := &GradientGenerator{}
	InfiniteLifeSpan()(g)
	if g.LifeSpan.Poll() != math.MaxFloat64 {
		t.Fatalf("Infinite Life Span did not poll math.MaxFloat64")
	}
}
