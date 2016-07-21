package event

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

var (
	thisBus = EventBus{make(map[string]map[int]*BindableStore)}
)

// This is a way of saying "Any function
// that takes a generic struct of data
// and returns an error can be bound".
type Bindable func(int, interface{}) error

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

// Stores data necessary
// to trace back to a bindable function
// and remove it from an eventBus
type Binding struct {
	BindingOption
	index int
}

type CID int

func GetEventBus() *EventBus {
	return &thisBus
}

func ResetEventBus() {
	thisBus = EventBus{make(map[string]map[int]*BindableStore)}
}

// Called by entities.
// Entities pass in a bindable function,
// and a set of options which
// are parsed out.
// Returns a binding that can used
// to unbind this binding later.
func (eb *EventBus) BindPriority(fn Bindable, opt BindingOption) (Binding, error) {

	list := eb.getBindableList(opt)
	i := list.storeBindable(fn)

	dlog.Info("Stored at", i)

	return Binding{opt, i}, nil
}

func GlobalBind(fn Bindable, name string) (Binding, error) {
	eb := GetEventBus()
	return eb.Bind(fn, name, 0)
}

func (eb *EventBus) Bind(fn Bindable, name string, callerID int) (Binding, error) {

	bOpt := BindingOption{}
	bOpt.Event = Event{
		Name:     name,
		CallerID: callerID,
	}

	dlog.Verb("Binding ", callerID, " with name ", name)

	return eb.BindPriority(fn, bOpt)
}

func (cid CID) Bind(fn Bindable, name string) (Binding, error) {
	eb := GetEventBus()
	return eb.Bind(fn, name, int(cid))
}

// Called by entities,
// for unbinding specific bindings.
func (eb *EventBus) Unbind(b Binding) error {

	list := eb.getBindableList(b.BindingOption)
	list.removeBindable(b)

	return nil
}

func (b Binding) Unbind() error {
	eb := GetEventBus()
	return eb.Unbind(b)
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

// Remove a bindable from a BindableList
func (bl *BindableList) removeBindable(b Binding) {
	i := b.index //store for messing with nextempty
	if len(bl.sl) < i+1 {
		return
	}
	bl.sl[i] = nil
	if i < bl.nextEmpty {
		bl.nextEmpty = i
	}
}

func (eb *EventBus) getBindableList(opt BindingOption) *BindableList {

	if _, ok := eb.bindingMap[opt.Name]; !ok {
		eb.bindingMap[opt.Name] = make(map[int]*BindableStore)
	}

	if _, ok := eb.bindingMap[opt.Name][opt.CallerID]; !ok {
		eb.bindingMap[opt.Name][opt.CallerID] = new(BindableStore)
		eb.bindingMap[opt.Name][opt.CallerID].defaultPriority = (new(BindableList))
	}

	store := eb.bindingMap[opt.Name][opt.CallerID]

	var list *BindableList

	// Default priority
	if opt.Priority == 0 {
		list = store.defaultPriority

		// High priority
	} else if opt.Priority > 0 {
		if store.highPriority[opt.Priority-1] == nil {
			store.highPriority[opt.Priority-1] = (new(BindableList))
		}

		if store.highIndex < opt.Priority {
			store.highIndex = opt.Priority
		}

		list = store.highPriority[opt.Priority-1]

		// Low priority
	} else {
		if store.lowPriority[(opt.Priority*-1)-1] == nil {
			store.lowPriority[(opt.Priority*-1)-1] = (new(BindableList))
		}

		if (store.lowIndex * -1) > opt.Priority {
			store.lowIndex = (-1 * opt.Priority)
		}

		list = store.lowPriority[(opt.Priority*-1)-1]
	}

	return list
}

// Called by entities or by game logic.
// Unbinds all events in this bus which
// match the given binding options.
func (eb *EventBus) UnbindAll(opt BindingOption) {

	var namekeys []string

	// If we were given a name,
	// we'll just iterate with that name.
	if opt.Name != "" {
		namekeys = append(namekeys, opt.Name)

		// Otherwise, iterate through all events.
	} else {
		for k := range eb.bindingMap {
			namekeys = append(namekeys, k)
		}
	}

	if opt.CallerID != 0 {
		for _, k := range namekeys {
			dlog.Verb("Deleting bindingMap: ", k, " callerId, ", opt.CallerID)
			delete(eb.bindingMap[k], opt.CallerID)
		}
	} else {
		for _, k := range namekeys {
			delete(eb.bindingMap, k)
		}
	}
	dlog.Verb(eb.bindingMap)
}

// Trigger an event, but only
// for one ID. Use case example:
// on onHit event
func (id CID) Trigger(eventName string, data interface{}) error {
	eb := GetEventBus()

	var err error

	if idMap, ok := eb.bindingMap[eventName]; ok {
		if bs, ok := idMap[int(id)]; ok {
			for i := bs.highIndex - 1; i >= 0; i-- {
				for _, bnd := range (*bs.highPriority[i]).sl {
					if bnd != nil {
						err = bnd(int(id), data)
						if err != nil {
							return err
						}
					}
				}
			}
			for _, bnd := range (bs.defaultPriority).sl {
				if bnd != nil {
					err = bnd(int(id), data)
					if err != nil {
						return err
					}
				}
			}
			for i := 0; i < bs.lowIndex; i++ {
				for _, bnd := range (*bs.lowPriority[i]).sl {
					if bnd != nil {
						err = bnd(int(id), data)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// Called externally by game logic
// and internally by plastic itself
// at specific integral points
func (eb_p *EventBus) Trigger(eventName string, data interface{}) error {

	eb := (*eb_p)

	var err error

	//dlog.Verb("Triggering, ", eventName)

	// Loop through all bindableStores for this eventName
	for id, bs := range eb.bindingMap[eventName] {
		// Loop through all bindables

		// Top to bottom, high priority
		for i := bs.highIndex - 1; i >= 0; i-- {
			for _, bnd := range (*bs.highPriority[i]).sl {
				if bnd != nil {
					err = bnd(id, data)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	for id, bs := range eb.bindingMap[eventName] {

		for _, bnd := range (bs.defaultPriority).sl {
			if bnd != nil {
				err = bnd(id, data)
				if err != nil {
					return err
				}
			}
		}
	}

	for id, bs := range eb.bindingMap[eventName] {
		// Bottom to top, low priority
		for i := 0; i < bs.lowIndex; i++ {
			for _, bnd := range (*bs.lowPriority[i]).sl {
				if bnd != nil {
					err = bnd(id, data)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
