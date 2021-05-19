package oak

import (
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/timing"

	"github.com/oakmound/oak/v3/dlog"
	okey "github.com/oakmound/oak/v3/key"
	omouse "github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/shiny/gesture"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

var (
	eFilter gesture.EventFilter
	eventFn func() interface{}
)

func (c *Controller) inputLoop() {
	for {
		switch e := c.windowControl.NextEvent().(type) {
		// We only currently respond to death lifecycle events.
		case lifecycle.Event:
			switch e.To {
			case lifecycle.StageDead:
				dlog.Info("Window closed.")
				// OnStop needs to be sent through TriggerBack, otherwise the
				// program will close before the stop events get propagated.
				dlog.Verb("Triggering OnStop.")
				<-c.logicHandler.TriggerBack(event.OnStop, nil)
				close(c.quitCh)
				return
			case lifecycle.StageFocused:
				c.inFocus = true
				// If you are in focused state, we don't care how you got there
				c.DrawTicker.SetTick(timing.FPSToFrameDelay(c.DrawFrameRate))
				c.logicHandler.Trigger(event.FocusGain, nil)
			case lifecycle.StageVisible:
				// If the last state was focused, this means the app is out of focus
				// otherwise, we're visible for the first time
				if e.From > e.To {
					c.inFocus = false
					c.DrawTicker.SetTick(timing.FPSToFrameDelay(c.IdleDrawFrameRate))
					c.logicHandler.Trigger(event.FocusLoss, nil)
				} else {
					c.inFocus = true
					c.DrawTicker.SetTick(timing.FPSToFrameDelay(c.DrawFrameRate))
					c.logicHandler.Trigger(event.FocusGain, nil)
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
			// TODO v3: reevaluate key bindings-- we need the rune this event has
			switch e.Direction {
			case key.DirPress:
				c.TriggerKeyDown(okey.Event(e))
			case key.DirRelease:
				c.TriggerKeyUp(okey.Event(e))
			default:
				c.TriggerKeyHeld(okey.Event(e))
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
			// Todo: consider incorporating viewport into the event, see the
			// workaround needed in mouseDetails, and how mouse events might not
			// propagate to their expected position.
			mevent := omouse.NewEvent(
				float64((((e.X - float32(c.windowRect.Min.X)) / float32(c.windowRect.Max.X-c.windowRect.Min.X)) * float32(c.ScreenWidth))),
				float64((((e.Y - float32(c.windowRect.Min.Y)) / float32(c.windowRect.Max.Y-c.windowRect.Min.Y)) * float32(c.ScreenHeight))),
				button,
				eventName,
			)
			c.TriggerMouseEvent(mevent)

		// There's something called a paint event that we don't respond to

		// Size events update what we scale the screen to
		case size.Event:
			//dlog.Verb("Got size event", e)
			c.ChangeWindow(e.WidthPx, e.HeightPx)
		}
	}
}

// TriggerKeyDown triggers a software-emulated keypress.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real keypress.
func (c *Controller) TriggerKeyDown(e okey.Event) {
	k := e.Code.String()[4:]
	c.SetDown(k)
	c.logicHandler.Trigger(okey.Down, e)
	c.logicHandler.Trigger(okey.Down+k, e)
}

// TriggerKeyUp triggers a software-emulated key release.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key release.
func (c *Controller) TriggerKeyUp(e okey.Event) {
	k := e.Code.String()[4:]
	c.SetUp(k)
	c.logicHandler.Trigger(okey.Up, e)
	c.logicHandler.Trigger(okey.Up+k, e)
}

// TriggerKeyHeld triggers a software-emulated key hold signal.
// This should be used cautiously when the keyboard is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key hold signal.
func (c *Controller) TriggerKeyHeld(e okey.Event) {
	k := e.Code.String()[4:]
	c.logicHandler.Trigger(okey.Held, e)
	c.logicHandler.Trigger(okey.Held+k, e)
}

// TriggerMouseEvent triggers a software-emulated mouse event.
// This should be used cautiously when the mouse is in use.
// From the perspective of the event handler this is indistinguishable
// from a real key mouse press or movement.
func (c *Controller) TriggerMouseEvent(mevent omouse.Event) {
	c.Propagate(mevent.Event+"On", mevent)
	c.logicHandler.Trigger(mevent.Event, mevent)

	mevent.Point2[0] += float64(c.viewPos[0])
	mevent.Point2[1] += float64(c.viewPos[1])
	c.Propagate(mevent.Event+"OnRelative", mevent)
}
