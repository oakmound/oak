package oak

import (
	"github.com/oakmound/oak/v2/event"

	"github.com/oakmound/oak/v2/dlog"
	okey "github.com/oakmound/oak/v2/key"
	omouse "github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/shiny/gesture"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

var (
	eFilter gesture.EventFilter
	eventFn func() interface{}
)

func inputLoop() {
	// Obtain input events in a manner dependant on config settings
	if conf.GestureSupport {
		eFilter = gesture.EventFilter{EventDeque: windowControl}
		eventFn = func() interface{} {
			return eFilter.Filter(eFilter.EventDeque.NextEvent())
		}
	} else {
		// Standard input
		eventFn = windowControl.NextEvent
	}
	for {
		switch e := eventFn().(type) {
		// We only currently respond to death lifecycle events.
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				dlog.Info("Window closed.")
				// OnStop needs to be sent through TriggerBack, otherwise the
				// program will close before the stop events get propagated.
				if fh, ok := logicHandler.(event.FullHandler); ok {
					dlog.Verb("Triggering OnStop.")
					<-fh.TriggerBack(event.OnStop, nil)
				}
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
			// key.Code strings all begin with "Code". This strips that off.
			k := GetKeyBind(e.Code.String()[4:])
			switch e.Direction {
			case key.DirPress:
				TriggerKeyDown(k)
			case key.DirRelease:
				TriggerKeyUp(k)
			default:
				TriggerKeyHeld(k)
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
			button := omouse.GetMouseButton(e.Button)
			eventName := omouse.GetEventName(e.Direction, e.Button)
			// The event triggered for mouse events has the same scaling as the
			// render and collision space. I.e. if the viewport is at 0, the mouse's
			// position is exactly the same as the position of a visible entity
			// on screen. When not at zero, the offset will be exactly the viewport.
			// Todo: consider incorporating viewport into the event, see the
			// workaround needed in mouseDetails, and how mouse events might not
			// propagate to their expected position.
			mevent := omouse.NewEvent(
				float64((((e.X - float32(windowRect.Min.X)) / float32(windowRect.Max.X-windowRect.Min.X)) * float32(ScreenWidth))),
				float64((((e.Y - float32(windowRect.Min.Y)) / float32(windowRect.Max.Y-windowRect.Min.Y)) * float32(ScreenHeight))),
				button,
				eventName,
			)
			TriggerMouseEvent(mevent)

		case gesture.Event:
			eventName := "Gesture" + e.Type.String()
			dlog.Verb(eventName)
			logicHandler.Trigger(eventName, omouse.FromShinyGesture(e))

		// There's something called a paint event that we don't respond to

		// Size events update what we scale the screen to
		case size.Event:
			//dlog.Verb("Got size event", e)
			ChangeWindow(e.WidthPx, e.HeightPx)
		}
	}
}

// TriggerKeyDown triggers a software-emulated keypress.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real keypress.
func TriggerKeyDown(k string) {
	SetDown(k)
	logicHandler.Trigger(okey.Down, k)
	logicHandler.Trigger(okey.Down+k, nil)
}

// TriggerKeyUp triggers a software-emulated key release.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key release.
func TriggerKeyUp(k string) {
	SetUp(k)
	logicHandler.Trigger(okey.Up, k)
	logicHandler.Trigger(okey.Up+k, nil)
}

// TriggerKeyHeld triggers a software-emulated key hold signal.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key hold signal.
func TriggerKeyHeld(k string) {
	logicHandler.Trigger(okey.Held, k)
	logicHandler.Trigger(okey.Held+k, nil)
}

// TriggerMouseEvent triggers a software-emulated mouse event.
// This should be used cautiously when the mouse is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key mouse press or movement.
func TriggerMouseEvent(mevent omouse.Event) {
	switch mevent.Event {
	case omouse.Press:
		SetDown(mevent.Button)
	case omouse.Release:
		SetUp(mevent.Button)
	}

	omouse.Propagate(mevent.Event+"On", mevent)
	logicHandler.Trigger(mevent.Event, mevent)
}
