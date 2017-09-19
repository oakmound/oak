package timing

import (
	"time"

	"github.com/oakmound/oak/dlog"
)

// BUG: There's some circumstance where DynamicTickers
// are not properly initialized/closed in the engine

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
	dt := &DynamicTicker{
		ticker:    time.NewTicker(1000 * time.Hour),
		C:         make(chan time.Time),
		resetCh:   make(chan *time.Ticker),
		forceTick: make(chan bool),
	}
	go func() {
		for dt.loop() {
		}
	}()
	return dt
}

func (dt *DynamicTicker) loop() bool {
	select {
	case v := <-dt.ticker.C:
		for {
			select {
			case r := <-dt.forceTick:
				if !r {
					dt.close()
					return false
				}
				continue
			case ticker := <-dt.resetCh:
				dt.ticker.Stop()
				dt.ticker = ticker
				return true
			case dt.C <- v:
				return true
			}
		}
	case ticker := <-dt.resetCh:
		dt.ticker.Stop()
		dt.ticker = ticker
	case r := <-dt.forceTick:
		if !r {
			dt.close()
			return false
		}
		for {
			select {
			case r := <-dt.forceTick:
				if !r {
					dt.close()
					return false
				}
				continue
			case ticker := <-dt.resetCh:
				dt.ticker.Stop()
				dt.ticker = ticker
				return true
			case dt.C <- time.Time{}:
				return true
			}
		}
	}
	return true
}

// SetTick changes the rate at which a dynamic ticker
// ticks
func (dt *DynamicTicker) SetTick(d time.Duration) {
	dt.resetCh <- time.NewTicker(d)
}

func (dt *DynamicTicker) close() {
	close(dt.C)
	close(dt.resetCh)
	close(dt.forceTick)
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
	dt.forceTick <- false
	<-dt.forceTick
}
