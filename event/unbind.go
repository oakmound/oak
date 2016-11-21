package event

import (
	"fmt"
)

// Called by entities,
// for unbinding specific bindings.
func (eb *EventBus) Unbind(b Binding) {
	fmt.Println("Locking Pending Mutex")
	pendingMutex.Lock()
	bindingsToUnbind = append(bindingsToUnbind, b)
	fmt.Println("Added a binding to unbind")
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
