package mouse

import (
	"time"

	"github.com/oakmound/oak/physics"
	"golang.org/x/exp/shiny/gesture"
)

// A GestureEvent is a conversion of a shiny Gesture to our local type so we
// don't need to import shiny variables in more places.
// GestureEvents contain information about mouse events that are not single actions,
// like drags, holds, and double clicks.
type GestureEvent struct {
	Drag        bool
	LongPress   bool
	DoublePress bool

	InitialPos physics.Vector
	CurrentPos physics.Vector

	Time time.Time
}

// FromShinyGesture converts a shiny gesture.Event to a GestureEvent
func FromShinyGesture(shinyGesture gesture.Event) GestureEvent {

	return GestureEvent{
		shinyGesture.Drag,
		shinyGesture.LongPress,
		shinyGesture.DoublePress,

		physics.NewVector(float64(shinyGesture.InitialPos.X), float64(shinyGesture.InitialPos.Y)),
		physics.NewVector(float64(shinyGesture.CurrentPos.X), float64(shinyGesture.CurrentPos.Y)),

		shinyGesture.Time,
	}
}
