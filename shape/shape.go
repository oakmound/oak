package shape

// A Shape represents a rectangle of width/height size
// where for each x,y coordinate, either that value lies
// inside the shape or outside of the shape, represented
// by true or false. Shapes can be fuzzed along their border
// to create gradients of floats, and shapes can be queried
// to just produce a 2d boolean array of width/height size.
// Todo: consider if the number of coordinate arguments
// should be variadic, if width/height should not be combined
// and/or variadic, for additional dimension support
type Shape interface {
	In(x, y int, sizes ...int) bool
	Outline(sizes ...int) ([]Point, error)
	Rect(sizes ...int) [][]bool
}
