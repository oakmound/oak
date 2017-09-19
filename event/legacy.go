package event

// As in collision and mouse, legacy.go lists functions that
// only operate on DefaultBus, a package global bus.

var (
	// DefaultBus is a bus that has additional operations for CIDs, and can
	// be called via event.Call as opposed to bus.Call
	DefaultBus = NewBus()
)

// Trigger an event, but only for one ID, on the default bus
func (cid CID) Trigger(eventName string, data interface{}) {

	go func(eventName string, data interface{}) {
		DefaultBus.mutex.RLock()
		iid := int(cid)
		if idMap, ok := DefaultBus.bindingMap[eventName]; ok {
			if bs, ok := idMap[iid]; ok {
				for i := bs.highIndex - 1; i >= 0; i-- {
					DefaultBus.triggerDefault((*bs.highPriority[i]).sl, iid, eventName, data)
				}
				DefaultBus.triggerDefault((bs.defaultPriority).sl, iid, eventName, data)

				for i := 0; i < bs.lowIndex; i++ {
					DefaultBus.triggerDefault((*bs.lowPriority[i]).sl, iid, eventName, data)
				}
			}
		}
		DefaultBus.mutex.RUnlock()
	}(eventName, data)
}

// Bind on a CID is shorthand for bus.Bind(fn, name, cid), on the default bus.
func (cid CID) Bind(fn Bindable, name string) {
	DefaultBus.Bind(fn, name, int(cid))
}

// BindPriority on a CID is shorthand for bus.BindPriority(fn, ...), on the default bus.
func (cid CID) BindPriority(fn Bindable, name string, priority int) {
	DefaultBus.BindPriority(fn, BindingOption{
		Event{
			name,
			int(cid),
		},
		priority,
	})
}

// UnbindAll removes all events with the given cid from the event bus
func (cid CID) UnbindAll() {
	DefaultBus.UnbindAll(BindingOption{
		Event{
			"",
			int(cid),
		},
		0,
	})
}

// UnbindAllAndRebind on a CID is equivalent to bus.UnbindAllAndRebind(..., cid)
func (cid CID) UnbindAllAndRebind(binds []Bindable, events []string) {
	DefaultBus.UnbindAllAndRebind(BindingOption{
		Event{
			"",
			int(cid),
		},
		0,
	}, binds, int(cid), events)
}

// Trigger is equivalent to bus.Trigger(...)
// Todo: move this to legacy.go, see mouse or collision
func Trigger(eventName string, data interface{}) {
	DefaultBus.Trigger(eventName, data)
}

// TriggerBack is equivalent to bus.TriggerBack(...)
func TriggerBack(eventName string, data interface{}) chan bool {
	return DefaultBus.TriggerBack(eventName, data)
}

// GlobalBind binds on the default bus to the cid 0, a non entity.
func GlobalBind(fn Bindable, name string) {
	DefaultBus.Bind(fn, name, 0)
}

// UnbindAll removes all events that match the given bindingOption from the
// default event bus
func UnbindAll(opt BindingOption) {
	DefaultBus.UnbindAll(opt)
}

// UnbindAllAndRebind is a way to reset the bindings on a CID efficiently,
// given a new set of equal length binding and event slices. This is equivalent
// to callign UnbindAll and then looping over Bind calls for the pairs of
// bindables and event names, but uses less mutex time.
func UnbindAllAndRebind(bo BindingOption, binds []Bindable, cid int, events []string) {
	DefaultBus.UnbindAllAndRebind(bo, binds, cid, events)
}

// UnbindBindable is a manual way to unbind a function Bindable. Use of
// this with closures will result in undefined behavior.
func UnbindBindable(opt UnbindOption) {
	DefaultBus.UnbindBindable(opt)
}

func Bind(fn Bindable, name string, callerID int) {
	DefaultBus.Bind(fn, name, callerID)
}

func BindPriority(fn Bindable, opt BindingOption) {
	DefaultBus.BindPriority(fn, opt)
}

func Flush() error {
	return DefaultBus.Flush()
}

func FramesElapsed() int {
	return DefaultBus.FramesElapsed()
}

func Reset() {
	DefaultBus.Reset()
}

func ResolvePending() {
	DefaultBus.ResolvePending()
}

func SetTick(framerate int) error {
	return DefaultBus.SetTick(framerate)
}

func Stop() error {
	return DefaultBus.Stop()
}

func Update() error {
	return DefaultBus.Update()
}

func UpdateLoop(framerate int, updateCh chan<- bool) error {
	return DefaultBus.UpdateLoop(framerate, updateCh)
}
