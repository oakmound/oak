package collision

const (
	// NilLabel is used internally for spaces that are otherwise not
	// given labels.
	// TODO V3: why is this exported
	NilLabel Label = -1
)

// Label is used to store type information for a given space
type Label int
