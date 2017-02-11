package render

import (
	"strconv"
)

type Layered struct {
	layer int
}

func (ld *Layered) GetLayer() int {
	return ld.layer
}

func (ld *Layered) SetLayer(l int) {
	ld.layer = l
}

func (ld *Layered) UnDraw() {
	ld.layer = -1
}

type LayeredPoint struct {
	Point
	Layered
}


func (ldp *LayeredPoint) String() string {
	x := strconv.FormatFloat(ldp.X, 'f', 2, 32)
	y := strconv.FormatFloat(ldp.Y, 'f', 2, 32)
	l := strconv.Itoa(ldp.layer)
	return "X: " + x + ", Y: " + y + ", L: " + l
}
