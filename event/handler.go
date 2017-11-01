package event

import (
	"math"
	"time"

	"github.com/oakmound/oak/timing"
)

var (
	_ Handler = &Bus{}
)

// Handler represents the necessary exported functions from an event.Bus
// for use in oak internally, and thus the functions that need to be replaced
// by alternative event handlers.
type Handler interface {
	UpdateLoop(framerate int, updateCh chan bool) error
	FramesElapsed() int
	SetTick(framerate int) error
	Update() error
	Flush() error
	Stop() error
	Reset()
	Trigger(event string, data interface{})
}

// UpdateLoop is expected to internally call Update()
// or do something equivalent at the given frameRate,
// sending signals to the sceneCh after each Update().
// Any flushing should be done as needed. This should
// not be called with `go`, if this requires goroutines
// it should create them itself.
// UpdateLoop is expected separately from Update() and
// Flush() because it will be more efficient for a Logical
// System to perform its own Updates outside of it’s exposed
// interface.
func (eb *Bus) UpdateLoop(framerate int, updateCh chan bool) error {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	ch := make(chan bool)
	eb.framesElapsed = 0
	eb.doneCh = ch
	eb.updateCh = updateCh
	go eb.ResolvePending()
	go func(doneCh chan bool) {
		eb.Ticker = timing.NewDynamicTicker()
		eb.Ticker.SetTick(timing.FPSToDuration(framerate))
		for {
			select {
			case <-eb.Ticker.C:
				<-eb.TriggerBack(Enter, eb.framesElapsed)
				eb.framesElapsed++
				eb.updateCh <- true
			case <-doneCh:
				eb.Ticker.Stop()
				doneCh <- true
				return
			}
		}
	}(ch)
	return nil
}

// Update updates all entities bound to this handler
func (eb *Bus) Update() error {
	<-eb.TriggerBack(Enter, eb.framesElapsed)
	return nil
}

// Flush refreshes any changes to the Handler’s bindings.
func (eb *Bus) Flush() error {
	if len(eb.unbindAllAndRebinds) > 0 {
		eb.resolveUnbindAllAndRebinds()
	}
	// Specific unbinds
	if len(eb.unbinds) > 0 {
		eb.resolveUnbinds()
	}

	// A full set of unbind settings
	if len(eb.fullUnbinds) > 0 {
		eb.resolveFullUnbinds()
	}

	// A partial set of unbind settings
	if len(eb.partUnbinds) > 0 {
		eb.resolvePartialUnbinds()
	}

	// Bindings
	if len(eb.binds) > 0 {
		eb.resolveBindings()
	}
	return nil
}

// Stop ceases anything spawned by an ongoing UpdateLoop
func (eb *Bus) Stop() error {
	eb.Ticker.SetTick(math.MaxInt32 * time.Second)
	select {
	case eb.doneCh <- true:
	case <-eb.updateCh:
		eb.doneCh <- true
	}
	<-eb.doneCh
	return nil
}

// FramesElapsed returns how many frames have elapsed since UpdateLoop was last called.
func (eb *Bus) FramesElapsed() int {
	return eb.framesElapsed
}

// SetTick optionally updates the Logical System’s tick rate
// (while it is looping) to be frameRate. If this operation is not
// supported, it should return an error.
func (eb *Bus) SetTick(framerate int) error {
	eb.Ticker.SetTick(timing.FPSToDuration(framerate))
	return nil
}
