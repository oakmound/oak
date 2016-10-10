package event

import (
	"sync"
)

var (
	highestID CID = 0
	callers       = make([]Entity, 0)
	idMutex       = sync.Mutex{}
)

type Entity interface {
	Init() CID
}

func NextID(e Entity) CID {
	idMutex.Lock()
	highestID++
	callers = append(callers, e)
	id := highestID
	idMutex.Unlock()
	return id
}

func GetEntity(i int) interface{} {
	if HasEntity(i) {
		return callers[i-1]
	}
	return nil
}

func HasEntity(i int) bool {
	return len(callers) >= i
}

func DestroyEntity(i int) {
	callers[i-1] = nil
}

func ResetEntities() {
	idMutex.Lock()
	highestID = 0
	callers = []Entity{}
	idMutex.Unlock()
}
