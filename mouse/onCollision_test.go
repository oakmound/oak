package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
)

type cphase struct {
	CollisionPhase
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
	s := collision.NewSpace(10, 10, 10, 10, cid)
	err := PhaseCollision(s, b)
	if err != nil {
		t.Fatalf("phase collision failed: %v", err)
	}
	activeCh := make(chan bool, 5)
	b1 := event.Bind(b, Start, cp, func(_ *cphase, _ *Event) event.Response {
		activeCh <- true
		return 0
	})
	b2 := event.Bind(b, Stop, cp, func(_ *cphase, _ *Event) event.Response {
		activeCh <- false
		return 0
	})
	<-b1.Bound
	<-b2.Bound
	LastEvent = Event{
		Point2: floatgeom.Point2{10, 10},
	}
	if active := <-activeCh; !active {
		t.Fatalf("collision should be active")
	}

	LastEvent = Event{
		Point2: floatgeom.Point2{21, 21},
	}
	time.Sleep(200 * time.Millisecond)
	if active := <-activeCh; active {
		t.Fatalf("collision should be inactive")
	}
}

func TestPhaseCollision_Unembedded(t *testing.T) {
	t.Parallel()
	s3 := collision.NewSpace(10, 10, 10, 10, 5)
	err := PhaseCollision(s3, event.DefaultBus)
	if err == nil {
		t.Fatalf("phase collision should have failed")
	}
}
