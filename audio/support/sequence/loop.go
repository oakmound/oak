package sequence

type Loop bool

type HasLoops interface {
	GetLoop() bool
	SetLoop(bool)
}

func (l *Loop) GetLoop() bool {
	return bool(*l)
}

func (l *Loop) SetLoop(b bool) {
	*l = Loop(b)
}

// Loops sets the generator's Loop
func Loops(b bool) Option {
	return func(g Generator) {
		if ht, ok := g.(HasLoops); ok {
			ht.SetLoop(b)
		}
	}
}
