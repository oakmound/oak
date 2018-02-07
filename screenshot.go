package oak

import (
	"image"

	"github.com/oakmound/shiny/screen"
)

// ScreenShot takes a snap shot of the window's image content.
// ScreenShot is not safe to call while an existing ScreenShot call has
// yet to finish executing. This could change in the future.
func ScreenShot() *image.RGBA {
	shotCh := make(chan *image.RGBA)
	// We need to take the shot when the screen is not being redrawn
	// We know the screen has everything drawn on it when it is published
	oldPublish := drawLoopPublish
	drawLoopPublish = func(tx screen.Texture) {
		// Copy the buffer
		rgba := winBuffer.RGBA()
		bds := rgba.Bounds()
		copy := image.NewRGBA(bds)
		for x := bds.Min.X; x < bds.Max.X; x++ {
			for y := bds.Min.Y; y < bds.Max.Y; y++ {
				copy.Set(x, y, rgba.RGBAAt(x, y))
			}
		}
		shotCh <- copy
		oldPublish(tx)
	}
	out := <-shotCh
	drawLoopPublish = oldPublish
	return out
}
