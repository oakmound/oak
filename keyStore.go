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

func getMouseButton(i int32) string {
	s := ""
	switch i {
	case 1:
		s = "Left"
	case 2:
		s = "Right"
	case 3:
		s = "Middle"
	default:
		s = ""
	}
	return s
}
