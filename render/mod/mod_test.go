package mod

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllModifications(t *testing.T) {
	in := setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255})
	type filterCase struct {
		Filter
		*image.RGBA
	}
	filterList := []filterCase{{
		ConformToPallete(color.Palette{color.RGBA{64, 0, 0, 128}}),
		setAll(newrgba(3, 3), color.RGBA{64, 0, 0, 128}),
	}, {
		Fade(10),
		setAll(newrgba(3, 3), color.RGBA{245, 0, 0, 245}),
	}, {
		Fade(500),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		FillMask(*setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255})),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		AndFilter(Fade(500), FillMask(*setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255}))),
		setAll(newrgba(3, 3), color.RGBA{0, 255, 0, 255}),
	}, {
		AndFilter(Fade(500), ApplyColor(color.RGBA{255, 255, 255, 255})),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		ApplyColor(color.RGBA{255, 255, 255, 255}),
		setAll(newrgba(3, 3), color.RGBA{127, 127, 127, 255}),
	}, {
		ApplyMask(*setAll(newrgba(3, 3), color.RGBA{255, 255, 255, 255})),
		setAll(newrgba(3, 3), color.RGBA{127, 127, 127, 255}),
	}, {
		AndFilter(Fade(500), ApplyMask(*setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}))),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 0}),
	}, {
		Brighten(-100),
		setAll(newrgba(3, 3), color.RGBA{0, 0, 0, 255}),
	}, {
		Saturate(-100),
		setAll(newrgba(3, 3), color.RGBA{128, 128, 128, 255}),
	}, {
		ColorBalance(0, 0, 0),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		InPlace(Scale(2, 2)),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}}
	for _, f := range filterList {
		in2 := copyrgba(in)
		f.Filter(in2)
		assert.Equal(t, in2, f.RGBA)
	}
	type modCase struct {
		Mod
		*image.RGBA
	}
	modList := []modCase{{
		TrimColor(color.RGBA{255, 0, 0, 255}),
		newrgba(0, 0),
	}, {
		TrimColor(color.RGBA{0, 0, 0, 0}),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		Zoom(1.0, 1.0, 1.0),
		setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
	}, {
		Cut(1, 1),
		setAll(newrgba(1, 1), color.RGBA{255, 0, 0, 255}),
	}, {
		CutRel(.66, .66),
		setAll(newrgba(2, 2), color.RGBA{255, 0, 0, 255}),
	}, {
		CutFromLeft(2, 2),
		setAll(newrgba(2, 2), color.RGBA{255, 0, 0, 255}),
	}, {
		CutRound(.5, .5),
		setOne(setOne(setOne(setOne(setOne(
			setAll(newrgba(3, 3), color.RGBA{255, 0, 0, 255}),
			color.RGBA{0, 0, 0, 0}, 0, 0),
			color.RGBA{0, 0, 0, 0}, 1, 0),
			color.RGBA{0, 0, 0, 0}, 2, 0),
			color.RGBA{0, 0, 0, 0}, 0, 1),
			color.RGBA{0, 0, 0, 0}, 0, 2),
	}}

	for _, m := range modList {
		assert.Equal(t, m.Mod(in), m.RGBA)
	}
}

// test utils

func copyrgba(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	rgba2 := newrgba(bounds.Max.X, bounds.Max.Y)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rgba2.Set(x, y, rgba.At(x, y))
		}
	}
	return rgba2
}

func newrgba(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func setAll(rgba *image.RGBA, c color.Color) *image.RGBA {
	bounds := rgba.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rgba.Set(x, y, c)
		}
	}
	return rgba
}

func setOne(rgba *image.RGBA, c color.Color, x, y int) *image.RGBA {
	rgba.Set(x, y, c)
	return rgba
}
