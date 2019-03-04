package btn

import (
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

// Box is a basic implementation of btn
type Box struct {
	entities.Solid
	mouse.CollisionPhase
}

// NewBox creates a new btn.box
func NewBox(cid event.CID, x, y, w, h float64, r render.Renderable, layers ...int) *Box {
	b := Box{}
	cid = cid.Parse(&b)
	b.Solid = *entities.NewSolid(x, y, w, h, r, mouse.DefTree, cid)
	if b.R != nil && len(layers) > 0 {
		render.Draw(b.R, layers...)
	}
	return &b
}

// Init intializes the btn.box
func (b *Box) Init() event.CID {
	b.CID = event.NextID(b)
	return b.CID
}

// GetRenderable returns the box's renderable
func (b *Box) GetRenderable() render.Renderable {
	return b.R
}
