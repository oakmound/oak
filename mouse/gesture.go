package mouse

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"golang.org/x/exp/shiny/gesture"
	"time"
)

type GestureEvent struct {
	Drag        bool
	LongPress   bool
	DoublePress bool

	InitialPos *physics.Vector
	CurrentPos *physics.Vector

	Time time.Time
}

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
