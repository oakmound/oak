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

func IsDown(key string) (k bool) {
	keyLock.RLock()
	k = keyState[key]
	keyLock.RUnlock()
	return
}
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

func BindKey(key string, binding string) {
	keyBinds[key] = binding
}

func GetKeyBind(key string) string {
	if v, ok := keyBinds[key]; ok {
		return v
	}
	return key
}

func KeyHoldLoop() {
	var next time.Time
	var delta time.Duration
	now := time.Now()
	for {
		next = time.Now()
		delta = next.Sub(now)
		now = next
		keyLock.RLock()
		for k, _ := range keyState {
			durationLock.Lock()
			keyDurations[k] += delta
			durationLock.Unlock()
		}
		keyLock.RUnlock()
	}
}
