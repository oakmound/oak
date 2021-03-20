package render

import (
	"image"
	"image/draw"
)

// DrawImage performs a draw operation at -x, -y, because
// shiny/screen represents quadrant 4 as negative in both axes.
// draw.Over will merge two pixels at a given position based on their
// alpha channel.
func DrawImage(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Over)
}

// OverwriteImage is equivalent to ShinyDraw, but uses draw.Src
// draw.Src will overwrite pixels beneath the given image regardless of
// the new image's alpha.
func OverwriteImage(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Src)
}
