package particle

import "github.com/oakmound/oak/shape"

// Shapeable generators can have the Shape option called on them
type Shapeable interface {
	SetShape(shape.Shape)
}

// Shape is an option to set a generator's shape
func Shape(sf shape.Shape) func(Generator) {
	return func(g Generator) {
		g.(Shapeable).SetShape(sf)
	}
}
