package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/event"
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
			<-event.TriggerBack(event.Enter, nil)
		}
	}()
	cp := cphase{}
	cid := cp.Init()
	s := NewSpace(10, 10, 10, 10, cid)
	err := PhaseCollision(s, nil)
	if err != nil {
		t.Fatalf("phase collision failed: %v", err)
	}
	var active bool
	cid.Bind("CollisionStart", func(event.CID, interface{}) int {
		active = true
		return 0
	})
	cid.Bind("CollisionStop", func(event.CID, interface{}) int {
		active = false
		return 0
	})

	s2 := NewLabeledSpace(15, 15, 10, 10, 5)
	Add(s2)
	time.Sleep(200 * time.Millisecond)
	if !active {
		t.Fatalf("collision should be active")
	}

	Remove(s2)
	time.Sleep(200 * time.Millisecond)
	if active {
		t.Fatalf("collision should be inactive")
	}

	s3 := NewSpace(10, 10, 10, 10, 5)
	err = PhaseCollision(s3, nil)
	if err == nil {
		t.Fatalf("phase collision should have failed")
	}

	err = PhaseCollision(s, DefaultTree)
	if err != nil {
		t.Fatalf("phase collision failed: %v", err)
	}
}
