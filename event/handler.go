package event

import (
	"time"
)

var (
	_ Handler = &Bus{}
)

// Handler represents the necessary exported functions from an event.Bus
// for use in oak internally, and thus the functions that need to be replaced
// by alternative event handlers.
type Handler interface {
	EnterLoop(time.Duration)
	SetEnterLoopRate(time.Duration) error
	Stop() error
	Reset()
	TriggerForCaller(cid CallerID, event UnsafeEventID, data interface{}) chan struct{}
	Trigger(event UnsafeEventID, data interface{}) chan struct{}
	UnsafeBind(UnsafeEventID, CallerID, UnsafeBindable) Binding
	Unbind(Binding)
	SetCallerMap(*CallerMap)
	GetCallerMap() *CallerMap
}
