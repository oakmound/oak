package render

import "github.com/oakmound/oak/event"

// NonStatic types are not always static. If something is not NonStatic,
// it is equivalent to having IsStatic always return true.
type NonStatic interface {
	IsStatic() bool
}

// Triggerable types can have an ID set so when their animations finish,
// they trigger AnimationEnd on that ID.
type Triggerable interface {
	SetTriggerID(event.CID)
}

type updates interface {
	update()
}
