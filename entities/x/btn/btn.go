// Package btn provides constructors for UI buttons
package btn

import (
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
)

// Btn defines a button for use in the UI
type Btn interface {
	event.Caller
	render.Positional
	GetRenderable() render.Renderable
	GetSpace() *collision.Space
	SetMetadata(string, string)
	Metadata(string) (string, bool)
	Destroy()
}
