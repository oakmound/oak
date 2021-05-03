package main

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/oak/v2/alg/range/floatrange"

	"image"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

// ~60 fps draw rate with these examples in testing
const (
	strRangeTop = 128
	strlen      = 250
	strSize     = 6
)

var (
	font       *render.Font
	r, g, b, a float64
	diff       = floatrange.NewSpread(0, 10)
	limit      = floatrange.NewLinear(0, 255)
	strs       []*render.Text
)

func randomStr(chars int) string {
	str := make([]rune, chars)
	// ascii
	for i := 0; i < chars; i++ {
		str[i] = rune(rand.Intn(strRangeTop))
	}
	return string(str)
}

func main() {
	oak.AddScene("demo",
		scene.Scene{Start: func(*scene.Context) {
			render.Draw(render.NewDrawFPS(.25, nil, 10, 10))

			r = 255
			// By not specifying "File", we use the default
			// font built into the engine
			fg := render.FontGenerator{
				Color:   image.NewUniform(color.RGBA{255, 0, 0, 255}),
				Size:    strSize,
				Hinting: "",
			}
			font, _ = fg.Generate()
			font.Unsafe = true

			for y := 0.0; y <= 480; y += strSize {
				str := randomStr(strlen)
				strs = append(strs, font.NewStrText(str, 0, y))
				render.Draw(strs[len(strs)-1], 0)
			}
		},
			Loop: func() bool {
				r = limit.EnforceRange(r + diff.Poll())
				g = limit.EnforceRange(g + diff.Poll())
				b = limit.EnforceRange(b + diff.Poll())
				// This should be a function in oak to just set color source
				// (or texture source)
				font.Drawer.Src = image.NewUniform(
					color.RGBA{
						uint8(r),
						uint8(g),
						uint8(b),
						255,
					},
				)
				for _, st := range strs {
					st.SetString(randomStr(strlen))
				}
				return true
			},
		})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("demo")
}
