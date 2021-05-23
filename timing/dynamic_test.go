package timing

import (
	"testing"
	"time"
)

func TestDynamicTickerFns(t *testing.T) {
	t.Parallel()
	dt := NewDynamicTicker()
	time.Sleep(10 * time.Second)
	select {
	case <-dt.C:
		t.Fatal("Dynamic Ticker should not initially send")
	default:
	}
	dt.Step()
	nextTime := <-dt.C
	// The above just needs to not time out
	now := time.Now()
	dt.SetTick(1 * time.Second)
	nextTime = <-dt.C
	got := nextTime.Sub(now)
	expectedLessThan := 1100 * time.Millisecond
	if got >= expectedLessThan {
		t.Fatalf("expected less than %v, got %v", expectedLessThan, got)
	}
	dt.Stop()
	select {
	case _, ok := <-dt.C:
		if ok {
			t.Fatal("Dynamic Ticker failed to stop")
		}
	default:
	}

	dt = NewDynamicTicker()
	dt.SetTick(1 * time.Second)
	time.Sleep(2 * time.Second)
	dt.Step()
	dt.SetTick(2 * time.Second)
	dt.ForceStep()
}

func TestDynamicTickerStop(t *testing.T) {
	t.Parallel()

	dt := NewDynamicTicker()
	dt.Stop()

	// Successive stops
	dt = NewDynamicTicker()
	dt.Step()
	dt.Stop()
	dt.Stop()

	// Unconsumed tick -> stop
	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	dt.SetTick(1 * time.Millisecond)
	time.Sleep(2 * time.Second)
	dt.Stop()

	// Unconsumed step -> stop
	dt = NewDynamicTicker()
	dt.SetTick(1 * time.Millisecond)
	time.Sleep(1 * time.Second)
	dt.SetTick(2 * time.Millisecond)
	dt.Step()
	dt.Stop()

	// Successive steps
	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	for i := 0; i < 20; i++ {
		dt.Step()
		time.Sleep(30 * time.Millisecond)
	}
	dt.Stop()

	// Unconsumed step -> Set Tick
	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	dt.Step()
	time.Sleep(1 * time.Second)
	dt.SetTick(1 * time.Millisecond)
}
