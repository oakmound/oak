package main

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/oak/v4/alg/span"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"

	"image"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

// ~60 fps draw rate with these examples in testing
const (
	strRangeTop = 128
	strlen      = 250
	strSize     = 6
)

var (
	font    *render.Font
	r, g, b float64
	diff    = span.NewSpread(0.0, 10.0)
	limit   = span.NewLinear(0.0, 255.0)
	strs    []*render.Text
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
		scene.Scene{Start: func(ctx *scene.Context) {
			render.Draw(render.NewDrawFPS(.25, nil, 10, 10))

			r = 255
			fg := render.DefFontGenerator
			fg.Color = image.NewUniform(color.RGBA{255, 0, 0, 255})
			fg.FontOptions = render.FontOptions{
				Size: strSize,
			}

			var err error
			font, err = fg.Generate()
			dlog.ErrorCheck(err)
			font.Unsafe = true

			for y := 0.0; y <= 480; y += strSize {
				str := randomStr(strlen)
				strs = append(strs, font.NewText(str, 0, y))
				render.Draw(strs[len(strs)-1], 0)
			}

			event.GlobalBind(ctx, event.Enter, func(_ event.EnterPayload) event.Response {
				r = limit.Clamp(r + diff.Poll())
				g = limit.Clamp(g + diff.Poll())
				b = limit.Clamp(b + diff.Poll())
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
				return 0
			})
		},
		})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("demo")
}
