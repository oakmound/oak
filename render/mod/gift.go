package mod

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
)

// GiftTransform converts any set of gift.Filters into a Mod.
func GiftTransform(fs ...gift.Filter) Mod {
	return func(rgba image.Image) *image.RGBA {
		filter := gift.New(fs...)
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// GiftFilter converts any set of gift.Filters into a Filter.
// if a filter is internally a transformation in gift, this will
// not work and GiftTransform should be used instead.
func GiftFilter(fs ...gift.Filter) Filter {
	return func(rgba *image.RGBA) {
		gift.New(fs...).Draw(rgba, rgba)
	}
}

// Brighten brightens an image between -100 and 100. 100 will be solid white,
// -100 will be solid black, for all colors not zero before filtering.
func Brighten(brightenBy float32) Filter {
	return GiftFilter(gift.Brightness(brightenBy))
}

// Saturate saturates the input between -100 and 500 percent.
func Saturate(saturateBy float32) Filter {
	return GiftFilter(gift.Saturation(saturateBy))
}

// ColorBalance takes in 3 numbers between -100 and 500 and applies it to the given image
func ColorBalance(r, g, b float32) Filter {
	return GiftFilter(gift.ColorBalance(r, g, b))
}

// Crop will return the given rectangle portion of transformed images. See gift.Crop
func Crop(rect image.Rectangle) Mod {
	return GiftTransform(gift.Crop(rect))
}

// CropToSize applies crop with an optional anchor. See gift.CropToSize
func CropToSize(width, height int, anchor gift.Anchor) Mod {
	return GiftTransform(gift.CropToSize(width, height, anchor))
}

// FlipX returns a new rgba which is flipped
// over the horizontal axis.
var FlipX = GiftTransform(gift.FlipHorizontal())

// FlipY returns a new rgba which is flipped
// over the vertical axis.
var FlipY = GiftTransform(gift.FlipVertical())

// Resize will transform images to match the input dimensions. See gift.Resize.
func Resize(width, height int, resampling gift.Resampling) Mod {
	return GiftTransform(gift.Resize(width, height, resampling))
}

// ResizeToFill will resize to fit and then crop using the given anchor. See gift.ResizeToFill.
func ResizeToFill(width, height int, resampling gift.Resampling, anchor gift.Anchor) Mod {
	return GiftTransform(gift.ResizeToFill(width, height, resampling, anchor))
}

// ResizeToFit will resize while perserving aspect ratio. See gift.ResizeToFit.
func ResizeToFit(width, height int, resampling gift.Resampling) Mod {
	return GiftTransform(gift.ResizeToFit(width, height, resampling))
}

// Rotate returns a rotated rgba.
func Rotate(degrees float32) Mod {
	return RotateInterpolated(degrees, gift.CubicInterpolation)
}

// RotateInterpolated acts as Rotate, but accepts an interpolation argument.
// standard rotation does this with Cubic Interpolation.
func RotateInterpolated(degrees float32, interpolation gift.Interpolation) Mod {
	return RotateBackground(degrees, color.RGBA{0, 0, 0, 0}, interpolation)
}

// RotateBackground acts as RotateInterpolated, but allows for supplying a specific
// background color to the rotation.
func RotateBackground(degrees float32, bckgrnd color.Color, interpolation gift.Interpolation) Mod {
	return GiftTransform(gift.Rotate(degrees, bckgrnd, interpolation))
}

// Rotate180 performs a specialized rotation for 180 degrees.
var Rotate180 = GiftTransform(gift.Rotate180())

// Rotate270 performs a specialized rotation for 270 degrees.
var Rotate270 = GiftTransform(gift.Rotate270())

// Rotate90 performs a specialized rotation for 360 degrees.
var Rotate90 = GiftTransform(gift.Rotate90())

// Transpose flips horizontally and rotates 90 degrees counter clockwise.
var Transpose = GiftTransform(gift.Transpose())

// Transverse flips vertically and rotates 90 degrees counter clockwise.
var Transverse = GiftTransform(gift.Transverse())
