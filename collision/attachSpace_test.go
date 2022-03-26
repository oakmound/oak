package collision

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
)

type aspace struct {
	AttachSpace
}

func (as *aspace) Init() event.CallerID {
	return event.NextID(as)
}

func TestAttachSpace(t *testing.T) {
	Clear()
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
	err := Attach(v, s, nil, 4, 4)
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

	err = Detach(s)
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
		t.Fatalf("unassinged space attach should have failed: %v", err)
	}
	err = Detach(s)
	if err == nil {
		t.Fatalf("unassinged space detach should have failed: %v", err)
	}
}
