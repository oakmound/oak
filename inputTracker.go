package oak

import (
	"sync"
	"time"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
)

// InputType expresses some form of input to the engine to represent a player
type InputType int

// Supported Input Types
const (
	KeyboardMouse InputType = iota
	Joystick      InputType = iota
)

var (
	// MostRecentInput tracks what input type was most recently detected.
	// This is only updated if TrackInputChanges is true in the config at startup
	MostRecentInput InputType
	recentInputLock sync.Mutex
)

func trackInputChanges() {
	event.GlobalBind(func(int, interface{}) int {
		recentInputLock.Lock()
		MostRecentInput = KeyboardMouse
		// Trigger that the most recent input changed?
		recentInputLock.Unlock()
		return 0
	}, key.Down)
	event.GlobalBind(func(int, interface{}) int {
		recentInputLock.Lock()
		MostRecentInput = KeyboardMouse
		recentInputLock.Unlock()
		return 0
	}, mouse.Press)
	event.GlobalBind(func(int, interface{}) int {
		recentInputLock.Lock()
		MostRecentInput = Joystick
		recentInputLock.Unlock()
		return 0
	}, "Tracking"+joystick.Change)
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
