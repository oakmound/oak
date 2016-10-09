package render

import (
	"image"
	"image/color"
	"image/draw"
)

type ColorBox struct {
	Sprite
}

func NewColorBox(w, h int, c color.Color) *ColorBox {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, image.NewUniform(c), image.Point{0, 0}, draw.Src)
	return &ColorBox{
		Sprite{
			LayeredPoint: LayeredPoint{
				Point: Point{
					X: 0.0,
					Y: 0.0,
				},
			},
			r: rgba,
		},
	}
}
