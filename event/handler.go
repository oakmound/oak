package event

import "github.com/oakmound/oak/timing"

// Handler represents the necessary exported functions from an event.Bus
// for use in oak internally, and thus the functions that need to be replaced
// by alternative event handlers.
type Handler interface {
	UpdateLoop(framerate int, updateCh chan<- bool) error
	FramesElapsed() int
	SetTick(framerate int) error
	Update() error
	Flush() error
	Stop() error
	Reset()
	Trigger(event string, data interface{})
}

func (eb *Bus) UpdateLoop(framerate int, updateCh chan<- bool) error {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	ch := make(chan bool)
	eb.framesElapsed = 0
	eb.doneCh = ch
	go func(doneCh chan bool) {
		eb.Ticker = timing.NewDynamicTicker()
		eb.Ticker.SetTick(timing.FPSToDuration(framerate))
		for {
			select {
			case <-eb.Ticker.C:
				<-eb.TriggerBack("EnterFrame", eb.framesElapsed)
				eb.framesElapsed++
				updateCh <- true
			case <-doneCh:
				eb.Ticker.Stop()
				return
			}
		}
	}(ch)
	return nil
}

func (eb *Bus) Update() error {
	<-eb.TriggerBack("EnterFrame", eb.framesElapsed)
	return nil
}

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

func (eb *Bus) Stop() error {
	close(eb.doneCh)
	return nil
}

func (eb *Bus) FramesElapsed() int {
	return eb.framesElapsed
}

func (eb *Bus) SetTick(framerate int) error {
	eb.Ticker.SetTick(timing.FPSToDuration(framerate))
	return nil
}
