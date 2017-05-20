package timing

import (
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

// A DynamicTicker is a ticker which can
// be sent signals in the form of durations to
// change how often it ticks.
type DynamicTicker struct {
	ticker    *time.Ticker
	C         chan time.Time
	resetCh   chan *time.Ticker
	forceTick chan bool
}

// NewDynamicTicker returns a null-initialized
// dynamic ticker
func NewDynamicTicker() *DynamicTicker {
	ch := make(chan time.Time)
	resetCh := make(chan *time.Ticker)
	forceTick := make(chan bool)
	dt := &DynamicTicker{
		// Please do not leave the application running
		// for a thousand hours without clicking on
		// the visualization knub, or else your next
		// visualization animation might skip a frame!
		// (We need -some- ticker defined or else
		// the program will crash in the following
		// routine on a nil pointer)
		ticker:    time.NewTicker(1000 * time.Hour),
		C:         ch,
		resetCh:   resetCh,
		forceTick: forceTick,
	}
	go func(dt *DynamicTicker) {
		for {
			select {
			case v := <-dt.ticker.C:
				select {
				case <-dt.forceTick:
					continue
				case dt.C <- v:
				case ticker := <-dt.resetCh:
					dt.ticker.Stop()
					dt.ticker = ticker
				}
			case ticker := <-dt.resetCh:
				dt.ticker.Stop()
				dt.ticker = ticker
			case r := <-dt.forceTick:
				if !r {
					close(dt.forceTick)
					close(dt.C)
					close(dt.resetCh)
					return
				}
				select {
				case <-dt.forceTick:
					continue
				case dt.C <- time.Time{}:
				}
			}
		}
	}(dt)
	return dt
}

// SetTick changes the rate at which a dynamic ticker
// ticks
func (dt *DynamicTicker) SetTick(d time.Duration) {
	dt.resetCh <- time.NewTicker(d)
}

// Step will force the dynamic ticker to tick, once.
// If the forced tick is not received, multiple calls
// to step will do nothing.
func (dt *DynamicTicker) Step() {
	select {
	case dt.forceTick <- true:
	default:
	}
}

// Stop closes all internal channels and stops dt's internal ticker
func (dt *DynamicTicker) Stop() {
	defer func() {
		if x := recover(); x != nil {
			dlog.Error("Dynamic Ticker stopped twice")
		}
	}()
	dt.ticker.Stop()
	select {
	case <-dt.C:
	default:
	}
	dt.forceTick <- false
}
