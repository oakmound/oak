package mouse

import (
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

//TODO V3: should event names be strings?

// GetEventName returns a string event name given some mobile/mouse information
func GetEventName(d mouse.Direction, b mouse.Button) string {
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
