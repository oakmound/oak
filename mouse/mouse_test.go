package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mobile/event/mouse"
)

func TestMouseClicks(t *testing.T) {
	sp := collision.NewFullSpace(0, 0, 100, 100, 1, 0)
	var triggered bool
	go event.ResolvePending()
	event.GlobalBind(func(int, interface{}) int {
		triggered = true
		return 0
	}, Click)
	DefTree.Add(sp)
	Propagate(PressOn, NewEvent(5, 5, "LeftMouse", PressOn))
	Propagate(ReleaseOn, NewEvent(5, 5, "LeftMouse", ReleaseOn))
	time.Sleep(2 * time.Second)
	assert.True(t, triggered)
}

func TestButtonIdentity(t *testing.T) {
	// This is a pretty worthless test
	assert.Equal(t, GetMouseButton(mouse.ButtonLeft), "LeftMouse")
	assert.Equal(t, GetMouseButton(mouse.ButtonRight), "RightMouse")
	assert.Equal(t, GetMouseButton(mouse.ButtonMiddle), "MiddleMouse")
	assert.Equal(t, GetMouseButton(mouse.ButtonWheelUp), "ScrollUpMouse")
	assert.Equal(t, GetMouseButton(mouse.ButtonWheelDown), "ScrollDownMouse")
	assert.Equal(t, GetMouseButton(mouse.ButtonWheelLeft), "")
	assert.Equal(t, GetMouseButton(mouse.ButtonWheelRight), "")
	assert.Equal(t, GetMouseButton(mouse.ButtonNone), "")
}

func TestEventNameIdentity(t *testing.T) {
	assert.Equal(t, GetEventName(mouse.DirPress, 0), "MousePress")
	assert.Equal(t, GetEventName(mouse.DirRelease, 0), "MouseRelease")
	assert.Equal(t, GetEventName(mouse.DirNone, -2), "MouseScrollDown")
	assert.Equal(t, GetEventName(mouse.DirNone, -1), "MouseScrollUp")
	assert.Equal(t, GetEventName(mouse.DirNone, 0), "MouseDrag")
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
	s.CID.Bind(func(int, interface{}) int {
		triggered = true
		return 0
	}, "MouseDownOn")
	Add(s)
	time.Sleep(200 * time.Millisecond)
	Propagate("MouseUpOn", NewEvent(15, 15, "LeftMouse", "MouseUp"))
	time.Sleep(200 * time.Millisecond)
	assert.False(t, triggered)
	time.Sleep(200 * time.Millisecond)
	Propagate("MouseDownOn", NewEvent(15, 15, "LeftMouse", "MouseDown"))
	time.Sleep(200 * time.Millisecond)
	assert.True(t, triggered)
}
