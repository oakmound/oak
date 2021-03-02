package oak

import (
	"sync"
	"time"
)

var (
	keyState     = make(map[string]bool)
	keyDurations = make(map[string]time.Time)
	lastUp       = make(map[string]time.Time)
	keyLock      = sync.RWMutex{}
	durationLock = sync.RWMutex{}
)

// SetUp, SetDown, and IsDown all
// control access to a keystate map
// from key strings to down or up boolean
// states.

// SetUp will cause later IsDown calls to report false
// for the given key. This is called internally when
// events are sent from the real keyboard and mouse.
// Calling this can interrupt real input or cause
// unintended behavior and should be done cautiously.
func SetUp(key string) {
	keyLock.Lock()
	durationLock.Lock()
	delete(keyState, key)

	// support GetAndClearLastUpTime
	lastDown := keyDurations[key]
	if !lastDown.IsZero() {
		lastUp[key] = time.Now()
		delete(keyDurations, key)
	}
	durationLock.Unlock()
	keyLock.Unlock()
}

// SetDown will cause later IsDown calls to report true
// for the given key. This is called internally when
// events are sent from the real keyboard and mouse.
// Calling this can interrupt real input or cause
// unintended behavior and should be done cautiously.
func SetDown(key string) {
	keyLock.Lock()
	keyState[key] = true
	keyDurations[key] = time.Now()
	keyLock.Unlock()
}

// IsDown returns whether a key is held down
func IsDown(key string) (k bool) {
	keyLock.RLock()
	k = keyState[key]
	keyLock.RUnlock()
	return
}

// GetAndClearLastUpTime returns the time of
// the last SetUp(key) event that followed
// a SetDown(key) event.
//
// It then clears the lastUp map of key.
// Subsequent queries will return the zero time.Time,
// as will keys that have never been pressed and
// released.
//
// My app was missing quick key presses, or thinking
// that one press was many. We solve both by
// having SetUp store into the lastUp map, and
// providing this method.
//
// Checking if a key was pressed becomes reliable. Simply:
//
// if !GetAndClearLastUpTime(key).IsZero() {
//     // key was released since last we checked.
// }
//
// Note that only one client, not many, can expect
// to do this and know if there has been at least one key up
// since the last poll, since we clear the lastUp time
// on purpose.
//
func GetAndClearLastUpTime(key string) (tm time.Time) {
	keyLock.Lock()
	tm = lastUp[key]
	delete(lastUp, key)
	keyLock.Unlock()
	return
}

// IsHeld returns whether a key is held down, and for how long
// it has been held.
func IsHeld(key string) (k bool, d time.Duration) {
	keyLock.RLock()
	k = keyState[key]
	keyLock.RUnlock()
	if k {
		durationLock.RLock()
		d = time.Since(keyDurations[key])
		durationLock.RUnlock()
	}
	return
}
