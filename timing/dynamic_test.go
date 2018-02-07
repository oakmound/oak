package timing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDynamicTickerFns(t *testing.T) {
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
	assert.True(t, nextTime.Sub(now) < 1100*time.Millisecond)
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
	dt := NewDynamicTicker()
	dt.Stop()

	dt = NewDynamicTicker()
	dt.Step()
	dt.Stop()
	dt.Stop()

	dt = NewDynamicTicker()
	go func() {
		<-dt.C
	}()
	dt.Stop()

	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	dt.SetTick(1 * time.Millisecond)
	time.Sleep(2 * time.Second)
	dt.Stop()

	dt = NewDynamicTicker()
	dt.SetTick(1 * time.Millisecond)
	time.Sleep(1 * time.Second)
	dt.SetTick(2 * time.Millisecond)
	dt.Step()
	dt.Stop()

	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	for i := 0; i < 20; i++ {
		dt.Step()
		time.Sleep(30 * time.Millisecond)
	}
	dt.Stop()

	dt = NewDynamicTicker()
	time.Sleep(1 * time.Second)
	dt.Step()
	time.Sleep(1 * time.Second)
	dt.SetTick(1 * time.Millisecond)
}
