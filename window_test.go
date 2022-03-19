package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
)

func TestMouseClicks(t *testing.T) {
	c1 := NewWindow()
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	var triggered bool
	go event.ResolveChanges()
	event.GlobalBind(mouse.Click, func(event.CallerID, interface{}) int {
		triggered = true
		return 0
	})
	time.Sleep(2 * time.Second)
	mouse.DefaultTree.Add(sp)
	c1.Propagate(mouse.PressOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.ReleaseOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatalf("propagation failed to trigger click binding")
	}
}

func TestMouseClicksRelative(t *testing.T) {
	c1 := NewWindow()
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	var triggered bool
	go c1.eventHandler.(*event.Bus).ResolveChanges()
	c1.eventHandler.GlobalBind(mouse.ClickOn+"Relative", func(event.CallerID, interface{}) int {
		triggered = true
		return 0
	})
	time.Sleep(2 * time.Second)
	c1.MouseTree.Add(sp)
	c1.Propagate(mouse.PressOn+"Relative", mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.ReleaseOn+"Relative", mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	time.Sleep(3 * time.Second)
	if !triggered {
		t.Fatalf("propagation failed to trigger click binding")
	}
}

type ent struct{}

func (e ent) Init() event.CallerID {
	return 0
}

func TestPropagate(t *testing.T) {
	c1 := NewWindow()
	go event.ResolveChanges()
	var triggered bool
	cid := event.CallerID(0).Parse(ent{})
	s := collision.NewSpace(10, 10, 10, 10, cid)
	s.CID.Bind("MouseDownOn", func(event.CallerID, interface{}) int {
		triggered = true
		return 0
	})
	mouse.Add(s)
	time.Sleep(200 * time.Millisecond)
	c1.Propagate("MouseUpOn", mouse.NewEvent(15, 15, mouse.ButtonLeft, "MouseUp"))
	time.Sleep(200 * time.Millisecond)
	if triggered {
		t.Fatalf("mouse up triggered binding")
	}
	time.Sleep(200 * time.Millisecond)
	c1.Propagate("MouseDownOn", mouse.NewEvent(15, 15, mouse.ButtonLeft, "MouseDown"))
	time.Sleep(200 * time.Millisecond)
	if !triggered {
		t.Fatalf("mouse down failed to trigger binding")
	}
}
