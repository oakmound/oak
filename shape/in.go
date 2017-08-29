package shape

import (
	"math"

	"github.com/oakmound/oak/alg/intgeom"
)

// In functions return whether the given coordinate lies
// in a shape.
type In func(x, y int, sizes ...int) bool

// AndIn will combine multiple In functions into one, where
// if any of the shapes are false the result is false.
func AndIn(is ...In) In {
	return func(x, y int, sizes ...int) bool {
		for _, i := range is {
			b := i(x, y, sizes...)
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
	return func(x, y int, sizes ...int) bool {
		for _, i := range is {
			b := i(x, y, sizes...)
			if b {
				return true
			}
		}
		return false
	}
}

// NotIn returns the opposite of a given In function for any query
func NotIn(i In) In {
	return func(x, y int, sizes ...int) bool {
		return !i(x, y, sizes...)
	}
}

// A JustIn lets an In function serve as a shape by automatically
// wrapping it in assistant functions for other utilites.
type JustIn In

// In acts as the underlying In function
func (ji JustIn) In(x, y int, sizes ...int) bool {
	return ji(x, y, sizes...)
}

// Rect calls InToRect on a JustIn's In
func (ji JustIn) Rect(sizes ...int) [][]bool {
	return InToRect(In(ji))(sizes...)
}

// Outline calls ToOutline on the JustIn
func (ji JustIn) Outline(sizes ...int) ([]intgeom.Point, error) {
	return ToOutline(ji)(sizes...)
}

var (
	// Square will return true for any [x][y]
	Square = JustIn(func(x, y int, sizes ...int) bool {
		return true
	})

	// Rectangle will return true for any [x][y] in w, h
	Rectangle = JustIn(func(x, y int, sizes ...int) bool {
		w := sizes[0]
		h := sizes[0]
		if len(sizes) > 1 {
			h = sizes[1]
		}
		if x < w && y < h && x >= 0 && y >= 0 {
			return true
		}
		return false
	})
	// Diamond has a shape like the following:
	// . . t . .
	// . t t t .
	// t t t t t
	// . t t t .
	// . . t . .
	Diamond = JustIn(func(x, y int, sizes ...int) bool {
		radius := sizes[0] / 2
		return math.Abs(float64(x-radius))+math.Abs(float64(y-radius)) < float64(radius)
	})
	// Circle has a shape like the following:
	// . . . . . . .
	// . . t t t . .
	// . t t t t t .
	// . t t t t t .
	// . t t t t t .
	// . . t t t . .
	// . . . . . . .
	Circle = JustIn(func(x, y int, sizes ...int) bool {
		radius := sizes[0] / 2
		dx := math.Abs(float64(x - radius))
		dy := math.Abs(float64(y - radius))
		radiusf64 := float64(radius)
		if dx+dy <= radiusf64 {
			return true
		}
		return math.Pow(dx, 2)+math.Pow(dy, 2) < math.Pow(radiusf64, 2)
	})
	// Checkered has a shape like the following:
	// t . t . t .
	// . t . t . t
	// t . t . t .
	// . t . t . t
	// t . t . t .
	// . t . t . t
	Checkered = JustIn(func(x, y int, sizes ...int) bool {
		return (x+y)%2 == 0
	})
)

// XRange is an example In utility which returns values within a given
// relative range (where 0 = 0 and 1 = size).
func XRange(a, b float64) In {
	return func(x, y int, sizes ...int) bool {
		xf := float64(x)
		sf := float64(sizes[0])
		return (xf >= sf*a) && (xf <= sf*b)
	}
}
