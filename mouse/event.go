package mouse

import (
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
)

var (
	// LastEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastEvent = NewZeroEvent(0, 0)
	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastPress = NewZeroEvent(0, 0)
)

// An Event is passed in through all Mouse related event bindings to
// indicate what type of mouse event was triggered, where it was triggered,
// and which mouse button it concerns.
type Event struct {
	floatgeom.Point2
	Button
	Event string

	// Set StopPropagation on a mouse event to prevent it from triggering on
	// lower layers of mouse collision spaces while in flight
	StopPropagation bool
}

// NewEvent creates an event.
func NewEvent(x, y float64, button Button, event string) Event {
	return Event{
		Point2: floatgeom.Point2{x, y},
		Button: button,
		Event:  event,
	}
}

// NewZeroEvent creates an event with no button or event name.
func NewZeroEvent(x, y float64) Event {
	return NewEvent(x, y, ButtonNone, "")
}

// ToSpace converts a mouse event into a collision space
func (e Event) ToSpace() *collision.Space {
	sp := collision.NewUnassignedSpace(e.X(), e.Y(), 0.1, 0.1)
	sp.Location.Max[2] = MaxZLayer
	sp.Location.Min[2] = MinZLayer
	return sp
}

// Min and Max Z layer inform what range of z layer values will be checked
// on mouse collision interactions. Mouse events will not propagate to
// elements with z layers outside of this range.
const (
	MinZLayer = 0
	MaxZLayer = 1000
)
