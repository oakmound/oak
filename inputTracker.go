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
type InputType = int32

// Supported Input Types
const (
	KeyboardMouse InputType = iota
	Joystick      InputType = iota
)

func (c *Controller) trackInputChanges() {
	c.logicHandler.GlobalBind(key.Down, func(event.CID, interface{}) int {
		old := atomic.SwapInt32(&c.mostRecentInput, KeyboardMouse)
		if old != KeyboardMouse {
			c.logicHandler.Trigger(event.InputChange, KeyboardMouse)
		}
		return 0
	})
	c.logicHandler.GlobalBind(mouse.Press, func(event.CID, interface{}) int {
		old := atomic.SwapInt32(&c.mostRecentInput, KeyboardMouse)
		if old != KeyboardMouse {
			c.logicHandler.Trigger(event.InputChange, KeyboardMouse)
		}
		return 0
	})
	c.logicHandler.GlobalBind("Tracking"+joystick.Change, func(event.CID, interface{}) int {
		old := atomic.SwapInt32(&c.mostRecentInput, Joystick)
		if old != Joystick {
			c.logicHandler.Trigger(event.InputChange, Joystick)
		}
		return 0
	})
}

type joyHandler struct {
	handler event.Handler
}

func (jh *joyHandler) Trigger(ev string, state interface{}) {
	jh.handler.Trigger("Tracking"+ev, state)
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
