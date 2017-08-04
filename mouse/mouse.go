package mouse

import (
	"github.com/oakmound/oak/collision"
	"golang.org/x/mobile/event/mouse"
)

// Propagate triggers direct mouse events on entities which are clicked
func Propagate(eventName string, me Event) {
	hits := DefTree.SearchIntersect(me.ToSpace().Bounds())
	for _, v := range hits {
		sp := v.(*collision.Space)
		sp.CID.Trigger(eventName, me)
	}
}

// GetMouseButton is a utitilty function which translates
// integer values of mouse keys from golang's event/mouse library
// into strings.
// Intended for internal use.
func GetMouseButton(b mouse.Button) (s string) {
	switch b {
	case mouse.ButtonLeft:
		s = "LeftMouse"
	case mouse.ButtonMiddle:
		s = "MiddleMouse"
	case mouse.ButtonRight:
		s = "RightMouse"
	case mouse.ButtonWheelUp:
		s = "ScrollUpMouse"
	case mouse.ButtonWheelDown:
		s = "ScrollDownMouse"
	default:
		s = ""
	}
	return
}
