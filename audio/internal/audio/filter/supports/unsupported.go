package supports

// Unsupported is an error type reporting that a filter was not supported
// by the Audio type it was used on
type Unsupported struct {
	filters []string
}

// NewUnsupported returns an Unsupported error with the input filters
func NewUnsupported(filters []string) Unsupported {
	return Unsupported{filters}
}

func (un Unsupported) Error() string {
	s := "Unsupported filters: "
	for _, f := range un.filters {
		s += f + " "
	}
	return s
}
