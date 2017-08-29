package render

// NonInterruptable types are not always interruptable.  If something is not
// NonInterruptable, it is equivalent to having IsInterruptable always return
// true.
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
