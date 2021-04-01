package event

import (
	"testing"
	"time"
)

func TestResolvePendingWithRefreshRate(t *testing.T) {
	b := NewBus()
	b.SetRefreshRate(6 * time.Second)
	b.ResolvePending()
	failed := false
	b.Bind("EnterFrame", 0, func(CID, interface{}) int {
		failed = true
		return 0
	})
	ch := make(chan struct{}, 1000)
	b.UpdateLoop(60, ch)
	time.Sleep(3 * time.Second)
	if failed {
		t.Fatal("binding was called before refresh rate should have added binding")
	}
}
