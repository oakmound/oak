package render

import (
	"image/color"
	"testing"

	"github.com/oakmound/oak/v3/shape"
)

func TestSimpleBezierLine(t *testing.T) {
	bz, err := shape.BezierCurve(0, 0, 10, 10)
	if err != nil {
		t.Fatalf("failed to create bezier curve: %v", err)
	}
	sp := BezierLine(bz, color.RGBA{255, 255, 255, 255})
	rgba := sp.GetRGBA()
	for i := 0; i < 10; i++ {
		if rgba.At(i, i) != (color.RGBA{255, 255, 255, 255}) {
			t.Fatalf("rgba not set at %v", i)
		}
	}

	bz, err = shape.BezierCurve(10, 10, 0, 0)
	if err != nil {
		t.Fatalf("failed to create bezier curve: %v", err)
	}
	sp = BezierLine(bz, color.RGBA{255, 255, 255, 255})
	rgba = sp.GetRGBA()
	for i := 0; i < 10; i++ {
		if rgba.At(i, i) != (color.RGBA{255, 255, 255, 255}) {
			t.Fatalf("rgba not set at %v", i)
		}
	}
}
