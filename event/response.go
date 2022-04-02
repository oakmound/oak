package event

type Response uint8

// Response types for bindables
const (
	// ResponseNone or 0, is returned by events that
	// don't want the event bus to do anything with
	// the event after they have been evaluated. This
	// is the usual behavior.
	ResponseNone Response = iota
	// ResponseUnbindThisBinding unbinds the one binding that returns it.
	ResponseUnbindThisBinding
	// ResponseUnbindThisCaller unbinds all of a caller's bindings when returned from any binding.
	ResponseUnbindThisCaller
)
