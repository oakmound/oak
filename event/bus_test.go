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
	b.updateCh <- struct{}{}
	<-b.doneCh
	phase = 1
	b.doneCh <- struct{}{}
	<-wait
	if topErr != nil {
		t.Fatal(topErr)
	}
}
