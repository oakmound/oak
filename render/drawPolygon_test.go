package render

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/floatgeom"
)

func TestDrawPolygon(t *testing.T) {
	rh := RenderableHeap{}

	r := rh.DrawPolygonDim()
	if r != (floatgeom.Rect2{floatgeom.Point2{0, 0}, floatgeom.Point2{0, 0}}) {
		t.Fatalf("draw polygon was not a zero to zero rectangle")
	}

	x := 10.0
	y := 10.0
	x2 := 20.0
	y2 := 20.0

	pgn := []floatgeom.Point2{{x, y}, {x, y2}, {x2, y2}, {x2, y}}
	rh.SetPolygon(pgn)

	r = rh.DrawPolygonDim()
	if r != (floatgeom.Rect2{floatgeom.Point2{x, y}, floatgeom.Point2{x2, y2}}) {
		t.Fatalf("draw polygon was not a x,y to x2,y2 rectangle")
	}

	type testcase struct {
		elems         [4]int
		shouldSucceed bool
	}

	tests := []testcase{
		{[4]int{0, 0, 0, 0}, false},
		{[4]int{0, 0, 30, 30}, true},
		{[4]int{15, 15, 17, 17}, true},
	}

	for _, cas := range tests {
		if cas.shouldSucceed != rh.InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]) {
			t.Fatalf("inDrawPolygon failed")
		}
	}

	rh.ClearDrawPolygon()

	for _, cas := range tests {
		if !rh.InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]) {
			t.Fatalf("inDrawPolygon with a cleared polygon failed")
		}
	}
}
