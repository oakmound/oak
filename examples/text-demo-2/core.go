package main

import (
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"image"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

// ~60 fps draw rate with these examples in testing
const (
	strRangeTop = 128
	strlen      = 360
	strSize     = 4
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
	strRange := intrange.NewLinear(0, strRangeTop)
	for i := 0; i < chars; i++ {
		str[i] = rune(strRange.Poll())
	}
	return string(str)
}

func main() {
	oak.Add("demo",
		// Init
		func(prevScene string, payload interface{}) {
			r = 255
			fg := render.FontGenerator{
				File:    "luxisr.ttf",
				Color:   image.NewUniform(color.RGBA{255, 0, 0, 255}),
				Size:    strSize,
				Hinting: "",
			}
			font = fg.Generate()

			for y := 0.0; y <= 480; y += strSize {
				str := randomStr(strlen)
				strs = append(strs, font.NewStrText(str, 0, y))
				render.Draw(strs[len(strs)-1], 0)
			}
		},
		// Loop
		func() bool {
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

		// End
		func() (string, *scene.Result) {
			return "demo", nil
		},
	)
	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
	)
	oak.Init("demo")
}
