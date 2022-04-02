package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/scene"
)

func TestTrackInputChanges(t *testing.T) {
	c1 := NewWindow()
	c1.SetLogicHandler(event.NewBus(event.NewCallerMap()))
	c1.AddScene("1", scene.Scene{})
	go c1.Init("1", func(c Config) (Config, error) {
		c.TrackInputChanges = true
		return c, nil
	})
	time.Sleep(2 * time.Second)
	expectedType := new(InputType)
	*expectedType = InputKeyboard
	failed := false
	event.GlobalBind(c1.eventHandler, InputChange, func(mode InputType) event.Response {
		if mode != *expectedType {
			failed = true
		}
		return 0
	})
	c1.TriggerKeyDown(key.Event{})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("keyboard change failed")
	}
	*expectedType = InputJoystick
	event.TriggerOn(c1.eventHandler, trackingJoystickChange, struct{}{})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("joystick change failed")
	}
	*expectedType = InputMouse
	c1.TriggerMouseEvent(mouse.Event{EventType: mouse.Press})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("mouse change failed")
	}
	*expectedType = InputKeyboard
	c1.mostRecentInput = int32(InputJoystick)
	c1.TriggerKeyDown(key.Event{})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("keyboard change failed")
	}
}
