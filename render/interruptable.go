package render

// NonInterruptable types are not always interruptable.  If something is not
// NonInterruptable, it is equivalent to having IsInterruptable always return
// true.
//
// The intended use of the Interruptable datatypes is entirely external--
// oak does not use them internally. The use case is for an entity that has
// a set of potential animations, and attempts to switch from one animation
// to another. The Interuptable boolean should represent whether that
// animation should be able to be switched out of before it ends.
//
// Because this use case is minor, this is a candidate for removal from render
// and moving into an auxillary package.
//
// Unless otherwie noted, all NonInterruptable types are interruptable when
// they are initialized and need to be switched (if the type supports it) to
// be non interruptbable.
type NonInterruptable interface {
	IsInterruptable() bool
}

// InterruptBool is a composable struct for NonInterruptable support
type InterruptBool struct {
	Interruptable bool
}

// IsInterruptable returns whether this can be interrupted.
func (ib InterruptBool) IsInterruptable() bool {
	return ib.Interruptable
}
