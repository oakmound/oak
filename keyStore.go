package oak

import (
	"time"
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
func (c *Controller) SetUp(key string) {
	c.keyLock.Lock()
	c.durationLock.Lock()
	delete(c.keyState, key)
	delete(c.keyDurations, key)
	c.durationLock.Unlock()
	c.keyLock.Unlock()
}

// SetDown will cause later IsDown calls to report true
// for the given key. This is called internally when
// events are sent from the real keyboard and mouse.
// Calling this can interrupt real input or cause
// unintended behavior and should be done cautiously.
func (c *Controller) SetDown(key string) {
	c.keyLock.Lock()
	c.keyState[key] = true
	c.keyDurations[key] = time.Now()
	c.keyLock.Unlock()
}

// IsDown returns whether a key is held down
func (c *Controller) IsDown(key string) (k bool) {
	c.keyLock.RLock()
	k = c.keyState[key]
	c.keyLock.RUnlock()
	return
}

// IsHeld returns whether a key is held down, and for how long
// it has been held.
func (c *Controller) IsHeld(key string) (k bool, d time.Duration) {
	c.keyLock.RLock()
	k = c.keyState[key]
	c.keyLock.RUnlock()
	if k {
		c.durationLock.RLock()
		d = time.Since(c.keyDurations[key])
		c.durationLock.RUnlock()
	}
	return
}
