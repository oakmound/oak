package render

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/stretchr/testify/assert"
)

var (
	widths  = intrange.NewLinear(1, 10)
	heights = intrange.NewLinear(1, 10)
	colors  = colorrange.NewLinear(color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255})
)

// Todo for color boxes, and things that take w/h --
// return an error for negative (or 0 in some cases) w / h. The engine assumes
// right now that the inputs will be valid, which is a mistake
// this is a breaking change for 2.0
func TestColorBoxFuzz(t *testing.T) {
	for i := 0; i < 100; i++ {
		w := widths.Poll()
		h := heights.Poll()
		c := colors.Poll()
		r, g, b, a := c.RGBA()
		fmt.Println(r, g, b, a)
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
