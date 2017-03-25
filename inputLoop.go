package oak

import (
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	pmouse "bitbucket.org/oakmoundstudio/oak/mouse"
	"golang.org/x/exp/shiny/gesture"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func InputLoop() {
	eFilter := gesture.EventFilter{EventDeque: windowControl}
	for {
		e := eFilter.Filter(eFilter.EventDeque.NextEvent()) //Filters an event to see if it fits a defined gesture

		switch e := e.(type) {

		// We only currently respond to death lifecycle events.
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				quitCh <- true
				return
			}

		// Send key events
		//
		// Key events have two varieties:
		// The "KeyDown" and "KeyUp" events, which trigger for all keys
		// and specific "KeyDown$key", etc events which trigger only for $key.
		// The specific key that is pressed is passed as the data interface for
		// the former events, but not for the latter.
		case key.Event:
			k := GetKeyBind(e.Code.String()[4:])
			if e.Direction == key.DirPress {
				dlog.Verb("--------------------", e.Code.String()[4:], k)
				setDown(k)
				eb.Trigger("KeyDown", k)
				eb.Trigger("KeyDown"+k, nil)
			} else if e.Direction == key.DirRelease {
				setUp(k)
				eb.Trigger("KeyUp", k)
				eb.Trigger("KeyUp"+k, nil)
			}

		// Send mouse events
		//
		// Mouse events are parsed based on their button
		// and direction into an event name and then triggered:
		// 'MousePress', 'MouseRelease', 'MouseScrollDown', 'MouseScrollUp', and 'MouseDrag'
		//
		// The basic event name is meant for entities which
		// want to respond to the mouse event happening -anywhere-.
		//
		// For events which have mouse collision enabled, they'll recieve
		// $eventName+"On" when the event occurs within their collision area.
		//
		// Mouse events all recieve an x, y, and button string.
		case mouse.Event:
			button := pmouse.GetMouseButton(int32(e.Button))
			dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
			mevent := pmouse.MouseEvent{e.X, e.Y, button}
			var eventName string
			if e.Direction == mouse.DirPress {
				setDown(button)
				eventName = "MousePress"
			} else if e.Direction == mouse.DirRelease {
				setUp(button)
				eventName = "MouseRelease"
			} else if e.Button == -2 {
				eventName = "MouseScrollDown"
			} else if e.Button == -1 {
				eventName = "MouseScrollUp"
			} else {
				eventName = "MouseDrag"
			}
			pmouse.LastMouseEvent = mevent

			eb.Trigger(eventName, mevent)
			pmouse.Propagate(eventName+"On", mevent)

		case gesture.Event:
			eventName := "Gesture" + e.Type.String()
			dlog.Verb(eventName)
			eb.Trigger(eventName, pmouse.FromShinyGesture(e))

		// I don't really know what a paint event is to be honest.
		// We hypothetically don't allow the user to manually resize
		// their window, so we don't do anything special for such events.
		case size.Event:
			dlog.Verb("Got size event", e)
		case paint.Event:
			dlog.Verb("Got paint event", e)
		case error:
			dlog.Error(e)
		}

		// This is a hardcoded quit function bound to the escape key.
		esc, dur := IsHeld("Escape")
		if esc && dur > time.Second*1 {
			dlog.Warn("Quiting oak from holding ESCAPE")
			windowControl.Send(lifecycle.Event{0, 0, nil})
		}
	}
}
