package event

import (
	"sync"
	"sync/atomic"
)

// A CID is a caller ID that entities use to trigger and bind functionality
type CallerID int64

const Global CallerID = 0

// A CallerMap tracks CID mappings to Entities.
// This is an alternative to passing in the entity via closure scoping,
// and allows for more general bindings as simple top level functions.
type CallerMap struct {
	highestID   *int64
	callersLock sync.RWMutex
	callers     map[CallerID]interface{}
}

// NewCallerMap creates a caller map. A CallerMap
// is not valid for use if not created via this function.
func NewCallerMap() *CallerMap {
	return &CallerMap{
		highestID: new(int64),
		callers:   map[CallerID]interface{}{},
	}
}

// DefaultCallerMap is the caller map used by all event package caller
// functions.
var DefaultCallerMap = NewCallerMap()

// NextID finds the next available caller id
// and returns it, after adding the given entity to
// the caller map.
func (cm *CallerMap) Register(e interface{}) CallerID {
	nextID := atomic.AddInt64(cm.highestID, 1)
	cm.callersLock.Lock()
	cm.callers[CallerID(nextID)] = e
	cm.callersLock.Unlock()
	return CallerID(nextID)
}

// GetEntity returns the entity corresponding to the given ID within
// the caller map. If no entity is found, it returns nil.
func (cm *CallerMap) GetEntity(id CallerID) interface{} {
	cm.callersLock.RLock()
	defer cm.callersLock.RUnlock()
	return cm.callers[id]
}

// HasEntity returns whether the given caller id is an initialized entity
// within the caller map.
func (cm *CallerMap) HasEntity(id CallerID) bool {
	cm.callersLock.RLock()
	defer cm.callersLock.RUnlock()
	_, ok := cm.callers[id]
	return ok
}

// DestroyEntity removes an entity from the caller map.
func (cm *CallerMap) DestroyEntity(id CallerID) {
	cm.callersLock.Lock()
	delete(cm.callers, id)
	cm.callersLock.Unlock()
}
