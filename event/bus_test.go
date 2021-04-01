package event

import (
	"testing"
)

func TestBusStop(t *testing.T) {
	b := NewBus()
	phase := 0
	wait := make(chan struct{})
	go func() {
		if err := b.Stop(); err != nil {
			t.Fatalf("stop errored: %v", err)
		}
		if phase != 1 {
			t.Fatalf("expected phase %v, got %v", 1, phase)
		}
		wait <- struct{}{}
	}()
	b.updateCh <- struct{}{}
	<-b.doneCh
	phase = 1
	b.doneCh <- struct{}{}
	<-wait
}
