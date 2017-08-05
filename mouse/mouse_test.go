package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mobile/event/mouse"
)

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
	Propagate("MouseUpOn", Event{15, 15, "LeftMouse", "MouseUp"})
	time.Sleep(200 * time.Millisecond)
	assert.False(t, triggered)
	time.Sleep(200 * time.Millisecond)
	Propagate("MouseDownOn", Event{15, 15, "LeftMouse", "MouseDown"})
	time.Sleep(200 * time.Millisecond)
	assert.True(t, triggered)
}
