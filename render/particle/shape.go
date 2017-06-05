package particle

import "math"

// A ShapeFunction takes in a total size of a shape and coordinates within the
// shape, and reports whether that coordinate pair lies in the shape.
type ShapeFunction func(x, y, size int) bool

// ShapeFunction block
var (
	Square = func(x, y, size int) bool {
		return true
	}
	Diamond = func(x, y, size int) bool {
		radius := size / 2
		return math.Abs(float64(x-radius))+math.Abs(float64(y-radius)) < float64(radius)
	}
	Circle = func(x, y, size int) bool {
		radius := size / 2
		dx := math.Abs(float64(x - radius))
		dy := math.Abs(float64(y - radius))
		radiusf64 := float64(radius)
		if dx+dy <= radiusf64 {
			return true
		}
		return math.Pow(dx, 2)+math.Pow(dy, 2) < math.Pow(radiusf64, 2)
	}
	Checkered = func(x, y, size int) bool {
		return (x+y)%2 == 0
	}
)

// AndShape will combine multiple shape functions into one shape, where
// if any of the shapes are false the result is false.
func AndShape(sfs ...ShapeFunction) ShapeFunction {
	return func(x, y, size int) bool {
		for _, sf := range sfs {
			b := sf(x, y, size)
			if !b {
				return false
			}
		}
		return true
	}
}

// OrShape will combine multiple shape functions into one shape, where
// if any of the shapes are true the result is true.
func OrShape(sfs ...ShapeFunction) ShapeFunction {
	return func(x, y, size int) bool {
		for _, sf := range sfs {
			b := sf(x, y, size)
			if b {
				return true
			}
		}
		return false
	}
}

// jeez these should go into their own subpackage I guess
func NotShape(sf ShapeFunction) ShapeFunction {
	return func(x, y, size int) bool {
		return !sf(x, y, size)
	}
}

// Shapeable generators can have the Shape option called on them
type Shapeable interface {
	SetShape(ShapeFunction)
}

// Shape is an option to set a generator's shape
func Shape(sf ShapeFunction) func(Generator) {
	return func(g Generator) {
		g.(Shapeable).SetShape(sf)
	}
}
