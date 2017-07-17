package oak

import (
	"image"
	"runtime"

	"github.com/oakmound/oak/dlog"
	pmouse "github.com/oakmound/oak/mouse"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

// Todo: bring this back in as an option, it was remove to reduce input lag
// var (
// 	eFilter gesture.EventFilter
// )

func inputLoop() {
	//eFilter = gesture.EventFilter{EventDeque: windowControl}
	schedCt := 0
	for {
		//e := eFilter.Filter(eFilter.EventDeque.NextEvent()) //Filters an event to see if it fits a defined gesture
		switch e := windowControl.NextEvent().(type) {
		// We only currently respond to death lifecycle events.
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				quitCh <- true
				return
			}
			// ... this is where we would respond to window focus events

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
				//dlog.Verb("--------------------", e.Code.String()[4:], k)
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
		// For events which have mouse collision enabled, they'll receive
		// $eventName+"On" when the event occurs within their collision area.
		//
		// Mouse events all receive an x, y, and button string.
		case mouse.Event:
			button := pmouse.GetMouseButton(int32(e.Button))
			//dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
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
			mevent := pmouse.Event{
				X:      e.X / float32(windowRect.Max.X) * float32(ScreenWidth),
				Y:      e.Y / float32(windowRect.Max.Y) * float32(ScreenHeight),
				Button: button,
				Event:  eventName,
			}

			pmouse.LastMouseEvent = mevent

			eb.Trigger(eventName, mevent)
			pmouse.Propagate(eventName+"On", mevent)

		// Uncomment this if using the filter
		// case gesture.Event:
		// 	eventName := "Gesture" + e.Type.String()
		// 	dlog.Verb(eventName)
		// 	eb.Trigger(eventName, pmouse.FromShinyGesture(e))

		// I don't really know what a paint event is to be honest.
		// We hypothetically don't allow the user to manually resize
		// their window, so we don't do anything special for such events.
		case size.Event:
			//dlog.Verb("Got size event", e)
			windowRect = image.Rect(0, 0, e.WidthPx, e.HeightPx)
		case error:
			dlog.Error(e)
		}
		// This loop can be tight enough that the go scheduler never gets
		// a chance to take control from this thread. This is a hack that
		// solves that.
		schedCt++
		if schedCt > 1000 {
			schedCt = 0
			runtime.Gosched()
		}
		/*
			//TODO: Reimplement outside of the input loop so that it doesnt slow down the input loop itself
				// This is a hardcoded quit function bound to the escape key.
				esc, dur := IsHeld("Escape")
				if esc && dur > time.Second*1 {
					dlog.Warn("Quiting oak from holding ESCAPE")
					windowControl.Send(lifecycle.Event{
						From:        0,
						To:          0,
						DrawContext: nil,
					})
				}
		*/
	}
}
