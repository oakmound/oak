package btn

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

// Btn defines a button for use in the UI
type Btn interface {
	event.Caller
	render.Positional
	GetRenderable() render.Renderable
	GetSpace() *collision.Space
}
