package oak

import (
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/timing"

	"github.com/oakmound/oak/v3/dlog"
	okey "github.com/oakmound/oak/v3/key"
	omouse "github.com/oakmound/oak/v3/mouse"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

func (w *Window) inputLoop() {
	for {
		switch e := w.windowControl.NextEvent().(type) {
		// We only currently respond to death lifecycle events.
		case lifecycle.Event:
			switch e.To {
			case lifecycle.StageDead:
				dlog.Info(dlog.WindowClosed)
				// OnStop needs to be sent through TriggerBack, otherwise the
				// program will close before the stop events get propagated.
				<-w.eventHandler.TriggerBack(event.OnStop, nil)
				close(w.quitCh)
				return
			case lifecycle.StageFocused:
				w.inFocus = true
				// If you are in focused state, we don't care how you got there
				w.DrawTicker.Reset(timing.FPSToFrameDelay(w.DrawFrameRate))
				w.eventHandler.Trigger(event.FocusGain, nil)
			case lifecycle.StageVisible:
				// If the last state was focused, this means the app is out of focus
				// otherwise, we're visible for the first time
				if e.From > e.To {
					w.inFocus = false
					w.DrawTicker.Reset(timing.FPSToFrameDelay(w.IdleDrawFrameRate))
					w.eventHandler.Trigger(event.FocusLoss, nil)
				} else {
					w.inFocus = true
					w.DrawTicker.Reset(timing.FPSToFrameDelay(w.DrawFrameRate))
					w.eventHandler.Trigger(event.FocusGain, nil)
				}
			}
		// Send key events
		//
		// Key events have two varieties:
		// The "KeyDown" and "KeyUp" events, which trigger for all keys
		// and specific "KeyDown$key", etc events which trigger only for $key.
		// The specific key that is pressed is passed as the data interface for
		// the former events, but not for the latter.
		case key.Event:
			switch e.Direction {
			case key.DirPress:
				w.TriggerKeyDown(okey.Event(e))
			case key.DirRelease:
				w.TriggerKeyUp(okey.Event(e))
			default:
				w.TriggerKeyHeld(okey.Event(e))
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
			button := omouse.Button(e.Button)
			eventName := omouse.GetEventName(e.Direction, e.Button)
			// The event triggered for mouse events has the same scaling as the
			// render and collision space. I.e. if the viewport is at 0, the mouse's
			// position is exactly the same as the position of a visible entity
			// on screen. When not at zero, the offset will be exactly the viewport.
			mevent := omouse.NewEvent(
				float64((((e.X - float32(w.windowRect.Min.X)) / float32(w.windowRect.Max.X-w.windowRect.Min.X)) * float32(w.ScreenWidth))),
				float64((((e.Y - float32(w.windowRect.Min.Y)) / float32(w.windowRect.Max.Y-w.windowRect.Min.Y)) * float32(w.ScreenHeight))),
				button,
				eventName,
			)
			w.TriggerMouseEvent(mevent)

		// Size events update what we scale the screen to
		case size.Event:
			if e.WidthPx == 0 || e.HeightPx == 0 {
				// The window has likely been minimized
				continue
			}
			w.eventHandler.Trigger(WindowSizeChange, intgeom.Point2{e.WidthPx, e.HeightPx})
			err := w.ChangeWindow(e.WidthPx, e.HeightPx)
			dlog.ErrorCheck(err)
		}
	}
}

const WindowSizeChange = "WindowSizeChange"

func SizeChangeEvent(f func(c event.CID, pt intgeom.Point2) int) event.Bindable {
	return func(c event.CID, ptData interface{}) int {
		pt, ok := ptData.(intgeom.Point2)
		if !ok {
			return event.UnbindSingle
		}
		return f(c, pt)
	}
}

// TriggerKeyDown triggers a software-emulated keypress.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real keypress.
func (w *Window) TriggerKeyDown(e okey.Event) {
	k := e.Code.String()[4:]
	w.SetDown(k)
	w.eventHandler.Trigger(okey.Down, e)
	w.eventHandler.Trigger(okey.Down+k, e)
}

// TriggerKeyUp triggers a software-emulated key release.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key release.
func (w *Window) TriggerKeyUp(e okey.Event) {
	k := e.Code.String()[4:]
	w.SetUp(k)
	w.eventHandler.Trigger(okey.Up, e)
	w.eventHandler.Trigger(okey.Up+k, e)
}

// TriggerKeyHeld triggers a software-emulated key hold signal.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key hold signal.
func (w *Window) TriggerKeyHeld(e okey.Event) {
	k := e.Code.String()[4:]
	w.eventHandler.Trigger(okey.Held, e)
	w.eventHandler.Trigger(okey.Held+k, e)
}

// TriggerMouseEvent triggers a software-emulated mouse event.
// This should be used cautiously when the mouse is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key mouse press or movement.
func (w *Window) TriggerMouseEvent(mevent omouse.Event) {
	w.LastMouseEvent = mevent
	omouse.LastEvent = mevent
	w.Propagate(mevent.Event+"On", mevent)
	w.eventHandler.Trigger(mevent.Event, &mevent)

	relativeEvent := mevent
	relativeEvent.Point2[0] += float64(w.viewPos[0])
	relativeEvent.Point2[1] += float64(w.viewPos[1])
	w.LastRelativeMouseEvent = relativeEvent
	w.Propagate(relativeEvent.Event+"OnRelative", relativeEvent)
}
