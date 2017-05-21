package event

// Unbind on a binding is a rewriting of bus.Unbind(b)
func (b Binding) Unbind() {
	thisBus.Unbind(b)
}

// UnbindAll removes all events with the given cid from the event bus
func (cid *CID) UnbindAll() {
	bo := BindingOption{
		Event{
			"",
			int(*cid),
		},
		0,
	}
	UnbindAll(bo)
}

// UnbindAll removes all events that match the given bindingOption from the
// event bus
func UnbindAll(opt BindingOption) {
	pendingMutex.Lock()
	partUnbinds = append(partUnbinds, opt)
	pendingMutex.Unlock()
}

// UnbindAllAndRebind on a CID is equivalent to bus.UnbindAllAndRebind(..., cid)
func (cid *CID) UnbindAllAndRebind(binds []Bindable, events []string) {
	bo := BindingOption{
		Event{
			"",
			int(*cid),
		},
		0,
	}
	UnbindAllAndRebind(bo, binds, int(*cid), events)
}

// UnbindAllAndRebind is a way to reset the bindings on a CID efficiently,
// given a new set of equal length binding and event slices. This is equivalent
// to callign UnbindAll and then looping over Bind calls for the pairs of
// bindables and event names, but uses less mutex time.
func UnbindAllAndRebind(bo BindingOption, binds []Bindable, cid int, events []string) {
	opts := make([]BindingOption, len(events))
	for k, v := range events {
		opts[k].Event = Event{
			Name:     v,
			CallerID: cid,
		}
	}

	pendingMutex.Lock()
	unbindAllAndRebinds = append(unbindAllAndRebinds, UnbindAllOption{
		ub:   bo,
		bs:   opts,
		bnds: binds,
	})
	pendingMutex.Unlock()
}

// Unbind removes the given binding (previously returned by .Bind) from
// the event bus, but because synchronously expecting something to
// be returned from .Bind is dangerous, that behavior was removed and it is now
// just used by the engine internals
func (eb *Bus) Unbind(b Binding) {
	pendingMutex.Lock()
	unbinds = append(unbinds, b)
	pendingMutex.Unlock()
}

// UnbindBindable is used by UNBIND_EVENT calls
func UnbindBindable(opt UnbindOption) {
	pendingMutex.Lock()
	fullUnbinds = append(fullUnbinds, opt)
	pendingMutex.Unlock()
}
