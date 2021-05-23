package render

import (
	"image/color"
	"testing"
)

func TestColorBox(t *testing.T) {
	c := color.RGBA{255, 200, 255, 255}
	sp := NewColorBox(5, 5, c)
	rgba := sp.GetRGBA()
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			if rgba.At(x, y) != c {
				t.Fatalf("rgba not set at %v,%v", x, y)
			}
		}
	}
	if rgba.At(6, 6) != (color.RGBA{0, 0, 0, 0}) {
		t.Fatalf("rgba exceeded w/h")
	}
}
