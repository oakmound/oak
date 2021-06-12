package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/joystick"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/scene"
)

func TestTrackInputChanges(t *testing.T) {
	c1 := NewController()
	c1.SetLogicHandler(event.NewBus())
	c1.AddScene("1", scene.Scene{})
	go c1.Init("1", func(c Config) (Config, error) {
		c.TrackInputChanges = true
		return c, nil
	})
	time.Sleep(2 * time.Second)
	expectedType := new(InputType)
	*expectedType = KeyboardMouse
	failed := false
	c1.logicHandler.GlobalBind(event.InputChange, func(_ event.CID, payload interface{}) int {
		mode := payload.(InputType)
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
	*expectedType = Joystick
	c1.logicHandler.Trigger("Tracking"+joystick.Change, &joystick.State{})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("joystick change failed")
	}
	*expectedType = KeyboardMouse
	c1.TriggerMouseEvent(mouse.Event{Event: mouse.Press})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("mouse change failed")
	}
	c1.mostRecentInput = Joystick
	c1.TriggerKeyDown(key.Event{})
	time.Sleep(2 * time.Second)
	if failed {
		t.Fatalf("keyboard change failed")
	}
}
