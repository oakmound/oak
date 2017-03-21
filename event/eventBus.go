// Package event propogates events through entities with given caller IDs.
// It sets up a subscribe-publish model with the Bind and Trigger functions.
// In a slight change to the sub-pub model, event allows bindings to occur
// in an explicit order through assigning priority to individual bind calls.
package event

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"

	"reflect"
	"strconv"
	"sync"
)

var (
	thisBus = EventBus{make(map[string]map[int]*BindableStore)}
	mutex   = sync.RWMutex{}
	rLocks  = 0
)

const (
	NO_RESPONSE = iota
	ERROR
	// UNBIND_EVENT unbinds everything for a specific
	// event name from an entity at the bindable's
	// priority.
	UNBIND_EVENT
	// We can't unbind a single bindable efficiently,
	// so UNBIND_EVENT is recommended.
	UNBIND_SINGLE
)

// This is a way of saying "Any function
// that takes a generic struct of data
// and returns an error can be bound".
type Bindable func(int, interface{}) int

// This just stores other relevant data
// that a list of bindables needs to
// operate efficiently
type BindableList struct {
	sl []Bindable
	// We keep track of where the next nil
	// element in our list is, so we
	// can let bindings know where they
	// are by index, (we don't shift to
	// fill empty spaces) and so we can
	// fill that slot next when a
	// new binding comes in.
	nextEmpty int
}

type BindableStore struct {
	// Priorities can be in a range
	// from -32 to 32. Below 0,
	// goes into lowPriority.
	// Above 0, goes into highPriority.
	// No priority makes it default to
	// take place between high and low.
	lowPriority     [32]*BindableList
	defaultPriority *BindableList
	highPriority    [32]*BindableList

	// These internally keep track
	// where we start / stop checking
	// our non-default binding lists.
	lowIndex  int
	highIndex int
}

type EventBus struct {
	bindingMap map[string]map[int]*BindableStore
}

// We keep a reference to caller
// in case an entity wants to
// unbind events related to itself
// (or some other entity)
type Event struct {
	Name     string
	CallerID int
}

// Populated by callers of Bind.
type BindingOption struct {
	Event
	Priority int
}

type UnbindOption struct {
	BindingOption
	fn Bindable
}

// Stores data necessary
// to trace back to a bindable function
// and remove it from an eventBus
type Binding struct {
	BindingOption
	index int
}

type CID int

func (cid CID) String() string {
	return strconv.Itoa(int(cid))
}

func (cid CID) E() interface{} {
	return GetEntity(int(cid))
}

func GetEventBus() *EventBus {
	return &thisBus
}

func ResetEventBus() {
	thisBus = EventBus{make(map[string]map[int]*BindableStore)}
	mutex.Lock()
	bindablesToBind = []Bindable{}
	optionsToBind = []BindingOption{}

	optionsToUnbind = []BindingOption{}
	ubOptionsToUnbind = []UnbindOption{}
	bindingsToUnbind = []Binding{}

	orderedUnbinds = []BindingOption{}
	orderedBindOptions = []BindingOption{}
	orderedBindables = []Bindable{}

	mutex.Unlock()
}

var (
	bindablesToBind = []Bindable{}
	optionsToBind   = []BindingOption{}

	optionsToUnbind   = []BindingOption{}
	ubOptionsToUnbind = []UnbindOption{}
	bindingsToUnbind  = []Binding{}

	orderedUnbinds     = []BindingOption{}
	orderedBindOptions = []BindingOption{}
	orderedBindables   = []Bindable{}

	pendingMutex = sync.Mutex{}
)

