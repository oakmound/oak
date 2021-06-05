package event

// As in collision and mouse, default.go lists functions that
// only operate on DefaultBus, a package global bus.

var (
	// DefaultBus is a bus that has additional operations for CIDs, and can
	// be called via event.Call as opposed to bus.Call
	DefaultBus = NewBus(DefaultCallerMap)
)

// Trigger an event, but only for one ID, on the default bus
func (cid CID) Trigger(eventName string, data interface{}) {
	go func(eventName string, data interface{}) {
		DefaultBus.mutex.RLock()
		if idMap, ok := DefaultBus.bindingMap[eventName]; ok {
			if bs, ok := idMap[cid]; ok {
				DefaultBus.triggerDefault(bs.sl, cid, eventName, data)
			}
		}
		DefaultBus.mutex.RUnlock()
	}(eventName, data)
}

func (cid CID) TriggerBus(eventName string, data interface{}, bus Handler) chan struct{} {
	return bus.TriggerCIDBack(cid, eventName, data)
}

// Bind on a CID is shorthand for bus.Bind(name, cid, fn), on the default bus.
func (cid CID) Bind(name string, fn Bindable) {
	DefaultBus.Bind(name, cid, fn)
}

// UnbindAll removes all events with the given cid from the event bus
func (cid CID) UnbindAll() {
	DefaultBus.UnbindAll(Event{
		Name:     "",
		CallerID: cid,
	})
}

// UnbindAllAndRebind on a CID is equivalent to bus.UnbindAllAndRebind(..., cid)
func (cid CID) UnbindAllAndRebind(binds []Bindable, events []string) {
	DefaultBus.UnbindAllAndRebind(Event{
		Name:     "",
		CallerID: cid,
	}, binds, cid, events)
}

// Trigger calls Trigger on the DefaultBus
func Trigger(eventName string, data interface{}) {
	DefaultBus.Trigger(eventName, data)
}

// TriggerBack calls TriggerBack on the DefaultBus
func TriggerBack(eventName string, data interface{}) chan struct{} {
	return DefaultBus.TriggerBack(eventName, data)
}

// GlobalBind calls GlobalBind on the DefaultBus
func GlobalBind(name string, fn Bindable) {
	DefaultBus.GlobalBind(name, fn)
}

// UnbindAll calls UnbindAll on the DefaultBus
func UnbindAll(opt Event) {
	DefaultBus.UnbindAll(opt)
}

// UnbindAllAndRebind calls UnbindAllAndRebind on the DefaultBus
func UnbindAllAndRebind(bo Event, binds []Bindable, cid CID, events []string) {
	DefaultBus.UnbindAllAndRebind(bo, binds, cid, events)
}

// UnbindBindable calls UnbindBindable on the DefaultBus
func UnbindBindable(opt UnbindOption) {
	DefaultBus.UnbindBindable(opt)
}

// Bind calls Bind on the DefaultBus
func Bind(name string, callerID CID, fn Bindable) {
	DefaultBus.Bind(name, callerID, fn)
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
func UpdateLoop(framerate int, updateCh chan struct{}) error {
	return DefaultBus.UpdateLoop(framerate, updateCh)
}
