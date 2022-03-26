package collision

import (
	"fmt"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
)

type aspace struct {
	AttachSpace
}

func TestAttachSpace(t *testing.T) {
	Clear()
	b := event.NewBus(event.NewCallerMap())
	go func() {
		for {
			<-time.After(5 * time.Millisecond)
			<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		}
	}()
	as := &aspace{}
	cid := b.GetCallerMap().Register(as)
	v := physics.NewVector(0, 0)
	s := NewSpace(100, 100, 10, 10, cid)
	Add(s)
	fmt.Println(s.CID)
	err := AttachWithBus(v, s, nil, b, 4, 4)
	if err != nil {
		t.Fatalf("attach failed: %v", err)
	}
	v.SetPos(5, 5)
	time.Sleep(200 * time.Millisecond)
	if s.X() != 9 {
		t.Fatalf("expected attached space to have x of 9, was %v", s.X())
	}
	if s.Y() != 9 {
		t.Fatalf("expected attached space to have y of 9, was %v", s.Y())
	}

	err = DetachWithBus(s, b)
	if err != nil {
		t.Fatalf("detach failed: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	v.SetPos(4, 4)
	time.Sleep(200 * time.Millisecond)
	if s.X() != 9 {
		t.Fatalf("expected attached space to have x of 9, was %v", s.X())
	}
	if s.Y() != 9 {
		t.Fatalf("expected attached space to have y of 9, was %v", s.Y())
	}

	// Failures
	s = NewUnassignedSpace(0, 0, 1, 1)
	err = Attach(v, s, nil)
	if err == nil {
		t.Fatalf("unassigned space attach should have failed: %v", err)
	}
	err = Detach(s)
	if err == nil {
		t.Fatalf("unassigned space detach should have failed: %v", err)
	}
}
