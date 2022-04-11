package oak

import (
	"image/color"

	"github.com/oakmound/oak/v3/shiny/screen"

	"github.com/oakmound/oak/v3/render/mod"
)

// SetPalette tells oak to conform the screen to the input color palette before drawing.
func (w *Window) SetPalette(palette color.Palette) {
	w.SetScreenFilter(mod.ConformToPalette(palette))
}

// SetScreenFilter will filter the screen by the given modification function prior
// to publishing the screen's rgba to be displayed.
func (w *Window) SetScreenFilter(screenFilter mod.Filter) {
	w.prePublish = func(w *Window, tx screen.Texture) {
		screenFilter(w.winBuffers[w.bufferIdx].RGBA())
	}
}

// ClearScreenFilter resets the draw function to no longer filter the screen before
// publishing it to the window.
func (w *Window) ClearScreenFilter() {
	w.prePublish = func(*Window, screen.Texture) {}
}
