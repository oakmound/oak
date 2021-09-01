package event

// Bind adds a function to the event bus tied to the given callerID
// to be called when the event name is triggered. It is equivalent to
// calling BindPriority with a zero Priority.
func (eb *Bus) Bind(name string, callerID CID, fn Bindable) {
	eb.pendingMutex.Lock()
	eb.binds = append(eb.binds, UnbindOption{
		Event: Event{
			Name:     name,
			CallerID: callerID,
		}, Fn: fn})
	eb.pendingMutex.Unlock()
}

// GlobalBind binds on the bus to the cid 0, a non entity.
func (eb *Bus) GlobalBind(name string, fn Bindable) {
	eb.Bind(name, 0, fn)
}

// Empty is a helper to convert a func() into a Bindable function signature.
func Empty(f func()) Bindable {
	return func(CID, interface{}) int {
		f()
		return 0
	}
}
