package event

import (
	"reflect"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/timing"
)

// Bindable is a way of saying "Any function
// that takes a generic struct of data
// and returns an error can be bound".
type Bindable func(CID, interface{}) int

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

// A Bus stores bindables to be triggered by events
type Bus struct {
	bindingMap          map[string]map[CID]*bindableList
	doneCh              chan struct{}
	updateCh            chan struct{}
	framesElapsed       int
	Ticker              *timing.DynamicTicker
	binds               []UnbindOption
	partUnbinds         []Event
	fullUnbinds         []UnbindOption
	unbinds             []binding
	unbindAllAndRebinds []UnbindAllOption
	framerate           int
	refreshRate         time.Duration

	mutex        sync.RWMutex
	pendingMutex sync.Mutex

	init sync.Once
}

// NewBus returns an empty event bus
func NewBus() *Bus {
	return &Bus{
		Ticker:              timing.NewDynamicTicker(),
		bindingMap:          make(map[string]map[CID]*bindableList),
		updateCh:            make(chan struct{}),
		doneCh:              make(chan struct{}),
		binds:               make([]UnbindOption, 0),
		partUnbinds:         make([]Event, 0),
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
	CallerID CID
}

// UnbindOption stores information necessary
// to unbind a bindable
type UnbindOption struct {
	Event
	Fn Bindable
}

// binding stores data necessary
// to trace back to a bindable function
// and remove it from a Bus.
type binding struct {
	Event
	index int
}

// Reset empties out all transient portions of the bus. It will not stop
// an ongoing loop.
func (eb *Bus) Reset() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	eb.bindingMap = make(map[string]map[CID]*bindableList)
	eb.binds = []UnbindOption{}
	eb.partUnbinds = []Event{}
	eb.fullUnbinds = []UnbindOption{}
	eb.unbinds = []binding{}
	eb.unbindAllAndRebinds = []UnbindAllOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

// UnbindAllOption stores information needed to unbind and rebind
type UnbindAllOption struct {
	ub   Event
	bs   []Event
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
	v := reflect.ValueOf(fn)
	for i := 0; i < len(bl.sl); i++ {
		v2 := reflect.ValueOf(bl.sl[i])
		if v2 == v {
			bl.removeIndex(i)
			return
		}
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

func (eb *Bus) getBindableList(opt Event) *bindableList {
	if m := eb.bindingMap[opt.Name]; m == nil {
		eb.bindingMap[opt.Name] = make(map[CID]*bindableList)
	}
	if m := eb.bindingMap[opt.Name][opt.CallerID]; m == nil {
		eb.bindingMap[opt.Name][opt.CallerID] = new(bindableList)
	}
	return eb.bindingMap[opt.Name][opt.CallerID]
}
