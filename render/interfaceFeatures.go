package render

// CanPause types have pause functions to start and stop animation
type CanPause interface {
	Pause()
	Unpause()
}

// NonStatic types are not always static. If something is not NonStatic,
// it is equivalent to having IsStatic always return true.
type NonStatic interface {
	IsStatic() bool
}

// NonInterruptable types are not always interruptable.  If something is not
// NonInterruptable, it is equivalent to having IsInterruptable always return
// true.
type NonInterruptable interface {
	IsInterruptable() bool
}

type updates interface {
	update()
}
