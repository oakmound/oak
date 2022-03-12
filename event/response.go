package event

type Response uint8

// Response types for bindables
const (
	// NoResponse or 0, is returned by events that
	// don't want the event bus to do anything with
	// the event after they have been evaluated. This
	// is the usual behavior.
	NoResponse Response = iota
	// UnbindThis unbinds the one binding that returns it.
	UnbindThis
)
