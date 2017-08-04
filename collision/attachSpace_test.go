package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"
)

type aspace struct {
	AttachSpace
}

func (as *aspace) Init() event.CID {
	return event.NextID(as)
}

func TestAttachSpace(t *testing.T) {
	Clear()
	go event.ResolvePending()
	go func() {
		for {
			<-time.After(5 * time.Millisecond)
			<-event.TriggerBack("EnterFrame", nil)
		}
	}()
	as := aspace{}
	v := physics.NewVector(0, 0)
	s := NewSpace(100, 100, 10, 10, as.Init())
	assert.Nil(t, Attach(v, s, 4, 4))
	v.SetPos(5, 5)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, s.GetX(), 9.0)
	assert.Equal(t, s.GetY(), 9.0)

	assert.Nil(t, Detach(s))
	v.SetPos(4, 4)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, s.GetX(), 9.0)
	assert.Equal(t, s.GetY(), 9.0)

	// Failures
	s = NewUnassignedSpace(0, 0, 1, 1)
	assert.NotNil(t, Attach(v, s))
	assert.NotNil(t, Detach(s))
}
