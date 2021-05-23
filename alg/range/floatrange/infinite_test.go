package floatrange

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestInfinite(t *testing.T) {
	rand.Seed(time.Now().Unix())
	inf := NewInfinite()
	if inf.Poll() != math.MaxFloat64 {
		t.Fatal("infinite.Poll did not return math.MaxFloat64")
	}
	inf2 := inf.Mult(rand.Float64())
	if inf2 != inf {
		t.Fatal("base infinite did not match multiplied infinite")
	}
	if inf.EnforceRange(rand.Float64()*10000) != math.MaxFloat64 {
		t.Fatal("infinite.EnforceRange did not return math.MaxFloat64")
	}
	if inf.Percentile(rand.Float64()) != math.MaxFloat64 {
		t.Fatal("infinite.Percentile did not return math.MaxFloat64")
	}
}
