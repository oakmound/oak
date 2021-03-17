package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"

	"golang.org/x/mobile/event/mouse"
)

func TestMouseClicks(t *testing.T) {
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	var triggered bool
	go event.ResolvePending()
	event.GlobalBind(Click, func(event.CID, interface{}) int {
		triggered = true
		return 0
	})
	DefTree.Add(sp)
	Propagate(PressOn, NewEvent(5, 5, ButtonLeft, PressOn))
	Propagate(ReleaseOn, NewEvent(5, 5, ButtonLeft, ReleaseOn))
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatalf("propagation failed to trigger click binding")
	}
}

func TestMouseClicksRelative(t *testing.T) {
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	var triggered bool
	go event.ResolvePending()
	event.GlobalBind(ClickOn+"Relative", func(event.CID, interface{}) int {
		triggered = true
		return 0
	})
	DefTree.Add(sp)
	Propagate(PressOn+"Relative", NewEvent(5, 5, ButtonLeft, PressOn))
	Propagate(ReleaseOn+"Relative", NewEvent(5, 5, ButtonLeft, ReleaseOn))
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatalf("propagation failed to trigger click binding")
	}
}

func TestEventNameIdentity(t *testing.T) {
	if GetEventName(mouse.DirPress, 0) != "MousePress" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirPress, "MousePress")
	}
	if GetEventName(mouse.DirRelease, 0) != "MouseRelease" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirRelease, "MouseRelease")
	}
	if GetEventName(mouse.DirNone, -2) != "MouseScrollDown" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollDown")
	}
	if GetEventName(mouse.DirNone, -1) != "MouseScrollUp" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollUp")
	}
	if GetEventName(mouse.DirNone, 0) != "MouseDrag" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseDrag")
	}
}

type ent struct{}

func (e ent) Init() event.CID {
	return event.NextID(e)
}
func TestPropagate(t *testing.T) {
	go event.ResolvePending()
	var triggered bool
	cid := event.CID(0).Parse(ent{})
	s := collision.NewSpace(10, 10, 10, 10, cid)
	s.CID.Bind("MouseDownOn", func(event.CID, interface{}) int {
		triggered = true
		return 0
	})
	Add(s)
	time.Sleep(200 * time.Millisecond)
	Propagate("MouseUpOn", NewEvent(15, 15, ButtonLeft, "MouseUp"))
	time.Sleep(200 * time.Millisecond)
	if triggered {
		t.Fatalf("mouse up triggered binding")
	}
	time.Sleep(200 * time.Millisecond)
	Propagate("MouseDownOn", NewEvent(15, 15, ButtonLeft, "MouseDown"))
	time.Sleep(200 * time.Millisecond)
	if !triggered {
		t.Fatalf("mouse down failed to trigger binding")
	}
}
