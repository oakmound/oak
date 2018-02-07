package render

import (
	"image/draw"

	"github.com/oakmound/oak/physics"
)

// A Renderable is anything which can
// be drawn at a given draw layer, undrawn,
// and set in a particular position.
//
// Basic Implementing struct: Sprite
type Renderable interface {
	Draw(buff draw.Image)
	DrawOffset(buff draw.Image, xOff, yOff float64)
	GetDims() (int, int)

	Positional
	Layered
	physics.Attachable
}

// Positional types have 2d positions on a screen and can be manipulated
// to be in a certain position on that screen.
//
// Basic Implementing struct: physics.Vector
type Positional interface {
	X() float64
	Y() float64
	ShiftX(x float64)
	ShiftY(y float64)
	SetPos(x, y float64)
}
