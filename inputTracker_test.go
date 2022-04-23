package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/scene"
)

func TestTrackInputChanges(t *testing.T) {
	inputChangeFailed := make(chan bool)

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
	event.GlobalBind(c1.eventHandler, InputChange, func(mode InputType) event.Response {
		inputChangeFailed <- mode != *expectedType
		return 0
	})
	c1.TriggerKeyDown(key.Event{})
	if <-inputChangeFailed {
		t.Fatalf("keyboard change failed")
	}
	*expectedType = InputJoystick
	event.TriggerOn(c1.eventHandler, trackingJoystickChange, struct{}{})
	if <-inputChangeFailed {
		t.Fatalf("joystick change failed")
	}
	c1.mostRecentInput = int32(InputJoystick)
	*expectedType = InputMouse
	c1.TriggerMouseEvent(mouse.Event{EventType: mouse.Press})
	if <-inputChangeFailed {
		t.Fatalf("mouse change failed")
	}
	*expectedType = InputKeyboard
	c1.mostRecentInput = int32(InputJoystick)
	c1.TriggerKeyDown(key.Event{})
	if <-inputChangeFailed {
		t.Fatalf("keyboard change failed")
	}
	if c1.MostRecentInput() != InputKeyboard {
		t.Fatalf("most recent input getter failed")
	}
}
