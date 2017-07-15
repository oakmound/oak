package timing

import (
	"testing"
	"time"
)

func TestDoAfter(t *testing.T) {
	go DoAfter(3*time.Second, func() {
		t.Fatal("DoAfter did not stop")
	})
	// Wait to make sure the routine started
	time.Sleep(1 * time.Second)
outer:
	for {
		select {
		case ClearDelayCh <- true:
		default:
			break outer
		}
	}
	time.Sleep(3 * time.Second)
}
