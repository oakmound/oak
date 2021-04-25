package oak

import (
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/joystick"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/mouse"
)

// InputType expresses some form of input to the engine to represent a player
type InputType = int32

// Supported Input Types
const (
	KeyboardMouse InputType = iota
	Joystick      InputType = iota
)

var (
	// MostRecentInput tracks what input type was most recently detected.
	// This is only updated if TrackInputChanges is true in the config at startup
	// TODO: scope this to controllers
	MostRecentInput InputType
)

func trackInputChanges() {
	event.GlobalBind(key.Down, func(event.CID, interface{}) int {
		atomic.SwapInt32(&MostRecentInput, KeyboardMouse)
		// TODO: Trigger that the most recent input changed?
		return 0
	})
	event.GlobalBind(mouse.Press, func(event.CID, interface{}) int {
		atomic.SwapInt32(&MostRecentInput, KeyboardMouse)
		return 0
	})
	event.GlobalBind("Tracking"+joystick.Change, func(event.CID, interface{}) int {
		atomic.SwapInt32(&MostRecentInput, Joystick)
		return 0
	})
}

type joyHandler struct{}

func (jh *joyHandler) Trigger(ev string, state interface{}) {
	event.Trigger("Tracking"+ev, state)
}

func trackJoystickChanges() {
	dlog.ErrorCheck(joystick.Init())
	go func() {
		jCh, _ := joystick.WaitForJoysticks(3 * time.Second)
		for j := range jCh {
			j.Handler = &joyHandler{}
			j.Listen(nil)
		}
	}()
}
