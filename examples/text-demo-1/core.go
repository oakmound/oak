package main

import (
	"image/color"
	"strconv"

	"github.com/200sc/go-dist/floatrange"

	"image"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

var (
	font       *render.Font
	r, g, b, a float64
	diff       = floatrange.NewSpread(0, 10)
	limit      = floatrange.NewLinear(0, 255)
)

type floatStringer struct {
	f *float64
}

func (fs floatStringer) String() string {
	return strconv.Itoa(int(*fs.f))
}

func main() {
	oak.Add("demo",
		// Init
		func(prevScene string, payload interface{}) {
			fg := render.FontGenerator{
				File:    "luxisr.ttf",
				Color:   image.NewUniform(color.RGBA{255, 0, 0, 255}),
				Size:    400,
				Hinting: "",
				DPI:     10,
			}
			r = 255
			font = fg.Generate()
			render.Draw(font.NewStrText("Rainbow", 200, 200), 0)
			render.Draw(font.NewText(floatStringer{&r}, 200, 250), 0)
			render.Draw(font.NewText(floatStringer{&g}, 320, 250), 0)
			render.Draw(font.NewText(floatStringer{&b}, 440, 250), 0)
			font2 := font.Copy()
			font2.Color = image.NewUniform(color.RGBA{255, 255, 255, 255})
			font2.Refresh()
			// Could give each r,g,b a color which is just the r,g,b value
			render.Draw(font2.NewStrText("r", 170, 250), 0)
			render.Draw(font2.NewStrText("g", 290, 250), 0)
			render.Draw(font2.NewStrText("b", 410, 250), 0)
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
