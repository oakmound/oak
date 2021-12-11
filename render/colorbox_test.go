package render

import (
	"image"
	"image/color"
	"testing"
	"testing/quick"
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

func TestColorBoxR(t *testing.T) {
	if err := quick.Check(testColorBoxRProperties, nil); err != nil {
		t.Error(err)
	}
}

// w and h are int8 because int16 can make us check 65536*65536 = 4 billion pixels which can time out
func testColorBoxRProperties(r, g, b, a uint8, w8, h8 int8) bool {
	c := color.RGBA{r, g, b, a}
	w := int(w8)
	h := int(h8)
	sp := NewColorBoxR(w, h, c)
	w2, h2 := sp.GetDims()
	if w2 != w {
		return false
	}
	if h2 != h {
		return false
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	sp.Draw(img, 0, 0)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cAt := img.RGBAAt(x, y)
			if c != cAt {
				return false
			}
		}
	}
	return true
}
