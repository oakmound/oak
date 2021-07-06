package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	okey "github.com/oakmound/oak/v3/key"
	"golang.org/x/mobile/event/key"
)

func TestInputLoop(t *testing.T) {
	c1 := blankScene(t)
	c1.SetLogicHandler(event.NewBus(nil))
	c1.windowControl.Send(okey.Event{
		Direction: key.DirPress,
		Code:      key.Code0,
	})
	c1.windowControl.Send(okey.Event{
		Direction: key.DirNone,
		Code:      key.Code0,
	})
	c1.windowControl.Send(okey.Event{
		Direction: key.DirRelease,
		Code:      key.Code0,
	})
	time.Sleep(2 * time.Second)
}
