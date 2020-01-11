package btn

import (
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
)

// Box is a basic implementation of btn
type Box struct {
	entities.Solid
	mouse.CollisionPhase
	metadata map[string]string
}

// NewBox creates a new btn.box
func NewBox(cid event.CID, x, y, w, h float64, r render.Renderable, layers ...int) *Box {
	b := Box{}
	cid = cid.Parse(&b)
	b.Solid = *entities.NewSolid(x, y, w, h, r, mouse.DefTree, cid)
	if b.R != nil && len(layers) > 0 {
		render.Draw(b.R, layers...)
	}
	b.metadata = make(map[string]string)
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

// SetMetadata sets the metadata for some key to some value. Empty value strings
// will not be stored.
func (b *Box) SetMetadata(k, v string) {
	if v == "" {
		delete(b.metadata, k)
	} else {
		b.metadata[k] = v
	}
}

// Metadata accesses the value, and whether it existed, for a given metadata key
func (b *Box) Metadata(k string) (v string, ok bool) {
	v, ok = b.metadata[k]
	return v, ok
}
