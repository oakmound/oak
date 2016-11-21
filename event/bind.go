package event

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
)

// Called by entities.
// Entities pass in a bindable function,
// and a set of options which
// are parsed out.
// Returns a binding that can used
// to unbind this binding later.
func (eb *EventBus) BindPriority(fn Bindable, opt BindingOption, ch chan Binding) chan Binding {

	pendingMutex.Lock()

	bindablesToBind = append(bindablesToBind, fn)
	optionsToBind = append(optionsToBind, opt)
	channelsToBind = append(channelsToBind, ch)

	pendingMutex.Unlock()

	return ch
}

func GlobalBindBack(fn Bindable, name string) chan Binding {
	return thisBus.BindBack(fn, name, 0)
}

func GlobalBind(fn Bindable, name string) {
	go thisBus.Bind(fn, name, 0)
}

func (eb *EventBus) BindBack(fn Bindable, name string, callerID int) chan Binding {

	bOpt := BindingOption{}
	bOpt.Event = Event{
		Name:     name,
		CallerID: callerID,
	}

	dlog.Verb("Binding ", callerID, " with name ", name)

	ch := make(chan Binding)

	go eb.BindPriority(fn, bOpt, ch)

	return ch
}

func (eb *EventBus) Bind(fn Bindable, name string, callerID int) {

	bOpt := BindingOption{}
	bOpt.Event = Event{
		Name:     name,
		CallerID: callerID,
	}

	dlog.Verb("Binding ", callerID, " with name ", name)

	go eb.BindPriority(fn, bOpt, nil)
}

func (cid CID) BindBack(fn Bindable, name string) chan Binding {
	return thisBus.BindBack(fn, name, int(cid))
}

func (cid CID) Bind(fn Bindable, name string) {
	go thisBus.Bind(fn, name, int(cid))
}
