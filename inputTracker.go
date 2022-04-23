package oak

import (
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/joystick"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
)

// InputType expresses some form of input to the engine to represent a player
type InputType int32

var trackingJoystickChange = event.RegisterEvent[struct{}]()

// The following constants define valid types of input sent via the InputChange event.
const (
	InputNone InputType = iota
	InputKeyboard
	InputMouse
	InputJoystick
)

func (w *Window) trackInputChanges() {
	event.GlobalBind(w.eventHandler, key.AnyDown, func(key.Event) event.Response {
		old := atomic.SwapInt32(&w.mostRecentInput, int32(InputKeyboard))
		if InputType(old) != InputKeyboard {
			event.TriggerOn(w.eventHandler, InputChange, InputKeyboard)
		}
		return 0
	})
	event.GlobalBind(w.eventHandler, mouse.Press, func(*mouse.Event) event.Response {
		old := atomic.SwapInt32(&w.mostRecentInput, int32(InputMouse))
		if InputType(old) != InputMouse {
			event.TriggerOn(w.eventHandler, InputChange, InputMouse)
		}
		return 0
	})
	event.GlobalBind(w.eventHandler, trackingJoystickChange, func(struct{}) event.Response {
		old := atomic.SwapInt32(&w.mostRecentInput, int32(InputMouse))
		if InputType(old) != InputJoystick {
			event.TriggerOn(w.eventHandler, InputChange, InputJoystick)
		}
		return 0
	})
}

type joyHandler struct {
	handler event.Handler
}

func (jh *joyHandler) Trigger(eventID event.UnsafeEventID, data interface{}) <-chan struct{} {
	return event.TriggerOn(jh.handler, trackingJoystickChange, struct{}{})
}

func trackJoystickChanges(handler event.Handler) {
	dlog.ErrorCheck(joystick.Init())
	go func() {
		jCh, _ := joystick.WaitForJoysticks(3 * time.Second)
		for j := range jCh {
			j.Handler = &joyHandler{
				handler: handler,
			}
			j.Listen(nil)
		}
	}()
}
