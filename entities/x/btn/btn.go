package btn

import (
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
)

// Btn defines a button for use in the UI
type Btn interface {
	event.Caller
	render.Positional
	GetRenderable() render.Renderable
	GetSpace() *collision.Space
	SetMetadata(string, string)
	Metadata(string) (string, bool)
}
