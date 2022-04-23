package mouse

import "github.com/oakmound/oak/v4/event"

var (
	// Press is triggered when a mouse key is pressed down
	Press = event.RegisterEvent[*Event]()
	// Release is triggered when a mouse key, pressed, is released
	Release = event.RegisterEvent[*Event]()
	// ScrollDown is triggered when a mouse's scroll wheel scrolls downward
	ScrollDown = event.RegisterEvent[*Event]()
	// ScrollUp is triggered when a mouse's scroll wheel scrolls upward
	ScrollUp = event.RegisterEvent[*Event]()
	// Click is triggered when a Release follows a press for the same mouse key without
	// other mouse key presses intertwining.
	Click = event.RegisterEvent[*Event]()
	// Drag is triggered when the mouse is moved.
	Drag = event.RegisterEvent[*Event]()

	// The 'On' Variants of all mouse events are triggered when a mouse event occurs on
	// a specific entity in a mouse collision tree.
	PressOn      = event.RegisterEvent[*Event]()
	ReleaseOn    = event.RegisterEvent[*Event]()
	ScrollDownOn = event.RegisterEvent[*Event]()
	ScrollUpOn   = event.RegisterEvent[*Event]()
	ClickOn      = event.RegisterEvent[*Event]()
	DragOn       = event.RegisterEvent[*Event]()

	// Relative variants are like 'On' variants, but their mouse position data is relative to
	// the window's current viewport. E.g. if the viewport is at 100,100 and a click happens at
	// 100,100 on the window-- Relative will report 100,100, and non-relative will report 200,200.
	// TODO: re-evaluate relative vs non-relative mouse events
	RelativePressOn      = event.RegisterEvent[*Event]()
	RelativeReleaseOn    = event.RegisterEvent[*Event]()
	RelativeScrollDownOn = event.RegisterEvent[*Event]()
	RelativeScrollUpOn   = event.RegisterEvent[*Event]()
	RelativeClickOn      = event.RegisterEvent[*Event]()
	RelativeDragOn       = event.RegisterEvent[*Event]()
)

// EventOn converts a generic positioned mouse event into its variant indicating
// it occurred on a CallerID targetted entity
func EventOn(ev event.EventID[*Event]) (event.EventID[*Event], bool) {
	switch ev {
	case Press:
		return PressOn, true
	case Release:
		return ReleaseOn, true
	case ScrollDown:
		return ScrollDownOn, true
	case ScrollUp:
		return ScrollUpOn, true
	case Click:
		return ClickOn, true
	case Drag:
		return DragOn, true
	}
	return event.EventID[*Event]{}, false
}

func EventRelative(ev event.EventID[*Event]) (event.EventID[*Event], bool) {
	switch ev {
	case PressOn:
		return RelativePressOn, true
	case ReleaseOn:
		return RelativeReleaseOn, true
	case ScrollDownOn:
		return RelativeScrollDownOn, true
	case ScrollUpOn:
		return RelativeScrollUpOn, true
	case ClickOn:
		return RelativeClickOn, true
	case DragOn:
		return RelativeDragOn, true
	}
	return event.EventID[*Event]{}, false
}
