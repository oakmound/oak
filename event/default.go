package event

// As in collision and mouse, default.go lists functions that
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
					lst := bs.highPriority[i]
					if lst != nil {
						DefaultBus.triggerDefault((*lst).sl, iid, eventName, data)
					}
				}
				DefaultBus.triggerDefault((bs.defaultPriority).sl, iid, eventName, data)

				for i := 0; i < bs.lowIndex; i++ {
					lst := bs.lowPriority[i]
					if lst != nil {
						DefaultBus.triggerDefault((*lst).sl, iid, eventName, data)
					}
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

// Trigger calls Trigger on the DefaultBus
func Trigger(eventName string, data interface{}) {
	DefaultBus.Trigger(eventName, data)
}

// TriggerBack calls TriggerBack on the DefaultBus
func TriggerBack(eventName string, data interface{}) chan bool {
	return DefaultBus.TriggerBack(eventName, data)
}

// GlobalBind calls GlobalBind on the DefaultBus
func GlobalBind(fn Bindable, name string) {
	DefaultBus.GlobalBind(fn, name)
}

// UnbindAll calls UnbindAll on the DefaultBus
func UnbindAll(opt BindingOption) {
	DefaultBus.UnbindAll(opt)
}

// UnbindAllAndRebind calls UnbindAllAndRebind on the DefaultBus
func UnbindAllAndRebind(bo BindingOption, binds []Bindable, cid int, events []string) {
	DefaultBus.UnbindAllAndRebind(bo, binds, cid, events)
}

// UnbindBindable calls UnbindBindable on the DefaultBus
func UnbindBindable(opt UnbindOption) {
	DefaultBus.UnbindBindable(opt)
}

// Bind calls Bind on the DefaultBus
func Bind(fn Bindable, name string, callerID int) {
	DefaultBus.Bind(fn, name, callerID)
}

// BindPriority calls BindPriority on the DefaultBus
func BindPriority(fn Bindable, opt BindingOption) {
	DefaultBus.BindPriority(fn, opt)
}

// Flush calls Flush on the DefaultBus
func Flush() error {
	return DefaultBus.Flush()
}

// FramesElapsed calls FramesElapsed on the DefaultBus
func FramesElapsed() int {
	return DefaultBus.FramesElapsed()
}

// Reset calls Reset on the DefaultBus
func Reset() {
	DefaultBus.Reset()
}

// ResolvePending calls ResolvePending on the DefaultBus
func ResolvePending() {
	DefaultBus.ResolvePending()
}

// SetTick calls SetTick on the DefaultBus
func SetTick(framerate int) error {
	return DefaultBus.SetTick(framerate)
}

// Stop calls Stop on the DefaultBus
func Stop() error {
	return DefaultBus.Stop()
}

// Update calls Update on the DefaultBus
func Update() error {
	return DefaultBus.Update()
}

// UpdateLoop calls UpdateLoop on the DefaultBus
func UpdateLoop(framerate int, updateCh chan bool) error {
	return DefaultBus.UpdateLoop(framerate, updateCh)
}
