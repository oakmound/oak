package main

import (
	"embed"
	"image/color"
	"path"
	"strconv"

	"github.com/oakmound/oak/v3/alg/range/floatrange"

	"image"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var (
	font    *render.Font
	r, g, b float64
	diff    = floatrange.NewSpread(0, 10)
	limit   = floatrange.NewLinear(0, 255)
)

type floatStringer struct {
	f *float64
}

func (fs floatStringer) String() string {
	return strconv.Itoa(int(*fs.f))
}

func main() {
	oak.AddScene("demo",
		scene.Scene{Start: func(ctx *scene.Context) {
			render.Draw(render.NewDrawFPS(0.25, nil, 10, 10))
			fg := render.FontGenerator{
				File:  path.Join("assets", "font", "luxisbi.ttf"),
				Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
				FontOptions: render.FontOptions{
					Size: 50,
					DPI:  72,
				},
			}
			r = 255
			font, _ = fg.Generate()
			font.Unsafe = true
			txts := []*render.Text{
				font.NewText("Rainbow", 200, 200),
				font.NewStringerText(floatStringer{&r}, 200, 260),
				font.NewStringerText(floatStringer{&g}, 320, 260),
				font.NewStringerText(floatStringer{&b}, 440, 260),
			}
			for _, txt := range txts {
				render.Draw(txt, 0)
			}
			font2, _ := font.RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
				fg.Color = image.NewUniform(color.RGBA{255, 255, 255, 255})
				return fg
			})
			render.Draw(font2.NewText("r", 160, 260), 0)
			render.Draw(font2.NewText("g", 280, 260), 0)
			render.Draw(font2.NewText("b", 400, 260), 0)

			ctx.DoEachFrame(func(){
				r = limit.EnforceRange(r + diff.Poll())
				g = limit.EnforceRange(g + diff.Poll())
				b = limit.EnforceRange(b + diff.Poll())
				font.Drawer.Src = image.NewUniform(
					color.RGBA{
						uint8(r),
						uint8(g),
						uint8(b),
						255,
					},
				)
			})
		},
		})
	oak.SetFS(assets)
	oak.Init("demo")
}

//go:embed assets
var assets embed.FS
