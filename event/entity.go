package event

import (
	"sync"
)

var (
	highestID CID
	callers   = make([]Entity, 0)
	idMutex   = sync.Mutex{}
)

// Todo: callers having assigned buses?

// An Entity is an element which can be bound to,
// in that it has a CID. All Entities need to implement
// is an Init function which should call NextID(e) and
// return that id.
type Entity interface {
	Init() CID
}

// NextID finds the next available caller id (always incrementing)
// and returns it, after adding the given entity to
// the slice of callers at the returned index.
func NextID(e Entity) CID {
	idMutex.Lock()
	highestID++
	callers = append(callers, e)
	id := highestID
	idMutex.Unlock()
	return id
}

// GetEntity either returns callers[i-1]
// or nil, if there is nothing at that index.
func GetEntity(i int) interface{} {
	if HasEntity(i) {
		return callers[i-1]
	}
	return nil
}

// HasEntity returns whether the given caller id is an initialized entity
func HasEntity(i int) bool {
	return len(callers) >= i && i != 0
}

// DestroyEntity sets the index within the caller list to nil. Note that this
// does not reduce the size of the caller list, a potential change in the
// future would be to A) use a map or B) reassign caller ids to not directly
// correspond to indices within callers
func DestroyEntity(i int) {
	callers[i-1] = nil
}

// ResetEntities resets callers and highestID, effectively dropping the
// remaining entities from accessible memory.
func ResetEntities() {
	idMutex.Lock()
	highestID = 0
	callers = []Entity{}
	idMutex.Unlock()
}
