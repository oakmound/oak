package timing

import (
	"time"

	"github.com/oakmound/oak/v3/dlog"
)

// A DynamicTicker is a ticker which can
// be sent signals in the form of durations to
// change how often it ticks.
type DynamicTicker struct {
	ticker    *time.Ticker
	C         chan time.Time
	resetCh   chan time.Duration
	forceTick chan bool
}

// NewDynamicTicker returns a null-initialized
// dynamic ticker.
func NewDynamicTicker() *DynamicTicker {
	dt := &DynamicTicker{
		ticker:    time.NewTicker(1000 * time.Hour),
		C:         make(chan time.Time),
		resetCh:   make(chan time.Duration),
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
	case at := <-dt.ticker.C:
		return dt.send(at)
	case tickDuration := <-dt.resetCh:
		dt.ticker.Stop()
		dt.ticker = time.NewTicker(tickDuration)
	case code := <-dt.forceTick:
		if code == forceStop {
			dt.close()
			return false
		}
		return dt.send(time.Time{})
	}
	return true
}

func (dt *DynamicTicker) send(v time.Time) bool {
	for {
		select {
		case r := <-dt.forceTick:
			if r == forceStop {
				dt.close()
				return false
			}
			continue
		case tickDuration := <-dt.resetCh:
			dt.ticker.Stop()
			dt.ticker = time.NewTicker(tickDuration)
			return true
		case dt.C <- v:
			return true
		}
	}
}

// SetTick changes the rate at which a dynamic ticker
// ticks
func (dt *DynamicTicker) SetTick(tickDuration time.Duration) {
	dt.resetCh <- tickDuration
}

func (dt *DynamicTicker) close() {
	close(dt.C)
	close(dt.resetCh)
	close(dt.forceTick)
}

const (
	forceTickOn = true
	forceStop   = false
)

// Step will force the dynamic ticker to tick, once.
// If the forced tick is not received, successive calls
// to Step will do nothing.
func (dt *DynamicTicker) Step() {
	select {
	case dt.forceTick <- forceTickOn:
	default:
	}
}

// ForceStep is the blocking equivalent to Step. After
// this is called, it won't return until the ticker has
// taken the forced step through. A potential use for this
// is in benchmarking how often the work between ticks
// can get done.
func (dt *DynamicTicker) ForceStep() {
	dt.forceTick <- forceTickOn
}

// Stop closes all internal channels and stops dt's internal ticker
func (dt *DynamicTicker) Stop() {
	defer func() {
		if x := recover(); x != nil {
			dlog.Error("Dynamic Ticker stopped twice")
		}
	}()
	dt.ticker.Stop()
	dt.forceTick <- forceStop
	<-dt.forceTick
}
