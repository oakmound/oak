package event

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
)

// BindPriority is called by entities. Entities pass in a bindable function,
// and a set of options which are parsed out.
// Returns a binding that can used to unbind this binding later.
func (eb *EventBus) BindPriority(fn Bindable, opt BindingOption) {
	pendingMutex.Lock()
	binds = append(binds, UnbindOption{opt, fn})
	pendingMutex.Unlock()
}

func GlobalBind(fn Bindable, name string) {
	thisBus.Bind(fn, name, 0)
}

func (eb *EventBus) Bind(fn Bindable, name string, callerID int) {

	bOpt := BindingOption{}
	bOpt.Event = Event{
		Name:     name,
		CallerID: callerID,
	}

	dlog.Verb("Binding ", callerID, " with name ", name)

	eb.BindPriority(fn, bOpt)
}

func (cid CID) Bind(fn Bindable, name string) {
	thisBus.Bind(fn, name, int(cid))
}
