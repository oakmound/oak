package main

import (
	"image/color"
	"strconv"

	"github.com/oakmound/oak/v3/alg/range/floatrange"

	"image"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
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
	oak.AddScene("demo",
		// Init
		scene.Scene{Start: func(*scene.Context) {
			render.Draw(render.NewDrawFPS(0.25, nil, 10, 10))
			// We use the font at ./assets/font/luxisbi.ttf
			// The /assets/font structure is determined by
			// oak.SetupConfig.Assets
			fg := render.FontGenerator{
				File:    "luxisbi.ttf",
				Color:   image.NewUniform(color.RGBA{255, 0, 0, 255}),
				Size:    400,
				Hinting: "",
				DPI:     10,
			}
			r = 255
			font, _ = fg.Generate()
			txts := []*render.Text{
				font.NewText("Rainbow", 200, 200),
				font.NewStringerText(floatStringer{&r}, 200, 250),
				font.NewStringerText(floatStringer{&g}, 320, 250),
				font.NewStringerText(floatStringer{&b}, 440, 250),
			}
			for _, txt := range txts {
				txt.SetFont(font)
				render.Draw(txt, 0)
			}
			font2 := font.Copy()
			font2.Color = image.NewUniform(color.RGBA{255, 255, 255, 255})
			font2, _ = font2.Generate()
			// Could give each r,g,b a color which is just the r,g,b value
			render.Draw(font2.NewText("r", 170, 250), 0)
			render.Draw(font2.NewText("g", 290, 250), 0)
			render.Draw(font2.NewText("b", 410, 250), 0)
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
				return true
			},
		})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("demo")
}
