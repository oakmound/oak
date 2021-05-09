package oak

import (
	"image/color"

	"github.com/oakmound/oak/v2/shiny/screen"

	"github.com/oakmound/oak/v2/render/mod"
)

// SetPalette tells oak to conform the screen to the input color palette before drawing.
func (c *Controller) SetPalette(palette color.Palette) {
	c.SetScreenFilter(mod.ConformToPallete(palette))
}

// SetScreenFilter will filter the screen by the given modification function prior
// to publishing the screen's rgba to be displayed.
func (c *Controller) SetScreenFilter(screenFilter mod.Filter) {
	c.prePublish = func(c *Controller, tx screen.Texture) {
		screenFilter(c.winBuffer.RGBA())
	}
}

// ClearScreenFilter resets the draw function to no longer filter the screen before
// publishing it to the window.
func (c *Controller) ClearScreenFilter() {
	c.prePublish = func(*Controller, screen.Texture) {}
}
