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
	// As the engine currently exists,
	// the buffer which is passed into draw
	// is always the same. This leads to
	// several parts of the engine being
	// reliant on shiny/screen when they
	// could call out to somewhere else to
	// determine what they are drawn onto.
	//
	// On the other hand, this allows manually
	// duplicating renderables onto multiple
	// buffers, but in certain implementations
	// (i.e Animation) would have unintended
	// consequences.
	Draw(buff draw.Image)
	DrawOffset(buff draw.Image, xOff, yOff float64)

	// Basic Implementing struct: Point
	ShiftX(x float64)
	GetX() float64
	ShiftY(y float64)
	GetY() float64
	SetPos(x, y float64)
	GetDims() (int, int)

	// Basic Implementing struct: Layered
	GetLayer() int
	SetLayer(l int)
	UnDraw()

	// Utilities
	String() string

	// Physics
	// Basic Implementing struct: physics.Vector
	physics.Attachable
}
