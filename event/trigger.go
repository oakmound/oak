package event

// Trigger will scan through the event bus and call all bindables found attached
// to the given event, with the passed in data.
func (bus *Bus) TriggerForCaller(callerID CallerID, eventID UnsafeEventID, data interface{}) chan struct{} {
	if callerID == Global {
		return bus.Trigger(eventID, data)
	}
	ch := make(chan struct{})
	go func() {
		bus.mutex.RLock()
		if idMap, ok := bus.bindingMap[eventID]; ok {
			if bs, ok := idMap[callerID]; ok {
				bus.trigger(bs, eventID, callerID, data)
			}
		}
		bus.mutex.RUnlock()
		close(ch)
	}()
	return ch
}

func (bus *Bus) Trigger(eventID UnsafeEventID, data interface{}) chan struct{} {
	ch := make(chan struct{})
	go func() {
		bus.mutex.RLock()
		for callerID, bs := range bus.bindingMap[eventID] {
			bus.trigger(bs, eventID, callerID, data)
		}
		bus.mutex.RUnlock()
		close(ch)
	}()
	return ch
}

func TriggerOn[T any](b Handler, ev EventID[T], data T) chan struct{} {
	return b.Trigger(ev.UnsafeEventID, data)
}

func TriggerForCallerOn[T any](b Handler, cid CallerID, ev EventID[T], data T) chan struct{} {
	return b.TriggerForCaller(cid, ev.UnsafeEventID, data)
}
