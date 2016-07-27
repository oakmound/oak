package plastic

import (
	"sync"
)

var (
	keyState = make(map[string]bool)
	keyLock  = sync.Mutex{}
)

// SetUp, SetDown, and IsDown all
// control access to a keystate map
// from key strings to down or up boolean
// states.
func setUp(key string) {
	keyLock.Lock()
	keyState[key] = false
	keyLock.Unlock()
}

func setDown(key string) {
	keyLock.Lock()
	keyState[key] = true
	keyLock.Unlock()
}

func IsDown(key string) bool {
	keyLock.Lock()
	k := keyState[key]
	keyLock.Unlock()
	return k
}
