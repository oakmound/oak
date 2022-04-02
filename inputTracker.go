package oak

import (
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/joystick"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
)

// InputType expresses some form of input to the engine to represent a player
type InputType int32

// InputChange is triggered when the most recent input device changes (e.g. keyboard to joystick or vice versa)
var InputChange = event.RegisterEvent[InputType]()

var trackingJoystickChange = event.RegisterEvent[struct{}]()

// Supported Input Types
const (
	InputKeyboard InputType = iota
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

func (jh *joyHandler) Trigger(eventID event.UnsafeEventID, data interface{}) chan struct{} {
	jh.handler.Trigger(trackingJoystickChange.UnsafeEventID, struct{}{})
	ch := make(chan struct{})
	close(ch)
	return ch
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
