package render

import (
	"image/color"
	"testing"

	"github.com/oakmound/oak/v4/alg/floatgeom"
)

func TestContains(t *testing.T) {
	p := NewPointsPolygon(
		floatgeom.Point2{10, 10},
		floatgeom.Point2{20, 10},
		floatgeom.Point2{10, 20},
	)
	p.Fill(color.RGBA{255, 0, 0, 255})
	if p.GetRGBA().At(1, 1) != (color.RGBA{255, 0, 0, 255}) {
		t.Fatalf("Fill did not hit 1,1")
	}
	p.FillInverse(color.RGBA{0, 255, 0, 255})
	if p.GetRGBA().At(1, 1) != (color.RGBA{0, 255, 0, 255}) {
		t.Fatalf("FillInverse did not hit 1,1")
	}
}

func TestPolygonFns(t *testing.T) {
	p := NewPointsPolygon(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{0, 10},
		floatgeom.Point2{10, 10},
		floatgeom.Point2{10, 0},
	)
	cmp := p.GetOutline(color.RGBA{255, 0, 0, 255})
	if len(cmp.rs) != 4 {
		t.Fatalf("composite did not contain four lines")
	}
	cmp = p.GetThickOutline(color.RGBA{255, 0, 0, 255}, 1)
	if len(cmp.rs) != 4 {
		t.Fatalf("composite did not contain four lines")
	}
	cmp = p.GetGradientOutline(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, 1)
	if len(cmp.rs) != 4 {
		t.Fatalf("composite did not contain four lines")
	}
}
