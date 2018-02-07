package event

import "github.com/oakmound/oak/dlog"

// BindPriority is called by entities. Entities pass in a bindable function,
// and a set of options which are parsed out.
// Returns a binding that can used to unbind this binding later.
func (eb *Bus) BindPriority(fn Bindable, opt BindingOption) {
	eb.pendingMutex.Lock()
	eb.binds = append(eb.binds, UnbindOption{opt, fn})
	eb.pendingMutex.Unlock()
}

// Bind adds a function to the event bus tied to the given callerID
// to be called when the event name is triggered.
func (eb *Bus) Bind(fn Bindable, name string, callerID int) {

	bOpt := BindingOption{}
	bOpt.Event = Event{
		Name:     name,
		CallerID: callerID,
	}

	dlog.Verb("Binding ", callerID, " with name ", name)

	eb.BindPriority(fn, bOpt)
}

// GlobalBind binds on the bus to the cid 0, a non entity.
func (eb *Bus) GlobalBind(fn Bindable, name string) {
	eb.Bind(fn, name, 0)
}
