package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v4/event"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

func TestInputLoop(t *testing.T) {
	c1 := blankScene(t)
	c1.SetLogicHandler(event.NewBus(nil))
	c1.Window.Send(key.Event{
		Direction: key.DirPress,
		Code:      key.Code0,
	})
	c1.Window.Send(key.Event{
		Direction: key.DirNone,
		Code:      key.Code0,
	})
	c1.Window.Send(key.Event{
		Direction: key.DirRelease,
		Code:      key.Code0,
	})
	c1.Window.Send(mouse.Event{})
	time.Sleep(2 * time.Second)
}
