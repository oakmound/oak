package event

import (
	"fmt"
	"testing"
	"time"
)

func TestBusStop(t *testing.T) {
	b := NewBus(nil)
	b.Ticker = time.NewTicker(10000 * time.Second)
	phase := 0
	wait := make(chan struct{})
	var topErr error
	go func() {
		if err := b.Stop(); err != nil {
			topErr = fmt.Errorf("stop errored: %v", err)
		}
		if phase != 1 {
			topErr = fmt.Errorf("expected phase %v, got %v", 1, phase)
		}
		wait <- struct{}{}
	}()
	phase = 1

	<-b.doneCh
	<-wait
	if topErr != nil {
		t.Fatal(topErr)
	}
}

func TestBusPersistentBind(t *testing.T) {
	t.Parallel()
	b := NewBus(nil)
	ev := "eventName"
	calls := 0
	b.PersistentBind(ev, 0, func(c CID, i interface{}) int {
		calls++
		return 0
	})
	b.Flush()
	<-b.TriggerBack(ev, nil)
	if calls != 1 {
		t.Fatalf("expected binding to be called once, was called %d time(s)", calls)
	}
	b.Reset()
	<-b.TriggerBack(ev, nil)
	if calls != 2 {
		t.Fatalf("expected binding to be called twice, was called %d time(s)", calls)
	}
}
