package mouse

import (
	"github.com/oakmound/oak/event"
	"golang.org/x/mobile/event/mouse"
)

var (
	//TrackMouseClicks enables the propagation of MouseClickOn during MouseRelease events
	TrackMouseClicks = true
)

// Propagate triggers direct mouse events on entities which are clicked
func Propagate(eventName string, me Event) {
	LastEvent = me

	hits := DefTree.SearchIntersect(me.ToSpace().Bounds())
	for _, sp := range hits {
		sp.CID.Trigger(eventName, me)
	}

	if TrackMouseClicks {
		if eventName == PressOn {
			LastPress = me
		} else if eventName == ReleaseOn {
			if me.Button == LastPress.Button {
				pressHits := DefTree.SearchIntersect(LastPress.ToSpace().Bounds())
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							event.Trigger(Click, me)
							sp1.CID.Trigger(ClickOn, me)
						}
					}
				}
			}
		}
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
