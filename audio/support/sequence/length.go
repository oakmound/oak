package sequence

type Length int

type HasLength interface {
	GetLength() int
	SetLength(int)
}

func (l *Length) GetLength() int {
	return int(*l)
}

func (l *Length) SetLength(i int) {
	*l = Length(i)
}

func PlayLength(i int) Option {
	return func(g Generator) {
		if l, ok := g.(HasLength); ok {
			l.SetLength(i)
		}
	}
}
