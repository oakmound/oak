package event

// Unbind on a binding is a rewriting of bus.Unbind(b)
func (b binding) unbind(eb *Bus) {
	eb.unbind(b)
}

func (eb *Bus) unbind(b binding) {
	eb.pendingMutex.Lock()
	eb.unbinds = append(eb.unbinds, b)
	eb.pendingMutex.Unlock()
}

// UnbindAllAndRebind is a way to reset the bindings on a CID efficiently,
// given a new set of equal length binding and event slices. This is equivalent
// to callign UnbindAll and then looping over Bind calls for the pairs of
// bindables and event names, but uses less mutex time.
func (eb *Bus) UnbindAllAndRebind(bo BindingOption, binds []Bindable, cid int, events []string) {
	opts := make([]BindingOption, len(events))
	for k, v := range events {
		opts[k].Event = Event{
			Name:     v,
			CallerID: cid,
		}
	}

	eb.pendingMutex.Lock()
	eb.unbindAllAndRebinds = append(eb.unbindAllAndRebinds, UnbindAllOption{
		ub:   bo,
		bs:   opts,
		bnds: binds,
	})
	eb.pendingMutex.Unlock()
}

// UnbindAll removes all events that match the given bindingOption from the
// default event bus
func (eb *Bus) UnbindAll(opt BindingOption) {
	eb.pendingMutex.Lock()
	eb.partUnbinds = append(eb.partUnbinds, opt)
	eb.pendingMutex.Unlock()
}

// UnbindBindable is a manual way to unbind a function Bindable. Use of
// this with closures will result in undefined behavior.
func (eb *Bus) UnbindBindable(opt UnbindOption) {
	eb.pendingMutex.Lock()
	eb.fullUnbinds = append(eb.fullUnbinds, opt)
	eb.pendingMutex.Unlock()
}
