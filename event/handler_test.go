package event

import (
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	updateCh := make(chan struct{})
	if UpdateLoop(60, updateCh) != nil {
		t.Fatalf("UpdateLoop failed")
	}
	triggers := 0
	Bind(Enter, 0, func(CID, interface{}) int {
		triggers++
		return 0
	})
	sleep()
	if triggers != 1 {
		t.Fatalf("expected update loop to increment triggers")
	}
	<-updateCh
	sleep()
	if triggers != 2 {
		t.Fatalf("expected update loop to increment triggers")
	}
	if FramesElapsed() != 2 {
		t.Fatalf("expected 2 update frames to have elapsed")
	}
	if SetTick(1) != nil {
		t.Fatalf("SetTick failed")
	}
	<-updateCh
	if Stop() != nil {
		t.Fatalf("Stop failed")
	}
	sleep()
	sleep()
	select {
	case <-updateCh:
		t.Fatal("Handler should be closed")
	default:
	}
	if Update() != nil {
		t.Fatalf("Update failed")
	}
	sleep()

	if triggers != 4 {
		t.Fatalf("expected update to increment triggers")
	}
	if Flush() != nil {
		t.Fatalf("Flush failed")
	}

	Flush()
	sleep()
	if Update() != nil {
		t.Fatalf("final Update failed")
	}
	sleep()
	sleep()
	Reset()
}

func BenchmarkHandler(b *testing.B) {
	triggers := 0
	entities := 10
	go DefaultBus.ResolvePending()
	for i := 0; i < entities; i++ {
		DefaultBus.GlobalBind(Enter, func(CID, interface{}) int {
			triggers++
			return 0
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-DefaultBus.TriggerBack(Enter, DefaultBus.framesElapsed)
	}
}

func TestPauseAndResume(t *testing.T) {
	b := NewBus()
	b.ResolvePending()
	triggerCt := 0
	b.Bind("EnterFrame", 0, func(CID, interface{}) int {
		triggerCt++
		return 0
	})
	ch := make(chan struct{}, 1000)
	b.UpdateLoop(60, ch)
	time.Sleep(1 * time.Second)
	b.Pause()
	time.Sleep(1 * time.Second)
	oldCt := triggerCt
	time.Sleep(1 * time.Second)
	if oldCt != triggerCt {
		t.Fatalf("pause did not stop enter frame from triggering: expected %v got %v", oldCt, triggerCt)
	}

	b.Resume()
	time.Sleep(1 * time.Second)
	newCt := triggerCt
	if newCt == oldCt {
		t.Fatalf("resume did not resume enter frame triggering: expected %v got %v", oldCt, newCt)
	}
}
