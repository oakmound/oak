package render

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v3/alg/range/colorrange"
	"github.com/oakmound/oak/v3/alg/range/intrange"
	"github.com/oakmound/oak/v3/render/mod"
)

var (
	// this is excessive for a lot of tests
	// but it takes away some decision making
	// and could reveal problems that probably aren't there
	widths  = intrange.NewLinear(1, 10)
	heights = intrange.NewLinear(1, 10)
	colors  = colorrange.NewLinear(color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255})
)

const (
	fuzzCt = 10
)

func TestColorBoxFuzz(t *testing.T) {
	for i := 0; i < fuzzCt; i++ {
		w := widths.Poll()
		h := heights.Poll()
		c := colors.Poll()
		r, g, b, a := c.RGBA()
		cb := NewColorBox(w, h, c)
		rgba := cb.GetRGBA()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				c2 := rgba.At(x, y)
				r2, g2, b2, a2 := c2.RGBA()
				if r != r2 {
					t.Fatalf("reds did not match")
				}
				if g != g2 {
					t.Fatalf("greens did not match")
				}
				if b != b2 {
					t.Fatalf("blues did not match")
				}
				if a != a2 {
					t.Fatalf("alphas did not match")
				}
			}
		}
	}
}

// GradientBoxes should use color ranges internally?
func TestGradientBoxFuzz(t *testing.T) {
	for i := 0; i < fuzzCt; i++ {
		w := widths.Poll()
		h := heights.Poll()
		c1 := colors.Poll()
		c2 := colors.Poll()
		r, g, b, a := c1.RGBA()
		r2, g2, b2, a2 := c2.RGBA()
		cb := NewHorizontalGradientBox(w, h, c1, c2)
		rgba := cb.GetRGBA()
		for x := 0; x < w; x++ {
			c3 := rgba.At(x, 0)
			r3, g3, b3, a3 := c3.RGBA()
			progress := float64(x) / float64(w)
			// This sort of color math is frustrating
			c4 := color.RGBA{
				uint8(uint16OnScale(r, r2, progress) / 256),
				uint8(uint16OnScale(g, g2, progress) / 256),
				uint8(uint16OnScale(b, b2, progress) / 256),
				uint8(uint16OnScale(a, a2, progress) / 256),
			}
			r4, g4, b4, a4 := c4.RGBA()
			if r3 != r4 {
				t.Fatalf("reds did not match")
			}
			if g3 != g4 {
				t.Fatalf("greens did not match")
			}
			if b3 != b4 {
				t.Fatalf("blues did not match")
			}
			if a3 != a4 {
				t.Fatalf("alphas did not match")
			}
		}
		cb = NewVerticalGradientBox(w, h, c1, c2)
		rgba = cb.GetRGBA()
		for y := 0; y < h; y++ {
			c3 := rgba.At(0, y)
			r3, g3, b3, a3 := c3.RGBA()
			progress := float64(y) / float64(h)
			// This sort of color math is frustrating
			c4 := color.RGBA{
				uint8(uint16OnScale(r, r2, progress) / 256),
				uint8(uint16OnScale(g, g2, progress) / 256),
				uint8(uint16OnScale(b, b2, progress) / 256),
				uint8(uint16OnScale(a, a2, progress) / 256),
			}
			r4, g4, b4, a4 := c4.RGBA()
			if r3 != r4 {
				t.Fatalf("reds did not match")
			}
			if g3 != g4 {
				t.Fatalf("greens did not match")
			}
			if b3 != b4 {
				t.Fatalf("blues did not match")
			}
			if a3 != a4 {
				t.Fatalf("alphas did not match")
			}
		}
		cb = NewCircularGradientBox(w, h, c1, c2)
		rgba = cb.GetRGBA()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				c3 := rgba.At(x, y)
				r3, g3, b3, a3 := c3.RGBA()
				progress := CircularProgress(x, y, w, h)
				// This sort of color math is frustrating
				c4 := color.RGBA{
					uint8(uint16OnScale(r, r2, progress) / 256),
					uint8(uint16OnScale(g, g2, progress) / 256),
					uint8(uint16OnScale(b, b2, progress) / 256),
					uint8(uint16OnScale(a, a2, progress) / 256),
				}
				r4, g4, b4, a4 := c4.RGBA()
				if r3 != r4 {
					t.Fatalf("reds did not match")
				}
				if g3 != g4 {
					t.Fatalf("greens did not match")
				}
				if b3 != b4 {
					t.Fatalf("blues did not match")
				}
				if a3 != a4 {
					t.Fatalf("alphas did not match")
				}
			}
		}
	}
}

func TestEmptySpriteFuzz(t *testing.T) {
	for i := 0; i < fuzzCt; i++ {
		w := widths.Poll()
		h := heights.Poll()
		s := NewEmptySprite(0, 0, w, h)
		rgba := s.GetRGBA()
		var zero uint32
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				c := rgba.At(x, y)
				r, g, b, a := c.RGBA()
				if r != zero {
					t.Fatalf("reds did not match")
				}
				if g != zero {
					t.Fatalf("greens did not match")
				}
				if b != zero {
					t.Fatalf("blues did not match")
				}
				if a != zero {
					t.Fatalf("alphas did not match")
				}
			}
		}
	}
}

func TestSpriteFuncs(t *testing.T) {
	s := NewEmptySprite(0, 0, 1, 1)
	s2 := Sprite{}
	s3 := s.Copy()

	// Dims

	w, h := s.GetDims()
	if w != 1 || h != 1 {
		t.Fatalf("get dims failed")
	}

	w, h = s2.GetDims()
	if w != 1 || h != 1 {
		t.Fatalf("get dims failed")
	}

	w, h = s3.GetDims()
	if w != 1 || h != 1 {
		t.Fatalf("get dims failed")
	}

	// Set/GetRGBA

	rgba := image.NewRGBA(image.Rect(0, 0, 4, 4))
	s.SetRGBA(rgba)
	rgba2 := s.GetRGBA()
	if !reflect.DeepEqual(rgba, rgba2) {
		t.Fatalf("sprite set rgba failed")
	}
}

func TestOverlaySprites(t *testing.T) {
	// This makes me wonder if overlay is easy enough to use
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgba.Set(0, 0, color.RGBA{255, 0, 0, 255})
	overlay := OverlaySprites([]*Sprite{
		NewColorBox(2, 2, color.RGBA{0, 255, 0, 255}),
		NewSprite(0, 0, rgba),
	})
	rgba = overlay.GetRGBA()
	shouldRed := rgba.At(0, 0)
	shouldGreen := rgba.At(0, 1)
	if shouldRed != (color.RGBA{255, 0, 0, 255}) {
		t.Fatalf("red was not red")
	}
	if shouldGreen != (color.RGBA{0, 255, 0, 255}) {
		t.Fatalf("tgreen was not green")
	}
}

func TestModifySprite(t *testing.T) {
	s := NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	s2 := s.Modify(mod.Cut(5, 5))
	w, h := s2.GetDims()
	if w != 5 || h != 5 {
		t.Fatalf("get dims failed")
	}
}

func TestSprite_ColorModel(t *testing.T) {
	s := NewEmptySprite(0, 0, 1, 1)
	if s.ColorModel() != color.RGBAModel {
		t.Fatalf("color model did not match expected")
	}
}
