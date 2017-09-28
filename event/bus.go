package event

import (
	"reflect"
	"sync"

	"github.com/oakmound/oak/timing"
)

// Bindable is a way of saying "Any function
// that takes a generic struct of data
// and returns an error can be bound".
type Bindable func(int, interface{}) int

// BindableList just stores other relevant data
// that a list of bindables needs to
// operate efficiently
type bindableList struct {
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

type bindableStore struct {
	// Priorities can be in a range
	// from -32 to 32. Below 0,
	// goes into lowPriority.
	// Above 0, goes into highPriority.
	// No priority makes it default to
	// take place between high and low.
	lowPriority     [32]*bindableList
	defaultPriority *bindableList
	highPriority    [32]*bindableList

	// These internally keep track
	// where we start / stop checking
	// our non-default binding lists.
	lowIndex  int
	highIndex int
}

// A Bus stores bindables to be triggered by events
type Bus struct {
	bindingMap          map[string]map[int]*bindableStore
	doneCh              chan bool
	updateCh            chan bool
	framesElapsed       int
	Ticker              *timing.DynamicTicker
	binds               []UnbindOption
	partUnbinds         []BindingOption
	fullUnbinds         []UnbindOption
	unbinds             []binding
	unbindAllAndRebinds []UnbindAllOption

	mutex        sync.RWMutex
	pendingMutex sync.Mutex

	init sync.Once
}

// NewBus returns an empty event bus
func NewBus() *Bus {
	return &Bus{
		bindingMap:          make(map[string]map[int]*bindableStore),
		doneCh:              make(chan bool),
		binds:               make([]UnbindOption, 0),
		partUnbinds:         make([]BindingOption, 0),
		fullUnbinds:         make([]UnbindOption, 0),
		unbinds:             make([]binding, 0),
		unbindAllAndRebinds: make([]UnbindAllOption, 0),
		mutex:               sync.RWMutex{},
		pendingMutex:        sync.Mutex{},
		init:                sync.Once{},
	}
}

// An Event is an event name and an associated caller id
type Event struct {
	Name     string
	CallerID int
}

// BindingOption is all the information required
// to bind something
type BindingOption struct {
	Event
	Priority int
}

// UnbindOption stores information necessary
// to unbind a bindable
type UnbindOption struct {
	BindingOption
	Fn Bindable
}

// binding stores data necessary
// to trace back to a bindable function
// and remove it from a Bus.
type binding struct {
	BindingOption
	index int
}

// Reset empties out all transient portions of the bus. It will not stop
// an ongoing loop.
func (eb *Bus) Reset() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	eb.bindingMap = make(map[string]map[int]*bindableStore)
	eb.binds = []UnbindOption{}
	eb.partUnbinds = []BindingOption{}
	eb.fullUnbinds = []UnbindOption{}
	eb.unbinds = []binding{}
	eb.unbindAllAndRebinds = []UnbindAllOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

// UnbindAllOption stores information needed to unbind and rebind
type UnbindAllOption struct {
	ub   BindingOption
	bs   []BindingOption
	bnds []Bindable
}

// Store a bindable into a BindableList.
func (bl *bindableList) storeBindable(fn Bindable) int {

	i := bl.nextEmpty
	if len(bl.sl) == i {
		bl.sl = append(bl.sl, fn)
	} else {
		bl.sl[i] = fn
	}

	// Find the next empty space
	for len(bl.sl) != bl.nextEmpty && bl.sl[bl.nextEmpty] != nil {
		bl.nextEmpty++
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
// returning UnbindSingle from the function
// itself is much safer!
func (bl *bindableList) removeBindable(fn Bindable) {
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
func (bl *bindableList) removeBinding(b binding) {
	bl.removeIndex(b.index)
}

func (bl *bindableList) removeIndex(i int) {
	if len(bl.sl) <= i {
		return
	}

	bl.sl[i] = nil

	if i < bl.nextEmpty {
		bl.nextEmpty = i
	}
}

func (eb *Bus) getBindableList(opt BindingOption) *bindableList {

	if m := eb.bindingMap[opt.Name]; m == nil {
		eb.bindingMap[opt.Name] = make(map[int]*bindableStore)
	}

	if m := eb.bindingMap[opt.Name][opt.CallerID]; m == nil {
		eb.bindingMap[opt.Name][opt.CallerID] = new(bindableStore)
		eb.bindingMap[opt.Name][opt.CallerID].defaultPriority = (new(bindableList))
	}

	store := eb.bindingMap[opt.Name][opt.CallerID]

	var prio *bindableList
	// Default priority
	if opt.Priority == 0 {
		prio = store.defaultPriority

		// High priority
	} else if opt.Priority > 0 {
		if store.highPriority[opt.Priority-1] == nil {
			store.highPriority[opt.Priority-1] = (new(bindableList))
		}

		if store.highIndex < opt.Priority {
			store.highIndex = opt.Priority
		}

		prio = store.highPriority[opt.Priority-1]

		// Low priority
	} else {
		if store.lowPriority[(opt.Priority*-1)-1] == nil {
			store.lowPriority[(opt.Priority*-1)-1] = (new(bindableList))
		}

		if (store.lowIndex * -1) > opt.Priority {
			store.lowIndex = (-1 * opt.Priority)
		}

		prio = store.lowPriority[(opt.Priority*-1)-1]
	}
	return prio
}
