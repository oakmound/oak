package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

// Basic Implementing struct: Sprite
type Renderable interface {
	Draw(buff screen.Buffer)
	GetRGBA() *image.RGBA
	// Basic Implementing struct: Point
	ShiftX(x float64)
	ShiftY(y float64)
	SetPos(x, y float64)
	// Basic Implementing struct: Layered
	GetLayer() int
	SetLayer(l int)
	UnDraw()
}
