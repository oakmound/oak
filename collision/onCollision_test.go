package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
)

type cphase struct {
	Phase
	callers *event.CallerMap
}

func TestCollisionPhase(t *testing.T) {
	b := event.NewBus(event.NewCallerMap())
	go func() {
		for {
			<-time.After(5 * time.Millisecond)
			<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		}
	}()
	cp := &cphase{}
	cid := b.GetCallerMap().Register(cp)
	s := NewSpace(10, 10, 10, 10, cid)
	tree := NewTree()
	err := PhaseCollisionWithBus(s, tree, b)
	if err != nil {
		t.Fatalf("phase collision failed: %v", err)
	}
	activeCh := make(chan bool, 5)
	b1 := event.Bind(b, Start, cp, func(_ *cphase, _ Label) event.Response {
		activeCh <- true
		return 0
	})
	b2 := event.Bind(b, Stop, cp, func(_ *cphase, _ Label) event.Response {
		activeCh <- false
		return 0
	})
	<-b1.Bound
	<-b2.Bound
	s2 := NewLabeledSpace(15, 15, 10, 10, 5)
	tree.Add(s2)
	if active := <-activeCh; !active {
		t.Fatalf("collision should be active")
	}

	tree.Remove(s2)
	time.Sleep(200 * time.Millisecond)
	if active := <-activeCh; active {
		t.Fatalf("collision should be inactive")
	}
}

func TestPhaseCollision_Unembedded(t *testing.T) {
	t.Parallel()
	s3 := NewSpace(10, 10, 10, 10, 5)
	err := PhaseCollision(s3, nil)
	if err == nil {
		t.Fatalf("phase collision should have failed")
	}
}
