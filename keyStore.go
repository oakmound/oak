package oak

import (
	"sync"
	"time"
)

var (
	keyBinds     = make(map[string]string)
	keyState     = make(map[string]bool)
	keyDurations = make(map[string]time.Duration)
	keyLock      = sync.RWMutex{}
	durationLock = sync.RWMutex{}
)

// SetUp, SetDown, and IsDown all
// control access to a keystate map
// from key strings to down or up boolean
// states.
func setUp(key string) {
	keyLock.Lock()
	durationLock.Lock()
	delete(keyState, key)
	delete(keyDurations, key)
	durationLock.Unlock()
	keyLock.Unlock()
}

func setDown(key string) {
	keyLock.Lock()
	keyState[key] = true
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
func IsHeld(key string) (k bool, d time.Duration) {
	keyLock.RLock()
	k = keyState[key]
	keyLock.RUnlock()
	if k {
		durationLock.RLock()
		d = keyDurations[key]
		durationLock.RUnlock()
	}
	return
}

// BindKey binds a name to be triggered when this
// key is triggered
func BindKey(key string, binding string) {
	keyBinds[key] = binding
}

// GetKeyBind returns either whatever name has been bound to
// a key or the key if nothing has been bound to it.
func GetKeyBind(key string) string {
	if v, ok := keyBinds[key]; ok {
		return v
	}
	return key
}

func keyHoldLoop() {
	var next time.Time
	var delta time.Duration
	now := time.Now()
	for {
		next = time.Now()
		delta = next.Sub(now)
		now = next
		keyLock.RLock()
		for k := range keyState {
			durationLock.Lock()
			keyDurations[k] += delta
			durationLock.Unlock()
		}
		keyLock.RUnlock()
	}
}
