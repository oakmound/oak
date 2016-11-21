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
	eb := GetEventBus()
	bo := BindingOption{
		Event{
			"",
			int(*cid),
		},
		0,
	}
	eb.UnbindAll(bo)
}

// Called by entities or by game logic.
// Unbinds all events in this bus which
// match the given binding options.
func (eb *EventBus) UnbindAll(opt BindingOption) {
	pendingMutex.Lock()
	optionsToUnbind = append(optionsToUnbind, opt)
	pendingMutex.Unlock()
}
