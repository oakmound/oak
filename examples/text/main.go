package main

import (
	"embed"
	"fmt"
	"image/color"
	"path"
	"strconv"

	"image"

	findfont "github.com/flopp/go-findfont"
	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/span"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

//go:embed assets
var assets embed.FS

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
			drawFallbackFonts(ctx)
			drawColorChangingText(ctx)
		},
		})
	oak.SetFS(assets)
	oak.Init("demo")
}

func drawFallbackFonts(ctx *scene.Context) {
	const fontHeight = 16

	fg := render.DefFontGenerator
	fg.Color = image.NewUniform(color.RGBA{255, 0, 0, 255})
	fg.FontOptions.Size = fontHeight
	font, _ := fg.Generate()

	fallbackFonts := []string{
		"Arial.ttf",
		"Yumin.ttf",
		// TODO: support multi-color glyphs
		"Seguiemj.ttf",
	}

	for _, fontname := range fallbackFonts {
		fontPath, err := findfont.Find(fontname)
		if err != nil {
			fmt.Println("Do you have ", fontname, "installed?")
			continue
		}
		fg := render.FontGenerator{
			File:  fontPath,
			Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
			FontOptions: render.FontOptions{
				Size: fontHeight,
			},
		}
		fallbackFont, err := fg.Generate()
		if err != nil {
			panic(err)
		}
		font.Fallbacks = append(font.Fallbacks, fallbackFont)
	}

	strings := []string{
		"Latin-lower: abcdefghijklmnopqrstuvwxyz",
		"Latin-upper: ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"Greek-lower: αβγδεζηθικλμνχοπρσςτυφψω",
		"Greek-upper: ΑΒΓΔΕΖΗΘΙΚΛΜΝΧΟΠΡΣΤΥΦΨΩ",
		"Japanese-kana: あいえおうかきけこくはひへほふさしせそすまみめもむ",
		"Kanji: 茂僕私華花日本英雄の時",
		"Emoji: 😀😃😄😁😆😅😂🤣🐶🐱🐭🐹🐰🦊🐻🐼",
	}

	y := 20.0
	for _, str := range strings {
		render.Draw(font.NewText(str, 10, y), 0)
		y += fontHeight
	}
}

func drawColorChangingText(ctx *scene.Context) {
	var (
		r, g, b float64
		diff    = span.NewSpread(0.0, 10.0)
		limit   = span.NewLinear(0.0, 255.0)
	)

	fg := render.FontGenerator{
		File:  path.Join("assets", "font", "luxisbi.ttf"),
		Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
		FontOptions: render.FontOptions{
			Size: 50,
			DPI:  72,
		},
	}
	r = 255
	font, _ := fg.Generate()
	font.Unsafe = true
	texts := []*render.Text{
		font.NewText("Color", 200, 200),
		font.NewStringerText(floatStringer{&r}, 200, 260),
		font.NewStringerText(floatStringer{&g}, 320, 260),
		font.NewStringerText(floatStringer{&b}, 440, 260),
	}
	for _, txt := range texts {
		render.Draw(txt, 0)
	}
	font2, _ := font.RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
		fg.Color = image.NewUniform(color.RGBA{255, 255, 255, 255})
		return fg
	})
	render.Draw(font2.NewText("r", 160, 260), 0)
	render.Draw(font2.NewText("g", 280, 260), 0)
	render.Draw(font2.NewText("b", 400, 260), 0)

	ctx.DoEachFrame(func() {
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
	})
}
