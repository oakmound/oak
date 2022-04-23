package event_test

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/oakmound/oak/v4/event"
)

func TestNewBus(t *testing.T) {
	t.Run("DefaultCallerMap", func(t *testing.T) {
		b := event.NewBus(nil)
		if b.GetCallerMap() != event.DefaultCallerMap {
			t.Fatal("nil caller map not turned into default caller map")
		}
	})
	t.Run("Basic", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		if b == nil {
			t.Fatal("NewBus created nil bus")
		}
	})
}

func TestBus_SetCallerMap(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		cm1 := event.NewCallerMap()
		b := event.NewBus(cm1)
		c1 := event.CallerID(rand.Intn(10000))
		b.GetCallerMap().Register(c1)
		cm2 := event.NewCallerMap()
		b.SetCallerMap(cm2)
		if b.GetCallerMap().HasEntity(c1) {
			t.Fatal("event had old entity after changed caller map")
		}
	})
}

func TestBus_ClearPersistentBindings(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var impersistentCalls int32
		var persistentCalls int32
		b1 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&impersistentCalls, 1)
			return 0
		})
		b2 := b.PersistentBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&persistentCalls, 1)
			return 0
		})
		<-b1.Bound
		<-b2.Bound
		<-b.Trigger(1, nil)
		if impersistentCalls != 1 {
			t.Fatal(expectedError("impersistent calls", 1, impersistentCalls))
		}
		if persistentCalls != 1 {
			t.Fatal(expectedError("persistent calls", 1, persistentCalls))
		}
		b.Reset()
		<-b.Trigger(1, nil)
		if impersistentCalls != 1 {
			t.Fatal(expectedError("impersistent calls", 1, impersistentCalls))
		}
		if persistentCalls != 2 {
			t.Fatal(expectedError("persistent calls", 2, persistentCalls))
		}
		b.ClearPersistentBindings()
		b.Reset()
		<-b.Trigger(1, nil)
		if impersistentCalls != 1 {
			t.Fatal(expectedError("impersistent calls", 1, impersistentCalls))
		}
		if persistentCalls != 2 {
			t.Fatal(expectedError("persistent calls", 2, persistentCalls))
		}
	})
}

func TestBus_EnterLoop(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var calls int32
		b1 := b.UnsafeBind(event.Enter.UnsafeEventID, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&calls, 1)
			return 0
		})
		<-b1.Bound
		cancel := event.EnterLoop(b, 50*time.Millisecond)
		time.Sleep(1*time.Second + 15*time.Millisecond)
		cancel()
		if calls != 20 {
			t.Fatal(expectedError("calls", 20, calls))
		}
	})
}
