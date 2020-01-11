package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/physics"
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
			<-event.TriggerBack(event.Enter, nil)
		}
	}()
	as := aspace{}
	v := physics.NewVector(0, 0)
	s := NewSpace(100, 100, 10, 10, as.Init())
	Add(s)
	assert.Nil(t, Attach(v, s, 4, 4))
	v.SetPos(5, 5)
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, s.X(), 9.0)
	assert.Equal(t, s.Y(), 9.0)

	assert.Nil(t, Detach(s))
	v.SetPos(4, 4)
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, s.X(), 9.0)
	assert.Equal(t, s.Y(), 9.0)

	// Failures
	s = NewUnassignedSpace(0, 0, 1, 1)
	assert.NotNil(t, Attach(v, s))
	assert.NotNil(t, Detach(s))
}
