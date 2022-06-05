package event_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/oakmound/oak/v4/event"
)

func TestBindingResponses(t *testing.T) {
	t.Run("UnbindThisBinding", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())
		li := 1
		var calls int32
		var innerCalls int32
		var b2 event.Binding
		b1 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&calls, 1)
			b2 = b.UnsafeBind(2, ci.CID(), func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
				atomic.AddInt32(&innerCalls, 1)
				return event.ResponseUnbindThisBinding
			})
			li++
			return event.ResponseUnbindThisBinding
		})
		<-b1.Bound
		<-b.Trigger(1, nil)
		if calls != 1 {
			t.Fatal(expectedError("calls", 1, calls))
		}
		// we do not get a signal for when this unbinding is finished
		time.Sleep(1 * time.Second)
		<-b.Trigger(1, nil)
		if calls != 1 {
			t.Fatal(expectedError("calls", 1, calls))
		}
		<-b2.Bound
		<-b.Trigger(2, nil)
		time.Sleep(1 * time.Second)
		<-b.Trigger(2, nil)
		t.Fatal(fmt.Errorf("%d", innerCalls))
		t.Fatal(fmt.Errorf("%d", li))
	})
	t.Run("UNbindThisCaller", func(t *testing.T) {
		b := event.NewBus(event.NewCallerMap())

		var calls int32
		b1 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&calls, 1)
			return event.ResponseUnbindThisCaller
		})
		<-b1.Bound
		b2 := b.UnsafeBind(1, 0, func(ci event.CallerID, h event.Handler, i interface{}) event.Response {
			atomic.AddInt32(&calls, 1)
			return 0
		})
		<-b2.Bound
		<-b.Trigger(1, nil)
		if calls != 2 {
			t.Fatal(expectedError("calls", 1, calls))
		}
		// we do not get a signal for when this unbinding is finished
		time.Sleep(1 * time.Second)
		<-b.Trigger(1, nil)
		if calls != 2 {
			t.Fatal(expectedError("calls", 2, calls))
		}
	})
}
