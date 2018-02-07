package show

import (
	"image"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/render"
)

var (
	// Default size 12
	// Default color white
	Express = (&render.FontGenerator{
		File: fpFilter("expressway rg.ttf"),
	}).Generate()
	Gnuolane = (&render.FontGenerator{
		File: fpFilter("gnuolane rg.ttf"),
	}).Generate()
	Libel = (&render.FontGenerator{
		File: fpFilter("libel-suit-rg.ttf"),
	}).Generate()
)

type FontMod func(*render.Font) *render.Font

func FontSet(set func(*render.Font)) FontMod {
	return func(f *render.Font) *render.Font {
		f = f.Copy()
		set(f)
		return f.Generate()
	}
}

func FontSize(size float64) FontMod {
	return FontSet(func(f *render.Font) {
		f.Size = size
	})
}

func FontColor(c color.Color) FontMod {
	return FontSet(func(f *render.Font) {
		f.Color = image.NewUniform(c)
	})
}

// todo: we need to do this because some things
// haven't started in the engine yet (the engine
// doesn't know what our directories are for assets)
// Can we change this?
func fpFilter(file string) string {
	return filepath.Join("assets", "font", file)
}