func ResolvePending() {

	if len(orderedUnbinds) > 0 {
		mutex.Lock()
		pendingMutex.Lock()
		for _, opt := range orderedUnbinds {
			var namekeys []string
			// If we were given a name,
			// we'll just iterate with that name.
			if opt.Name != "" {
				namekeys = append(namekeys, opt.Name)

				// Otherwise, iterate through all events.
			} else {
				for k := range thisBus.bindingMap {
					namekeys = append(namekeys, k)
				}
			}

			if opt.CallerID != 0 {
				for _, k := range namekeys {
					delete(thisBus.bindingMap[k], opt.CallerID)
				}
			} else {
				for _, k := range namekeys {
					delete(thisBus.bindingMap, k)
				}
			}
			dlog.Verb(thisBus.bindingMap)
		}

		// Bindings
		for i := 0; i < len(orderedBindables); i++ {
			fn := orderedBindables[i]
			opt := orderedBindOptions[i]
			list := thisBus.getBindableList(opt)
			list.storeBindable(fn)
		}

		mutex.Unlock()

		orderedUnbinds = []BindingOption{}
		orderedBindables = []Bindable{}
		orderedBindOptions = []BindingOption{}
		pendingMutex.Unlock()
	}

	// Unbinds
	if len(bindingsToUnbind) > 0 {
		mutex.Lock()
		pendingMutex.Lock()
		for _, b := range bindingsToUnbind {
			thisBus.getBindableList(b.BindingOption).removeBinding(b)
		}
		mutex.Unlock()

		bindingsToUnbind = []Binding{}
		pendingMutex.Unlock()
	}

	if len(optionsToUnbind) > 0 {
		mutex.Lock()
		pendingMutex.Lock()
		for _, opt := range optionsToUnbind {

			var namekeys []string

			// If we were given a name,
			// we'll just iterate with that name.
			if opt.Name != "" {
				namekeys = append(namekeys, opt.Name)

				// Otherwise, iterate through all events.
			} else {
				for k := range thisBus.bindingMap {
					namekeys = append(namekeys, k)
				}
			}

			if opt.CallerID != 0 {
				for _, k := range namekeys {
					delete(thisBus.bindingMap[k], opt.CallerID)
				}
			} else {
				for _, k := range namekeys {
					delete(thisBus.bindingMap, k)
				}
			}
			dlog.Verb(thisBus.bindingMap)
		}
		mutex.Unlock()

		optionsToUnbind = []BindingOption{}
		pendingMutex.Unlock()
	}

	// ubOptions need to be fully populated,
	// unlike optionsToUnbind
	if len(ubOptionsToUnbind) > 0 {
		mutex.Lock()
		pendingMutex.Lock()

		for _, opt := range ubOptionsToUnbind {
			thisBus.getBindableList(opt.BindingOption).removeBindable(opt.fn)
		}
		mutex.Unlock()

		ubOptionsToUnbind = []UnbindOption{}
		pendingMutex.Unlock()
	}

	if len(bindablesToBind) > 0 {
		mutex.Lock()
		pendingMutex.Lock()
		// Bindings
		for i := 0; i < len(bindablesToBind); i++ {
			fn := bindablesToBind[i]
			opt := optionsToBind[i]
			list := thisBus.getBindableList(opt)
			list.storeBindable(fn)
		}
		mutex.Unlock()

		bindablesToBind = []Bindable{}
		optionsToBind = []BindingOption{}
		pendingMutex.Unlock()
	}
}

// Store a bindable into a BindableList.
func (bl_p *BindableList) storeBindable(fn Bindable) int {

	bl := (*bl_p)

	i := bl.nextEmpty
	if len(bl.sl) == i {
		bl_p.sl = append(bl.sl, fn)
	} else {
		bl_p.sl[i] = fn
	}

	// Find the next empty space
	for len(bl_p.sl) != bl_p.nextEmpty && bl_p.sl[bl_p.nextEmpty] != nil {
		bl_p.nextEmpty++
	}

	return i
}

// This scans linearly for the bindable
// This will cause an issue with closures!
// You can't unbind closures that don't have the
// same variable reference because this compares
// pointers!
//
// At all costs, this should be avoided, and
// returning "UNBIND_SINGLE" from the function
// itself is much safer!
func (bl *BindableList) removeBindable(fn Bindable) {
	i := 0
	v := reflect.ValueOf(fn)
	for {
		v2 := reflect.ValueOf(bl.sl[i])
		if v2 == v {
			bl.removeIndex(i)
			return
		}
		i++
	}
}

// Remove a bindable from a BindableList
func (bl *BindableList) removeBinding(b Binding) {
	bl.removeIndex(b.index)
}

func (bl *BindableList) removeIndex(i int) {
	if len(bl.sl) < i+1 {
		return
	}

	bl.sl[i] = nil

	if i < bl.nextEmpty {
		bl.nextEmpty = i
	}
}

func (eb *EventBus) getBindableList(opt BindingOption) *BindableList {

	if m, ok := eb.bindingMap[opt.Name]; !ok || m == nil {
		eb.bindingMap[opt.Name] = make(map[int]*BindableStore)
	}

	if m, ok := eb.bindingMap[opt.Name][opt.CallerID]; !ok || m == nil {
		eb.bindingMap[opt.Name][opt.CallerID] = new(BindableStore)
		eb.bindingMap[opt.Name][opt.CallerID].defaultPriority = (new(BindableList))
	}

	store := eb.bindingMap[opt.Name][opt.CallerID]

	// Default priority
	if opt.Priority == 0 {
		return store.defaultPriority

		// High priority
	} else if opt.Priority > 0 {
		if store.highPriority[opt.Priority-1] == nil {
			store.highPriority[opt.Priority-1] = (new(BindableList))
		}

		if store.highIndex < opt.Priority {
			store.highIndex = opt.Priority
		}

		return store.highPriority[opt.Priority-1]

		// Low priority
	} else {
		if store.lowPriority[(opt.Priority*-1)-1] == nil {
			store.lowPriority[(opt.Priority*-1)-1] = (new(BindableList))
		}

		if (store.lowIndex * -1) > opt.Priority {
			store.lowIndex = (-1 * opt.Priority)
		}

		return store.lowPriority[(opt.Priority*-1)-1]
	}
}
