package collision

import (
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/oakmound/oak/alg/floatgeom"
)

type SplitTree struct {
	staticTree  *Tree
	dynamicTree *Tree
}

var (
	staticElements  = 1000
	dynamicElements = 10
	xRange          = floatrange.NewLinear(0, 10000)
	yRange          = floatrange.NewLinear(0, 10000)
	wRange          = floatrange.NewLinear(1, 50)
	hRange          = floatrange.NewLinear(1, 50)
	xChange         = floatrange.NewLinear(-3, 3)
	yChange         = floatrange.NewLinear(-3, 3)
)

func randomSpace() *Space {
	return NewUnassignedSpace(xRange.Poll(), yRange.Poll(), wRange.Poll(), hRange.Poll())
}

func BenchmarkSplitTreeHits(b *testing.B) {
	t1, _ := NewTree(2, 4)
	t2, _ := NewTree(20, 40)
	st := SplitTree{
		t1,
		t2,
	}
	dynSpc := []*Space{}

	for i := 0; i < staticElements; i++ {
		st.staticTree.Add(randomSpace())
	}
	for i := 0; i < dynamicElements; i++ {
		s := randomSpace()
		dynSpc = append(dynSpc, s)
		st.dynamicTree.Add(s)
	}
	for i := 0; i < b.N; i++ {
		s := randomSpace()
		st.staticTree.Hits(s)
		st.dynamicTree.Hits(s)
		for _, d := range dynSpc {
			st.dynamicTree.Remove(d)
			p := floatgeom.Point3{xChange.Poll(), yChange.Poll(), 0}
			d.Location.Min.Add(p)
			d.Location.Max.Add(p)
			st.dynamicTree.Add(d)
		}
	}
}

func BenchmarkTreeHits(b *testing.B) {
	t2, _ := NewTree(20, 40)
	dynSpc := []*Space{}

	for i := 0; i < staticElements; i++ {
		t2.Add(randomSpace())
	}
	for i := 0; i < dynamicElements; i++ {
		s := randomSpace()
		dynSpc = append(dynSpc, s)
		t2.Add(s)
	}
	for i := 0; i < b.N; i++ {
		s := randomSpace()
		t2.Hits(s)
		for _, d := range dynSpc {
			t2.Remove(d)
			p := floatgeom.Point3{xChange.Poll(), yChange.Poll(), 0}
			d.Location.Min.Add(p)
			d.Location.Max.Add(p)
			t2.Add(d)
		}
	}
}
