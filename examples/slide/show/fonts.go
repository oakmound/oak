package show

import (
	"image"
	"image/color"
	"path"

	"github.com/oakmound/oak/v3/render"
)

func InitFonts() (err error) {
	Express, err = (&render.FontGenerator{
		File:  fpFilter("expressway rg.ttf"),
		Color: image.NewUniform(color.RGBA{255, 255, 255, 255}),
	}).Generate()
	if err != nil {
		return
	}
	Gnuolane, err = (&render.FontGenerator{
		File:  fpFilter("gnuolane rg.ttf"),
		Color: image.NewUniform(color.RGBA{255, 255, 255, 255}),
	}).Generate()
	if err != nil {
		return
	}
	Libel, err = (&render.FontGenerator{
		File:  fpFilter("libel-suit-rg.ttf"),
		Color: image.NewUniform(color.RGBA{255, 255, 255, 255}),
	}).Generate()
	if err != nil {
		return
	}
	return nil
}

var (
	Express  *render.Font
	Gnuolane *render.Font
	Libel    *render.Font
)

//FontSize sets size on a font
func FontSize(size float64) func(render.FontGenerator) render.FontGenerator {
	return func(f render.FontGenerator) render.FontGenerator {
		f.Size = size
		return f
	}
}

//FontColor sets the color on a font
func FontColor(c color.Color) func(render.FontGenerator) render.FontGenerator {
	return func(f render.FontGenerator) render.FontGenerator {
		f.Color = image.NewUniform(c)
		return f
	}
}

func fpFilter(file string) string {
	return path.Join("assets", "font", file)
}
