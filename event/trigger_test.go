package event_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/event"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

func TestBus_TriggerForCaller(t *testing.T) {
	t.Run("NoBinding", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		ch := b.TriggerForCaller(0, id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
	})
	t.Run("GlobalWithBinding", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		errs := make(chan error)
		binding := b.UnsafeBind(id, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			defer close(errs)
			if ci != 0 {
				errs <- expectedError("callerID", 0, ci)
			}
			if h != b {
				errs <- expectedError("bus", b, h)
			}
			if i != nil {
				errs <- expectedError("payload", nil, i)
			}
			return 0
		})
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.TriggerForCaller(0, id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		for err := range errs {
			t.Error(err)
		}
	})
	t.Run("WithMissingCallerID", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		callerID := event.CallerID(rand.Intn(100000))
		errs := make(chan error)
		binding := b.UnsafeBind(id, callerID, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			errs <- fmt.Errorf("binding should not be triggered")
			return 0
		})
		_ = binding
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.TriggerForCaller(callerID, id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		select {
		case err := <-errs:
			t.Error(err)
		default:
		}
	})
	t.Run("WithValidCallerID", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var cid event.CallerID
		callerID := b.GetCallerMap().Register(cid)
		id := event.UnsafeEventID(rand.Intn(100000))
		errs := make(chan error)
		binding := b.UnsafeBind(id, callerID, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			defer close(errs)
			if ci != callerID {
				errs <- expectedError("callerID", callerID, ci)
			}
			if h != b {
				errs <- expectedError("bus", b, h)
			}
			if i != nil {
				errs <- expectedError("payload", nil, i)
			}
			return 0
		})
		_ = binding
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.TriggerForCaller(callerID, id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		for err := range errs {
			t.Error(err)
		}
	})
}

func TestBus_Trigger(t *testing.T) {
	t.Run("NoBinding", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		ch := b.Trigger(id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
	})
	t.Run("GlobalWithBinding", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		errs := make(chan error)
		binding := b.UnsafeBind(id, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			defer close(errs)
			if ci != 0 {
				errs <- expectedError("callerID", 0, ci)
			}
			if h != b {
				errs <- expectedError("bus", b, h)
			}
			if i != nil {
				errs <- expectedError("payload", nil, i)
			}
			return 0
		})
		_ = binding
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.Trigger(id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		for err := range errs {
			t.Error(err)
		}
	})
	t.Run("WithMissingCallerID", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		id := event.UnsafeEventID(rand.Intn(100000))
		callerID := rand.Intn(100000)
		errs := make(chan error)
		binding := b.UnsafeBind(id, event.CallerID(callerID), func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			errs <- fmt.Errorf("binding should not be triggered")
			return 0
		})
		_ = binding
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.Trigger(id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		select {
		case err := <-errs:
			t.Error(err)
		default:
		}
	})
	t.Run("WithValidCallerID", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		var cid event.CallerID
		callerID := b.GetCallerMap().Register(cid)
		id := event.UnsafeEventID(rand.Intn(100000))
		errs := make(chan error)
		binding := b.UnsafeBind(id, event.CallerID(callerID), func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			defer close(errs)
			if ci != callerID {
				errs <- expectedError("callerID", callerID, ci)
			}
			if h != b {
				errs <- expectedError("bus", b, h)
			}
			if i != nil {
				errs <- expectedError("payload", nil, i)
			}
			return 0
		})
		_ = binding
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for bind to close channel")
		case <-binding.Bound:
		}
		ch := b.Trigger(id, nil)
		select {
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for trigger to close channel")
		case <-ch:
		}
		for err := range errs {
			t.Error(err)
		}
	})
}

// TriggerOn and TriggerForCallerOn are simple wrappers of the tested methods above, so
// they are not tested thoroughly.

func TestTriggerOn(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		eventID := event.RegisterEvent[struct{}]()
		event.TriggerOn(b, eventID, struct{}{})
	})
}

func TestTriggerForCallerOn(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		eventID := event.RegisterEvent[struct{}]()
		event.TriggerForCallerOn(b, 0, eventID, struct{}{})
	})
}

func expectedError(name string, expected, got interface{}) error {
	return fmt.Errorf("expected %s to be %v, got %v", name, expected, got)
}
