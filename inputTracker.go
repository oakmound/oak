package oak

import (
	"sync"
	"time"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/joystick"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/mouse"
)

// InputType expresses some form of input to the engine to represent a player
// Todo v3: convert into int32 for use with atomic.SwapInt32
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
