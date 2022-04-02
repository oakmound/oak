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
	c1.MouseTree = collision.NewTree()
	ch := make(chan struct{})
	c1.eventHandler = event.NewBus(event.NewCallerMap())
	bnd := event.GlobalBind(c1.eventHandler, mouse.Click, func(_ *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("click binding never bound")
	case <-bnd.Bound:
	}
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	c1.MouseTree.Add(sp)
	c1.Propagate(mouse.PressOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.ReleaseOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger click binding")
	case <-ch:
	}
}

func TestMouseClicksRelative(t *testing.T) {
	c1 := NewWindow()
	c1.MouseTree = collision.NewTree()
	ch := make(chan struct{})
	c1.eventHandler = event.NewBus(event.NewCallerMap())
	bnd := event.GlobalBind(c1.eventHandler, mouse.RelativeClickOn, func(_ *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("click binding never bound")
	case <-bnd.Bound:
	}
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	c1.MouseTree.Add(sp)
	defer c1.MouseTree.Clear()
	c1.Propagate(mouse.RelativePressOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.PressOn))
	c1.Propagate(mouse.RelativeReleaseOn, mouse.NewEvent(5, 5, mouse.ButtonLeft, mouse.ReleaseOn))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger click binding")
	case <-ch:
	}
}

type ent struct {
	event.CallerID
}

func TestPropagate(t *testing.T) {
	c1 := NewWindow()
	c1.eventHandler = event.NewBus(event.NewCallerMap())

	thisEnt := ent{}
	thisEnt.CallerID = c1.eventHandler.GetCallerMap().Register(thisEnt)
	ch := make(chan struct{})
	s := collision.NewSpace(10, 10, 10, 10, thisEnt.CallerID)
	event.Bind(c1.eventHandler, mouse.PressOn, thisEnt, func(ent, *mouse.Event) event.Response {
		close(ch)
		return 0
	})
	c1.MouseTree = collision.NewTree()
	c1.MouseTree.Add(s)
	c1.Propagate(mouse.ReleaseOn, mouse.NewEvent(15, 15, mouse.ButtonLeft, mouse.Release))
	select {
	case <-ch:
		t.Fatalf("release propagation triggered press binding")
	case <-time.After(1 * time.Second):
	}
	c1.Propagate(mouse.PressOn, mouse.NewEvent(15, 15, mouse.ButtonLeft, mouse.Press))
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("propagation failed to trigger press binding")
	case <-ch:
	}
}
