package oak

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"

	"github.com/oakmound/oak/render"
)

var (
	paletteMod func(*image.RGBA)
	// ColorPalette is the current color palette oak is set to conform to. Modification of this
	// value directly will not effect oak's palette, use SetPalette for that.
	ColorPalette color.Palette
)

// SetPalette tells oak to conform the screen to the input color palette before drawing.
func SetPalette(palette color.Palette) {
	ColorPalette = palette
	paletteMod = render.ConformToPalleteInPlace(ColorPalette)
	drawLoopPublish = func(tx screen.Texture) {
		paletteMod(winBuffer.RGBA())
		drawLoopPublishDef(tx)
	}
}

// ClearPalette stops conforming draw frames to a palette.
func ClearPalette() {
	drawLoopPublish = drawLoopPublishDef
}
