package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
)

type cphase struct {
	CollisionPhase
}

func (cp *cphase) Init() event.CID {
	return event.NextID(cp)
}

func TestCollisionPhase(t *testing.T) {
	go event.ResolvePending()
	go func() {
		for {
			<-time.After(5 * time.Millisecond)
			<-event.TriggerBack("EnterFrame", nil)
		}
	}()
	cp := cphase{}
	cid := cp.Init()
	s := collision.NewSpace(10, 10, 10, 10, cid)
	assert.Nil(t, PhaseCollision(s))
	var active bool
	cid.Bind(func(int, interface{}) int {
		active = true
		return 0
	}, "MouseCollisionStart")
	cid.Bind(func(int, interface{}) int {
		active = false
		return 0
	}, "MouseCollisionStop")
	time.Sleep(200 * time.Millisecond)
	LastMouseEvent = Event{10, 10, "", ""}
	time.Sleep(200 * time.Millisecond)
	assert.True(t, active)
	LastMouseEvent = Event{21, 21, "", ""}
	time.Sleep(200 * time.Millisecond)
	assert.False(t, active)
	s = collision.NewSpace(10, 10, 10, 10, 5)
	assert.NotNil(t, PhaseCollision(s))
}
