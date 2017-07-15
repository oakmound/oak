package intgeom

import "math/rand"

// Rect is a basic integer pair with width / height
type Rect struct {
	Min Point
	Max Point
}

// NewRect returns an (X,Y):(X2,Y2) rectangle
func NewRect(x, y, x2, y2 int) Rect {
	return Rect{
		Min: Point{x, y},
		Max: Point{x2, y2},
	}
}

// NewRectWH returns an (X,Y):(X+W,Y+H) rectangle
func NewRectWH(x, y, w, h int) Rect {
	return Rect{
		Min: Point{x, y},
		Max: Point{x + w, y + h},
	}
}

// ShuffleRects is a utility function to randomize the order of a rectangle set
func ShuffleRects(rs []Rect) []Rect {
	for i := len(rs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		rs[i], rs[j] = rs[j], rs[i]
	}
	return rs
}
