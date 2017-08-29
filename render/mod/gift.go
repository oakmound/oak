package mod

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
)

// GiftFilter converts any set of gift.Filters into a Mod.
func GiftFilter(fis ...gift.Filter) Mod {
	return func(rgba image.Image) *image.RGBA {
		filter := gift.New(fis...)
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// Brighten brightens an image between -100 and 100. 100 will be solid white,
// -100 will be solid black.
func Brighten(brightenBy float32) Mod {
	return GiftFilter(gift.Brightness(brightenBy))
}

// Saturate saturates the input between -100 and 500 percent.
func Saturate(saturateBy float32) Mod {
	return GiftFilter(gift.Saturation(saturateBy))
}

// FlipX returns a new rgba which is flipped
// over the horizontal axis.
var FlipX = GiftFilter(gift.FlipHorizontal())

// FlipY returns a new rgba which is flipped
// over the vertical axis.
var FlipY = GiftFilter(gift.FlipVertical())

// ColorBalance takes in 3 numbers between -100 and 500 and applies it to the given image
func ColorBalance(r, g, b float32) Mod {
	return GiftFilter(gift.ColorBalance(r, g, b))
}

// Rotate returns a rotated rgba.
func Rotate(degrees int) Mod {
	return RotateInterpolated(degrees, gift.CubicInterpolation)
}

// RotateInterpolated acts as Rotate, but accepts an interpolation argument.
// standard rotation does this with Cubic Interpolation.
func RotateInterpolated(degrees int, interpolation gift.Interpolation) Mod {
	return GiftFilter(gift.Rotate(float32(degrees), transparent, interpolation))
}

var (
	transparent = color.RGBA{0, 0, 0, 0}
)
