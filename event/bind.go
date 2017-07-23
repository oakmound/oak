package event

import "github.com/oakmound/oak/dlog"

// BindPriority is called by entities. Entities pass in a bindable function,
// and a set of options which are parsed out.
// Returns a binding that can used to unbind this binding later.
func (eb *Bus) BindPriority(fn Bindable, opt BindingOption) {
	pendingMutex.Lock()
	binds = append(binds, UnbindOption{opt, fn})
	pendingMutex.Unlock()
}

// GlobalBind binds to the cid 0, a non entity.
func GlobalBind(fn Bindable, name string) {
	thisBus.Bind(fn, name, 0)
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

// Bind on a CID is shorthand for bus.Bind(fn, name, cid)
func (cid CID) Bind(fn Bindable, name string) {
	thisBus.Bind(fn, name, int(cid))
}

// BindPriority on a CID is shorthand for bus.BindPriority(fn, ...)
func (cid CID) BindPriority(fn Bindable, name string, priority int) {
	thisBus.BindPriority(fn, BindingOption{
		Event{
			name,
			int(cid),
		},
		priority,
	})
}
