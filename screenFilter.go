package oak

import (
	"image/color"

	"golang.org/x/exp/shiny/screen"

	"github.com/oakmound/oak/render/mod"
)

var (
	// ColorPalette is the current color palette oak is set to conform to. Modification of this
	// value directly will not effect oak's palette, use SetPalette instead. If SetPallete is never called,
	// this is the zero value ([]Color of length 0).
	ColorPalette color.Palette
)

// SetPalette tells oak to conform the screen to the input color palette before drawing.
func SetPalette(palette color.Palette) {
	ColorPalette = palette
	SetScreenFilter(mod.ConformToPalleteInPlace(ColorPalette))
}

// SetScreenFilter will filter the screen by the given modification function prior
// to publishing the screen's rgba to be displayed.
func SetScreenFilter(screenFilter mod.Filter) {
	drawLoopPublish = func(tx screen.Texture) {
		screenFilter(winBuffer.RGBA())
		drawLoopPublishDef(tx)
	}
}

// ClearScreenFilter resets the draw function to no longer filter the screen before
// publishing it to the window.
func ClearScreenFilter() {
	drawLoopPublish = drawLoopPublishDef
}
