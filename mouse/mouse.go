package mouse

import (
	"github.com/oakmound/oak/v3/event"
	"golang.org/x/mobile/event/mouse"
)

// Button represents a mouse interaction type, like a left button or
// mouse wheel movement.
type Button = mouse.Button

// Valid Button event types
const (
	ButtonLeft       = mouse.ButtonLeft
	ButtonMiddle     = mouse.ButtonMiddle
	ButtonRight      = mouse.ButtonRight
	ButtonWheelDown  = mouse.ButtonWheelDown
	ButtonWheelUp    = mouse.ButtonWheelUp
	ButtonWheelLeft  = mouse.ButtonWheelLeft
	ButtonWheelRight = mouse.ButtonWheelRight
	ButtonNone       = mouse.ButtonNone
)

// GetEventName returns a string event name given some mobile/mouse information
func GetEvent(d mouse.Direction, b mouse.Button) event.EventID[*Event] {
	switch d {
	case mouse.DirPress:
		return Press
	case mouse.DirRelease:
		return Release
	default:
		switch b {
		case -2:
			return ScrollDown
		case -1:
			return ScrollUp
		}
	}
	return Drag
}
