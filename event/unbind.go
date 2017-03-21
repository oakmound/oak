package event

// Called by entities,
// for unbinding specific bindings.
func (eb *EventBus) Unbind(b Binding) {
	pendingMutex.Lock()
	bindingsToUnbind = append(bindingsToUnbind, b)
	pendingMutex.Unlock()
}

func (b Binding) Unbind() {
	thisBus.Unbind(b)
}

// Unbind all events for
// the given CID
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

func UnbindAllAndRebind(bo BindingOption, binds []Bindable, cid int, events []string) {
	opt := make([]BindingOption, len(events))
	for k, v := range events {
		opt[k].Event = Event{
			Name:     v,
			CallerID: cid,
		}
	}

	pendingMutex.Lock()
	orderedUnbinds = append(orderedUnbinds, bo)
	orderedBindOptions = append(orderedBindOptions, opt...)
	orderedBindables = append(orderedBindables, binds...)
	pendingMutex.Unlock()
}

// Called by entities or by game logic.
// Unbinds all events in this bus which
// match the given binding options.
func UnbindAll(opt BindingOption) {
	pendingMutex.Lock()
	optionsToUnbind = append(optionsToUnbind, opt)
	pendingMutex.Unlock()
}

func UnbindBindable(opt UnbindOption) {
	pendingMutex.Lock()
	ubOptionsToUnbind = append(ubOptionsToUnbind, opt)
	pendingMutex.Unlock()
}
