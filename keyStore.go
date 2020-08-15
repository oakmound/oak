package oak

import (
	"sync"
	"time"
)

var (
	keyState     = make(map[string]bool)
	keyDurations = make(map[string]time.Time)
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
	delete(keyDurations, key)
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
