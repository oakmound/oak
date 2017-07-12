package mouse

import "bitbucket.org/oakmoundstudio/oak/collision"

var (
	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent Event
)

// An Event is passed in through all Mouse related event bindings to
// indicate what type of mouse event was triggered, where it was triggered,
// and which mouse button it concerns.
// this is a candidate for merging with physics.Vector
type Event struct {
	X, Y   float32
	Button string
	Event  string
}

// ToSpace converts a mouse event into a collision space
func (e Event) ToSpace() *collision.Space {
	return collision.NewUnassignedSpace(float64(e.X), float64(e.Y), 0.1, 0.1)
}
