package render

import (
	"sync"
)

var (
	dLock    = sync.Mutex{}
	debugMap = make(map[string]Renderable)
)

func UpdateDebugMap(rName string, r Renderable) {
	dLock.Lock()
	debugMap[rName] = r
	dLock.Unlock()
}

func GetDebugRenderable(rName string) (Renderable, bool) {
	dLock.Lock()
	tmp, ok := debugMap[rName]
	dLock.Unlock()
	return tmp, ok
}
