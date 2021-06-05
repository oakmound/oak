package sequence

// A Generator stores settings to create a sequence
type Generator interface {
	Generate() *Sequence
}

// Option types are inserted into Constructors to create generators
type Option func(Generator)

// And combines any number of options into a single option.
// And is a reminder that you can store combined settings to avoid
// having to rewrite them
func And(opts ...Option) Option {
	return func(g Generator) {
		for _, opt := range opts {
			opt(g)
		}
	}
}
