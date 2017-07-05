package intgeom

import "math/rand"

// Rect is a basic integer pair with width / height
type Rect struct {
	Min Point
	Max Point
}

func NewRect(x, y, x2, y2 int) Rect {
	return Rect{
		Min: Point{x, y},
		Max: Point{x2, y2},
	}
}

func NewRectWH(x, y, w, h int) Rect {
	return Rect{
		Min: Point{x, y},
		Max: Point{x + w, y + h},
	}
}

func ShuffleRects(rs []Rect) []Rect {
	for i := len(rs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		rs[i], rs[j] = rs[j], rs[i]
	}
	return rs
}
