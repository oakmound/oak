package event

import (
	"math"
	"time"

	"github.com/oakmound/oak/v3/timing"
)

var (
	_ Handler = &Bus{}
)

// Handler represents the necessary exported functions from an event.Bus
// for use in oak internally, and thus the functions that need to be replaced
// by alternative event handlers.
// TODO V3: consider breaking down the bus into smaller components
// for easier composition for external handler implementations
type Handler interface {
	// <Handler>
	UpdateLoop(framerate int, updateCh chan struct{}) error
	FramesElapsed() int
	SetTick(framerate int) error
	Update() error
	Flush() error
	Stop() error
	Reset()
	SetRefreshRate(time.Duration)
	// <Triggerer>
	Trigger(event string, data interface{})
	TriggerBack(event string, data interface{}) chan struct{}
	// <Pauser>
	Pause()
	Resume()
	// <Binder>
	Bind(string, CID, Bindable)
	GlobalBind(string, Bindable)
	UnbindAll(Event)
	UnbindAllAndRebind(Event, []Bindable, CID, []string)
	UnbindBindable(UnbindOption)
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
func (eb *Bus) UpdateLoop(framerate int, updateCh chan struct{}) error {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	ch := make(chan struct{})
	eb.framesElapsed = 0
	eb.doneCh = ch
	eb.updateCh = updateCh
	eb.framerate = framerate
	eb.Ticker = timing.NewDynamicTicker()
	go eb.ResolvePending()
	go func(doneCh chan struct{}) {
		eb.Ticker.SetTick(timing.FPSToFrameDelay(framerate))
		for {
			select {
			case <-eb.Ticker.C:
				<-eb.TriggerBack(Enter, eb.framesElapsed)
				eb.framesElapsed++
				eb.updateCh <- struct{}{}
			case <-doneCh:
				eb.Ticker.Stop()
				doneCh <- struct{}{}
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
	case eb.doneCh <- struct{}{}:
	case <-eb.updateCh:
		eb.doneCh <- struct{}{}
	}
	<-eb.doneCh
	return nil
}

// Pause stops the event bus from running any further enter events
func (eb *Bus) Pause() {
	eb.Ticker.SetTick(math.MaxInt32 * time.Second)
}

// Resume will resume emitting enter events
func (eb *Bus) Resume() {
	eb.Ticker.SetTick(timing.FPSToFrameDelay(eb.framerate))
}

// FramesElapsed returns how many frames have elapsed since UpdateLoop was last called.
func (eb *Bus) FramesElapsed() int {
	return eb.framesElapsed
}

// SetTick optionally updates the Logical System’s tick rate
// (while it is looping) to be frameRate. If this operation is not
// supported, it should return an error.
func (eb *Bus) SetTick(framerate int) error {
	eb.Ticker.SetTick(timing.FPSToFrameDelay(framerate))
	return nil
}
