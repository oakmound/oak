package event_test

import (
	"sync/atomic"
	"testing"

	"github.com/oakmound/oak/v4/event"
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
			<-b2.Unbind()
		}
		if goodCalls != 1000 {
			t.Fatal("a pre-reset unbind unbound a post-reset binding", goodCalls)
		}
	})
}

func TestBind(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var cid event.CallerID
		id := b.GetCallerMap().Register(cid)
		var calls int32
		b1 := event.Bind(b, event.Enter, id, func(event.CallerID, event.EnterPayload) event.Response {
			atomic.AddInt32(&calls, 1)
			return 0
		})
		<-b1.Bound
		<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		if calls != 1 {
			t.Fatal(expectedError("calls", 1, calls))
		}
	})
}

func TestGlobalBind(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var calls int32
		b1 := event.GlobalBind(b, event.Enter, func(event.EnterPayload) event.Response {
			atomic.AddInt32(&calls, 1)
			return 0
		})
		<-b1.Bound
		<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		if calls != 1 {
			t.Fatal(expectedError("calls", 1, calls))
		}
	})
}

func TestBus_UnbindAllFrom(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var cid event.CallerID
		id := b.GetCallerMap().Register(cid)
		var calls int32
		for i := 0; i < 5; i++ {
			b1 := event.Bind(b, event.Enter, id, func(event.CallerID, event.EnterPayload) event.Response {
				atomic.AddInt32(&calls, 1)
				return 0
			})
			<-b1.Bound
		}
		id2 := b.GetCallerMap().Register(cid)
		b1 := event.Bind(b, event.Enter, id2, func(event.CallerID, event.EnterPayload) event.Response {
			atomic.AddInt32(&calls, 1)
			return 0
		})
		<-b1.Bound
		<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		if calls != 6 {
			t.Fatal(expectedError("calls", 1, calls))
		}
		<-b.UnbindAllFrom(id)
		<-event.TriggerOn(b, event.Enter, event.EnterPayload{})
		if calls != 7 {
			t.Fatal(expectedError("calls", 1, calls))
		}
	})
}
