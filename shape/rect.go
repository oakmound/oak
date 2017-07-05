package shape

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

// For this type to work we'd need the ability to pick
// (and apply) a scaling algorithm
// type SelfRect [][]bool

// func (sr SelfRect) In(x, y, size) bool {

// }

// func (sr SelfRect) Rect(size) [][]bool {

// }
