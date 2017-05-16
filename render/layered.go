package render

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/physics"
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
	physics.Vector
	Layered
}

func (ldp *LayeredPoint) ShiftX(x float64) {
	ldp.X += x
}
func (ldp *LayeredPoint) ShiftY(y float64) {
	ldp.Y += y
}

func (ldp *LayeredPoint) SetPos(x, y float64) {
	ldp.X = x
	ldp.Y = y
}

func (ldp *LayeredPoint) GetDims() (int, int) {
	return 6, 6
}

func (ldp *LayeredPoint) String() string {
	x := strconv.FormatFloat(ldp.X, 'f', 2, 32)
	y := strconv.FormatFloat(ldp.Y, 'f', 2, 32)
	l := strconv.Itoa(ldp.layer)
	return "X: " + x + ", Y: " + y + ", L: " + l
}
