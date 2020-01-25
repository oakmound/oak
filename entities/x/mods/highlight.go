package mods

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/render/mod"
)

func HighlightOff(c color.Color, thickness, xOff, yOff int) mod.Mod {
	return func(img image.Image) *image.RGBA {
		bds := img.Bounds()

		w := bds.Max.X + thickness*2 + xOff
		h := bds.Max.Y + thickness*2 + yOff

		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		highlight := image.NewRGBA(image.Rect(0, 0, w, h))

		for x := thickness; x < w-thickness; x++ {
			for y := thickness; y < h-thickness; y++ {
				newRgba.Set(x, y, img.At(x-thickness, y-thickness))
			}
		}
		for x := thickness; x < w-thickness; x++ {
			for y := thickness; y < h-thickness; y++ {
				if _, _, _, a := newRgba.At(x, y).RGBA(); a > 0 {
					for x2 := x - thickness; x2 <= x+thickness; x2++ {
						for y2 := y - thickness; y2 <= y+thickness; y2++ {
							highlight.Set(x2, y2, c)
						}
					}
				}
			}
		}
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				hc := highlight.At(x, y)
				if _, _, _, a := hc.RGBA(); a != 0 {
					if _, _, _, a2 := newRgba.At(x+xOff, y+yOff).RGBA(); a2 == 0 {
						newRgba.Set(x+xOff, y+yOff, hc)
					}
				}
			}
		}
		return newRgba
	}
}

func InnerHighlightOff(c color.Color, thickness, xOff, yOff int) mod.Mod {
	return func(img image.Image) *image.RGBA {
		bds := img.Bounds()

		w := bds.Max.X
		h := bds.Max.Y

		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		highlight := image.NewRGBA(image.Rect(0, 0, w, h))

		for x := thickness; x < w-thickness; x++ {
			for y := thickness; y < h-thickness; y++ {
				newRgba.Set(x, y, img.At(x-thickness, y-thickness))
			}
		}
		for x := thickness; x < w-thickness; x++ {
			for y := thickness; y < h-thickness; y++ {
				if _, _, _, a := newRgba.At(x, y).RGBA(); a == 0 {
					for x2 := x - thickness; x2 <= x+thickness; x2++ {
						for y2 := y - thickness; y2 <= y+thickness; y2++ {
							highlight.Set(x2, y2, c)
						}
					}
				}
			}
		}
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				hc := highlight.At(x, y)
				if _, _, _, a := hc.RGBA(); a != 0 {
					if _, _, _, a2 := newRgba.At(x+xOff, y+yOff).RGBA(); a2 != 0 {
						newRgba.Set(x+xOff, y+yOff, hc) // todo overlay instead
					}
				}
			}
		}
		return newRgba
	}
}

func InnerHighlight(c color.Color, thickness int) mod.Mod {
	return InnerHighlightOff(c, thickness, 0, 0)
}

func Highlight(c color.Color, thickness int) mod.Mod {
	return HighlightOff(c, thickness, 0, 0)
}

type Filter func(color.Color) color.Color

func Inset(fn Filter, dir intgeom.Dir2) mod.Mod {
	return func(img image.Image) *image.RGBA {
		bds := img.Bounds()

		w := bds.Max.X
		h := bds.Max.Y

		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				// todo: depth
				_, _, _, a := img.At(x+dir.X(), y+dir.Y()).RGBA()
				if a == 0 {
					newRgba.Set(x, y, fn(img.At(x, y)))
				} else {
					newRgba.Set(x, y, img.At(x, y))
				}
			}
		}
		return newRgba
	}
}

// Darker produces a darker color by f percentage (0 to 1) difference
func Darker(c color.Color, f float64) color.Color {
	r, g, b, a := c.RGBA()
	diff := uint32(65535 * f)
	r -= diff
	g -= diff
	b -= diff
	// Don't touch alpha
	if r > 65535 {
		r = 0
	}
	if g > 65535 {
		g = 0
	}
	if b > 65535 {
		b = 0
	}
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// Fade produces a color with more transparency by f percentage (0 to 1)
func Fade(c color.Color, f float64) color.Color {
	r, g, b, a := c.RGBA()
	diff := uint32(65535 * f)
	r -= diff
	g -= diff
	b -= diff
	a -= diff
	if r > 65535 {
		r = 0
	}
	if g > 65535 {
		g = 0
	}
	if b > 65535 {
		b = 0
	}
	if a > 65535 {
		a = 0
	}
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}
