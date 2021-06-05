package show

import (
	"image"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/v3/render"
)

// Loaded Fonts
var (
	// Default size 12
	// Default color white

	Express, _ = (&render.FontGenerator{
		File: fpFilter("expressway rg.ttf"),
	}).Generate()
	Gnuolane, _ = (&render.FontGenerator{
		File: fpFilter("gnuolane rg.ttf"),
	}).Generate()
	Libel, _ = (&render.FontGenerator{
		File: fpFilter("libel-suit-rg.ttf"),
	}).Generate()
)

// FontMod modifies a font
type FontMod func(*render.Font) *render.Font

// FontSet applies
func FontSet(set func(*render.Font)) FontMod {
	return func(f *render.Font) *render.Font {
		f = f.Copy()
		set(f)
		f2, _ := f.Generate()
		return f2
	}
}

//FontSize sets size on a font
func FontSize(size float64) FontMod {
	return FontSet(func(f *render.Font) {
		f.Size = size
	})
}

//FontColor sets the color on a font
func FontColor(c color.Color) FontMod {
	return FontSet(func(f *render.Font) {
		f.Color = image.NewUniform(c)
	})
}

func fpFilter(file string) string {
	return filepath.Join("assets", "font", file)
}
