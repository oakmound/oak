package event_test

import (
	"sync/atomic"
	"testing"

	"github.com/oakmound/oak/v3/event"
)

func TestBus_UnsafeBind(t *testing.T) {
	t.Run("ConcurrentReset", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())

		var calls int32
		for i := 0; i < 1000; i++ {
			b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
				atomic.AddInt32(&calls, 1)
				return 0
			})
			b.Reset()
			// No matter what happens with thread scheduling above, this trigger should never increment calls
			<-b.Trigger(1, nil)
		}
		if calls != 0 {
			t.Fatal("a pre-reset binding was triggered after a bus reset")
		}
	})
}

func TestBus_Unbind(t *testing.T) {
	t.Run("ConcurrentReset", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())

		var goodCalls int32
		for i := 0; i < 1000; i++ {
			b1 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
				return 0
			})
			b.Unbind(b1)
			b.Reset()
			// b1 and b2 will share a bindID
			b2 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
				atomic.AddInt32(&goodCalls, 1)
				return 0
			})
			<-b2.Bound
			<-b.Trigger(1, nil)
			b2.Unbind()
		}
		if goodCalls != 1000 {
			t.Fatal("a pre-reset unbind unbound a post-reset binding", goodCalls)
		}
	})
}

func TestBind(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func TestGlobalBind(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func TestBus_UnbindAllFrom(t *testing.T) {
	t.Skip("TODO")
}
