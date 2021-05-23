package event

// Response types from bindables
// reponses are not their own type because func(event.CID, interface{}) int
// is easier to write than func(event.CID, interface{}) event.Response. This may
// yet change.
const (
	// NoResponse or 0, is returned by events that
	// don't want the event bus to do anything with
	// the event after they have been evaluated. This
	// is the usual behavior.
	NoResponse = iota
	// UnbindEvent unbinds everything for a specific
	// event name from an entity at the bindable's
	// priority.
	UnbindEvent
	// UnbindSingle just unbinds the one binding that
	// it is returned from
	UnbindSingle
)
