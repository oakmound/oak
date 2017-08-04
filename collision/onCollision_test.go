package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
)

type cphase struct {
	Phase
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
	s := NewSpace(10, 10, 10, 10, cid)
	assert.Nil(t, PhaseCollision(s))
	var active bool
	cid.Bind(func(int, interface{}) int {
		active = true
		return 0
	}, "CollisionStart")
	cid.Bind(func(int, interface{}) int {
		active = false
		return 0
	}, "CollisionStop")

	s = NewLabeledSpace(15, 15, 10, 10, 5)
	Add(s)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, active)

	Remove(s)
	time.Sleep(10 * time.Millisecond)
	assert.False(t, active)

	s = NewSpace(10, 10, 10, 10, 5)
	assert.NotNil(t, PhaseCollision(s))
}
