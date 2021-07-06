package event

import (
	"sync"
	"sync/atomic"
)

// A CallerMap tracks CID mappings to Entities. Its intended use is
// to be a source of truth within event bindings for what entity the
// binding is triggering on:
// 		var cm *event.CallerMap
//		func(cid event.CID, payload interface{}) int {
//			ent := cm.GetEntity(cid)
//          f, ok := ent.(*Foo)
//          if !ok {
//				// bound to an unexpected entity type!
//				return event.UnbindSingle
//			}
//          // ...
//		}
// This is an alternative to passing in the entity via closure scoping,
// and allows for more general bindings as simple top level functions.
type CallerMap struct {
	highestID   *int64
	callersLock sync.RWMutex
	callers     map[CID]Entity
}

// NewCallerMap creates a caller map. A CallerMap
// is not valid for use if not created via this function.
func NewCallerMap() *CallerMap {
	return &CallerMap{
		highestID: new(int64),
		callers:   map[CID]Entity{},
	}
}

var DefaultCallerMap = NewCallerMap()

// NextID finds the next available caller id
// and returns it, after adding the given entity to
// the caller map.
func (cm *CallerMap) NextID(e Entity) CID {
	nextID := atomic.AddInt64(cm.highestID, 1)
	cm.callersLock.Lock()
	cm.callers[CID(nextID)] = e
	cm.callersLock.Unlock()
	return CID(nextID)
}

// GetEntity returns the entity corresponding to the given ID within
// the caller map. If no entity is found, it returns nil.
func (cm *CallerMap) GetEntity(id CID) Entity {
	cm.callersLock.RLock()
	defer cm.callersLock.RUnlock()
	return cm.callers[id]
}

// HasEntity returns whether the given caller id is an initialized entity
// within the caller map.
func (cm *CallerMap) HasEntity(id CID) bool {
	cm.callersLock.RLock()
	defer cm.callersLock.RUnlock()
	_, ok := cm.callers[id]
	return ok
}

// DestroyEntity removes an entity from the caller map.
func (cm *CallerMap) DestroyEntity(id CID) {
	cm.callersLock.Lock()
	delete(cm.callers, id)
	cm.callersLock.Unlock()
}

// NextID finds the next available caller id
// and returns it, after adding the given entity to
// the default caller map.
func NextID(e Entity) CID {
	return DefaultCallerMap.NextID(e)
}

// GetEntity returns the entity corresponding to the given ID within
// the default caller map. If no entity is found, it returns nil.
func GetEntity(id CID) Entity {
	return DefaultCallerMap.GetEntity(id)
}

// HasEntity returns whether the given caller id is an initialized entity
// within the default caller map.
func HasEntity(id CID) bool {
	return DefaultCallerMap.HasEntity(id)
}

// DestroyEntity removes an entity from the default caller map.
func DestroyEntity(id CID) {
	DefaultCallerMap.DestroyEntity(id)
}

// ResetCallerMap resets the DefaultCallerMap to be empty.
func ResetCallerMap() {
	*DefaultCallerMap = *NewCallerMap()
}
