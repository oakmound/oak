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
	TriggerCIDBack(cid CID, eventName string, data interface{}) chan struct{}
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
	eb.framesElapsed = 0
	eb.framerate = framerate
	frameDelay := timing.FPSToFrameDelay(framerate)
	if eb.Ticker == nil {
		eb.Ticker = time.NewTicker(frameDelay)
	}
	go eb.ResolveChanges()
	go func() {
		eb.Ticker.Reset(frameDelay)
		frameDelayF64 := float64(frameDelay)
		lastTick := time.Now()
		for {
			select {
			case now := <-eb.Ticker.C:
				deltaTime := now.Sub(lastTick)
				lastTick = now
				<-eb.TriggerBack(Enter, EnterPayload{
					FramesElapsed:  eb.framesElapsed,
					SinceLastFrame: deltaTime,
					TickPercent:    float64(deltaTime) / frameDelayF64,
				})
				eb.framesElapsed++
				select {
				case updateCh <- struct{}{}:
				case <-eb.doneCh:
					return
				}
			case <-eb.doneCh:
				return
			}
		}
	}()
	return nil
}

type EnterPayload struct {
	FramesElapsed  int
	SinceLastFrame time.Duration
	TickPercent    float64
}

// Update updates all entities bound to this handler
func (eb *Bus) Update() error {
	<-eb.TriggerBack(Enter, EnterPayload{
		FramesElapsed: eb.framesElapsed,
	})
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
	if eb.Ticker != nil {
		eb.Ticker.Stop()
	}
	eb.doneCh <- struct{}{}
	return nil
}

// Pause stops the event bus from running any further enter events
func (eb *Bus) Pause() {
	eb.Ticker.Reset(math.MaxInt32 * time.Second)
}

// Resume will resume emitting enter events
func (eb *Bus) Resume() {
	eb.Ticker.Reset(timing.FPSToFrameDelay(eb.framerate))
}

// FramesElapsed returns how many frames have elapsed since UpdateLoop was last called.
func (eb *Bus) FramesElapsed() int {
	return eb.framesElapsed
}

// SetTick optionally updates the Logical System’s tick rate
// (while it is looping) to be frameRate. If this operation is not
// supported, it should return an error.
func (eb *Bus) SetTick(framerate int) error {
	eb.Ticker.Reset(timing.FPSToFrameDelay(framerate))
	return nil
}
