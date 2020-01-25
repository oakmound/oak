package render

import (
	"image"

	"github.com/oakmound/oak/v2/render/mod"
	//"image/draw"
)

// A Modifiable is a Renderable that has functions to change its
// underlying image.
type Modifiable interface {
	Renderable
	GetRGBA() *image.RGBA
	Modify(...mod.Mod) Modifiable
	Filter(...mod.Filter)
	Copy() Modifiable
}
