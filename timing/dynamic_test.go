package timing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDynamicTicker(t *testing.T) {
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
	case <-dt.C:
		t.Fatal("Dynamic Ticker failed to stop")
	default:
	}
}
