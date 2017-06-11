package shape

import "math"

// In functions return whether the given coordinate lies
// in a shape.
type In func(x, y, size int) bool

// AndIn will combine multiple In functions into one, where
// if any of the shapes are false the result is false.
func AndIn(is ...In) In {
	return func(x, y, size int) bool {
		for _, i := range is {
			b := i(x, y, size)
			if !b {
				return false
			}
		}
		return true
	}
}

// OrIn will combine multiple In functions into one, where
// if any of the shapes are true the result is true.
func OrIn(is ...In) In {
	return func(x, y, size int) bool {
		for _, i := range is {
			b := i(x, y, size)
			if b {
				return true
			}
		}
		return false
	}
}

// NotIn returns the opposite of a given In function for any query
func NotIn(i In) In {
	return func(x, y, size int) bool {
		return !i(x, y, size)
	}
}

// A JustIn lets an In function serve as a shape by automatically
// wrapping it in assistant functions for other utilites.
type JustIn In

// In acts as the underlying In function
func (ji JustIn) In(x, y, size int) bool {
	return ji(x, y, size)
}

// Rect calls InToRect on a JustIn's In
func (ji JustIn) Rect(size int) [][]bool {
	return InToRect(In(ji))(size)
}

var (
	// Square will return true for any [x][y]
	Square = JustIn(func(x, y, size int) bool {
		return true
	})
	// Diamond
	// . . t . .
	// . t t t .
	// t t t t t
	// . t t t .
	// . . t . .
	Diamond = JustIn(func(x, y, size int) bool {
		radius := size / 2
		return math.Abs(float64(x-radius))+math.Abs(float64(y-radius)) < float64(radius)
	})
	// Circle
	// . . . . . . .
	// . . t t t . .
	// . t t t t t .
	// . t t t t t .
	// . t t t t t .
	// . . t t t . .
	// . . . . . . .
	Circle = JustIn(func(x, y, size int) bool {
		radius := size / 2
		dx := math.Abs(float64(x - radius))
		dy := math.Abs(float64(y - radius))
		radiusf64 := float64(radius)
		if dx+dy <= radiusf64 {
			return true
		}
		return math.Pow(dx, 2)+math.Pow(dy, 2) < math.Pow(radiusf64, 2)
	})
	// Checkered
	// t . t . t .
	// . t . t . t
	// t . t . t .
	// . t . t . t
	// t . t . t .
	// . t . t . t
	Checkered = JustIn(func(x, y, size int) bool {
		return (x+y)%2 == 0
	})
)

// XRange is an example In utility which returns values within a given
// relative range (where 0 = 0 and 1 = size).
func XRange(a, b float64) In {
	return func(x, y, size int) bool {
		xf := float64(x)
		sf := float64(size)
		return (xf >= sf*a) && (xf <= sf*b)
	}
}
