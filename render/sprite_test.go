package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/stretchr/testify/assert"
)

var (
	// this is excessive for a lot of tests
	// but it takes away some decision making
	// and could reveal problems that probably aren't there
	// but hey you never know
	widths  = intrange.NewLinear(1, 10)
	heights = intrange.NewLinear(1, 10)
	colors  = colorrange.NewLinear(color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255})
)

const (
	fuzzCt = 10
)

// Todo for color boxes, and things that take w/h --
// return an error for negative (or 0 in some cases) w / h. The engine assumes
// right now that the inputs will be valid, which is a mistake
// this is a breaking change for 2.0
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
				assert.Equal(t, r, r2)
				assert.Equal(t, g, g2)
				assert.Equal(t, b, b2)
				assert.Equal(t, a, a2)
			}
		}
	}
}

// GradientBoxes should use color ranges internally

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
				assert.Equal(t, r, zero)
				assert.Equal(t, g, zero)
				assert.Equal(t, b, zero)
				assert.Equal(t, a, zero)
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
	assert.Equal(t, w, 1)
	assert.Equal(t, h, 1)

	w, h = s2.GetDims()
	assert.Equal(t, w, 6)
	assert.Equal(t, h, 6)

	w, h = s3.GetDims()
	assert.Equal(t, w, 1)
	assert.Equal(t, h, 1)

	// IsNil

	assert.Equal(t, false, s.IsNil())
	assert.Equal(t, true, s2.IsNil())
	assert.Equal(t, false, s3.(*Sprite).IsNil())

	// Set/GetRGBA

	rgba := image.NewRGBA(image.Rect(0, 0, 4, 4))
	s.SetRGBA(rgba)
	rgba2 := s.GetRGBA()
	assert.Equal(t, rgba, rgba2)
}

func TestOverlaySprites(t *testing.T) {
	// This makes me wonder if overlay is easy enough to use
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgba.Set(0, 0, color.RGBA{255, 0, 0, 255})
	// It should probably take in pointers
	sprites := []Sprite{
		*NewColorBox(2, 2, color.RGBA{0, 255, 0, 255}),
		*NewSprite(0, 0, rgba),
	}
	overlay := OverlaySprites(sprites)
	rgba = overlay.GetRGBA()
	shouldRed := rgba.At(0, 0)
	shouldGreen := rgba.At(0, 1)
	assert.Equal(t, shouldRed, color.RGBA{255, 0, 0, 255})
	assert.Equal(t, shouldGreen, color.RGBA{0, 255, 0, 255})
}

// Can't test ParseSubSprite without loading in something for it to return,
// ParseSubSprite also ignores an error for no good reason?
func TestParseSubSprite(t *testing.T) {
	loadedImages["test"] = NewColorBox(100, 100, color.RGBA{255, 0, 0, 255}).GetRGBA()
	sp := ParseSubSprite("test", 0, 0, 25, 25, 0)
	rgba := sp.GetRGBA()
	for x := 0; x < 25; x++ {
		for y := 0; y < 25; y++ {
			c := rgba.At(x, y)
			r, g, b, a := c.RGBA()
			assert.Equal(t, r, uint32(65535))
			assert.Equal(t, g, uint32(0))
			assert.Equal(t, b, uint32(0))
			assert.Equal(t, a, uint32(65535))
		}
	}

}

func TestModifySprite(t *testing.T) {
	s := NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	s2 := s.Modify(Cut(5, 5))
	w, h := s2.GetDims()
	assert.Equal(t, 5, w)
	assert.Equal(t, 5, h)
}

// We'll cover drawing elsewhere
