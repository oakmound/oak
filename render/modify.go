package render

import
// This file is being slowly converted to use gift over manual math and loops,
// because our math / loops will be more likely to have (and have already had)
// missable bugs.
//"github.com/anthonynsimon/bild/blend"

(
	"image"

	"github.com/oakmound/oak/render/mod"
	//"image/draw"
)

// A Modifiable is a Renderable that has functions to change its
// underlying image.
// This may be replaced with the gift library down the line
type Modifiable interface {
	Renderable
	GetRGBA() *image.RGBA
	Modify(...mod.Mod) Modifiable
	Copy() Modifiable
}
