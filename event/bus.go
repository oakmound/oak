package event

import (
	"sync"
	"time"

	"github.com/oakmound/oak/v3/oakerr"
)

// A Bus stores bindables to be triggered by events
type Bus struct {
	nextBindID         *int64
	bindingMap         map[UnsafeEventID]map[CallerID]bindableList
	persistentBindings []persistentBinding
	doneCh             chan struct{}
	framesElapsed      int
	ticker             *time.Ticker
	callerMap          *CallerMap

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
		doneCh:     make(chan struct{}),
		callerMap:  callerMap,
	}
}

// SetCallerMap updates a bus to use a specific set of callers.
func (bus *Bus) SetCallerMap(cm *CallerMap) {
	bus.callerMap = cm
}

// ClearPersistentBindings removes all persistent bindings. It will not unbind them
// from the bus, but they will not be bound following the next bus reset.
func (eb *Bus) ClearPersistentBindings() {
	eb.mutex.Lock()
	eb.persistentBindings = eb.persistentBindings[:0]
	eb.mutex.Unlock()
}

// Reset unbinds all present, non-persistent bindings on the bus.
func (bus *Bus) Reset() {
	bus.mutex.Lock()
	bus.bindingMap = make(map[UnsafeEventID]map[CallerID]bindableList)
	for _, pb := range bus.persistentBindings {
		bus.UnsafeBind(pb.eventID, pb.callerID, pb.fn)
	}
	bus.mutex.Unlock()
}

// EnterLoop triggers Enter events at the specified rate
func (bus *Bus) EnterLoop(frameDelay time.Duration) {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	bus.framesElapsed = 0
	if bus.ticker == nil {
		bus.ticker = time.NewTicker(frameDelay)
	}
	bus.doneCh = make(chan struct{})
	go func() {
		bus.ticker.Reset(frameDelay)
		frameDelayF64 := float64(frameDelay)
		lastTick := time.Now()
		for {
			select {
			case now := <-bus.ticker.C:
				deltaTime := now.Sub(lastTick)
				lastTick = now
				<-bus.Trigger(Enter.UnsafeEventID, EnterPayload{
					FramesElapsed:  bus.framesElapsed,
					SinceLastFrame: deltaTime,
					TickPercent:    float64(deltaTime) / frameDelayF64,
				})
				bus.framesElapsed++
			case <-bus.doneCh:
				return
			}
		}
	}()
}

// Stop ceases anything spawned by an ongoing UpdateLoop
func (bus *Bus) Stop() error {
	if bus.ticker != nil {
		bus.ticker.Stop()
		bus.ticker = nil
	}
	close(bus.doneCh)
	return nil
}

// SetTick optionally updates the Logical System’s tick rate
// (while it is looping) to be frameRate. If this operation is not
// supported, it should return an error.
func (bus *Bus) SetEnterLoopRate(frameDelay time.Duration) error {
	if bus.ticker == nil {
		return oakerr.NotFound{
			InputName: "bus.ticker",
		}
	}
	bus.ticker.Reset(frameDelay)
	return nil
}