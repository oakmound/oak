package event

import (
	"sync"
	"time"
)

// A Bus stores bindables to be triggered by events.
type Bus struct {
	// nextBindID is an atomically incrementing value to track bindings within this structure
	nextBindID *int64

	// resetCount increments every time the bus is reset. bindings and unbindings make sure that
	// they are called on a bus with an unchanged reset count, and become NOPs if performed on
	// a bus with a different reset count to ensure they do not interfere with a bus using different
	// bind IDs.
	resetCount         int64
	bindingMap         map[UnsafeEventID]map[CallerID]bindableList
	persistentBindings []persistentBinding

	callerMap *CallerMap

	mutex sync.RWMutex
}

// a persistentBinding is rebound every time the bus is reset.
type persistentBinding struct {
	eventID  UnsafeEventID
	callerID CallerID
	fn       UnsafeBindable
}

// NewBus returns an empty event bus with an assigned caller map. If nil
// is provided, the caller map used will be DefaultCallerMap
func NewBus(callerMap *CallerMap) *Bus {
	if callerMap == nil {
		callerMap = DefaultCallerMap
	}
	return &Bus{
		nextBindID: new(int64),
		bindingMap: make(map[UnsafeEventID]map[CallerID]bindableList),
		callerMap:  callerMap,
	}
}

// SetCallerMap updates a bus to use a specific set of callers.
func (bus *Bus) SetCallerMap(cm *CallerMap) {
	bus.callerMap = cm
}

// GetCallerMap returns this bus's caller map.
func (b *Bus) GetCallerMap() *CallerMap {
	return b.callerMap
}

// ClearPersistentBindings removes all persistent bindings. It will not unbind them
// from the bus, but they will not be bound following the next bus reset.
func (eb *Bus) ClearPersistentBindings() {
	eb.mutex.Lock()
	eb.persistentBindings = eb.persistentBindings[:0]
	eb.mutex.Unlock()
}

// Reset unbinds all present, non-persistent bindings on the bus. It will block until
// persistent bindings are in place.
func (bus *Bus) Reset() {
	bus.mutex.Lock()
	bus.resetCount++
	bus.bindingMap = make(map[UnsafeEventID]map[CallerID]bindableList)
	repersist := make([]Binding, len(bus.persistentBindings))
	for i, pb := range bus.persistentBindings {
		repersist[i] = bus.UnsafeBind(pb.eventID, pb.callerID, pb.fn)
	}
	bus.mutex.Unlock()
	for _, bnd := range repersist {
		<-bnd.Bound
	}
}

// EnterLoop triggers Enter events at the specified rate until the returned cancel is called.
func EnterLoop(bus Handler, frameDelay time.Duration) (cancel func()) {
	ch := make(chan struct{})
	go func() {
		ticker := time.NewTicker(frameDelay)
		frameDelayF64 := float64(frameDelay)
		lastTick := time.Now()
		framesElapsed := 0
		for {
			select {
			case now := <-ticker.C:
				deltaTime := now.Sub(lastTick)
				lastTick = now
				<-bus.Trigger(Enter.UnsafeEventID, EnterPayload{
					FramesElapsed:  framesElapsed,
					SinceLastFrame: deltaTime,
					TickPercent:    float64(deltaTime) / frameDelayF64,
				})
				framesElapsed++
			case <-ch:
				ticker.Stop()
				return
			}
		}
	}()
	return func() {
		// Q: why send here as well as close
		// A: to ensure that no more ticks are sent, the above goroutine has to
		//    acknowledge that it should stop and return-- just closing would
		//    enable code following this cancel function to assume no enters were
		//    being triggered when they still are.
		ch <- struct{}{}
		close(ch)
	}
}
