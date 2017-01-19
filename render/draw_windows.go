// +build windows

package render

import (
	"image"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

// ShinyDraw performs a draw operation at -x, -y, because
// shiny/screen represents quadrant 4 as negative in both axes.
// draw.Over will merge two pixels at a given position based on their
// alpha channel.
func ShinyDraw(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Over)
}

// draw.Src will overwrite pixels beneath the given image regardless of
// the new image's alpha.
func ShinyOverwrite(buff screen.Buffer, img image.Image, x, y int) {
	draw.Draw(buff.RGBA(), buff.Bounds(),
		img, image.Point{-x, -y}, draw.Src)
}
