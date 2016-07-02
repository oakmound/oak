package plastic

import (
	"sync"
)

var (
	keyState = make(map[string]bool)
	keyLock  = sync.Mutex{}
)

func SetUp(key string) {
	keyLock.Lock()
	keyState[key] = false
	keyLock.Unlock()
}

func SetDown(key string) {
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
