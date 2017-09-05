package shape

import (
	"github.com/oakmound/oak/alg/intgeom"
)

// A Rect is a function that returns a 2d boolean array
// of booleans for a given size, where true represents
// that the bounded shape contains the point [x][y].
type Rect func(sizes ...int) [][]bool

// InToRect converts an In function into a Rect function.
// Know that, if you are planning on looping over this only
// once, it's better to just use the In function. The use
// case for this is if the same size rect will be queried
// on some function multiple times, and just having the booleans
// to re-access is needed.
func InToRect(i In) Rect {
	return func(sizes ...int) [][]bool {
		w := sizes[0]
		h := sizes[0]
		if len(sizes) > 1 {
			h = sizes[1]
		}
		out := make([][]bool, w)
		for x := range out {
			out[x] = make([]bool, h)
			for y := range out[x] {
				out[x][y] = i(x, y, sizes...)
			}
		}
		return out
	}
}

// A StrictRect is a shape that ignores input width and height given to it.
type StrictRect [][]bool

// NewStrictRect returns a StrictRect with the given strict dimensions, all
// values set fo false.
func NewStrictRect(w, h int) StrictRect {
	sh := make(StrictRect, w)
	for x := range sh {
		sh[x] = make([]bool, h)
	}
	return sh
}

// In returns whether the input x and y are within this StrictRect's shape.
// If the shape is undefined for the input values, it returns false.
func (sr StrictRect) In(x, y int, sizes ...int) bool {
	if x > len(sr) {
		return false
	}
	if y > len(sr[x]) {
		return false
	}
	return sr[x][y]
}

// Outline returns this StrictRect's outline, ignoring the input dimensions.
func (sr StrictRect) Outline(sizes ...int) ([]intgeom.Point, error) {
	return ToOutline(sr)(len(sr), len(sr[0]))
}

// Rect returns the StrictRect itself.
func (sr StrictRect) Rect(sizes ...int) [][]bool {
	return sr
}
