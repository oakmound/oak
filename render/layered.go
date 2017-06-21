package render

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/physics"
)

const (
	Undraw = -1000
)

type Layered struct {
	layer int
}

func (ld *Layered) GetLayer() int {
	if ld == nil {
		return Undraw
	}
	return ld.layer
}

func (ld *Layered) SetLayer(l int) {
	ld.layer = l
}

func (ld *Layered) UnDraw() {
	ld.layer = Undraw
}

type LayeredPoint struct {
	physics.Vector
	Layered
}

func NewLayeredPoint(x, y float64, l int) LayeredPoint {
	return LayeredPoint{
		Vector:  physics.NewVector(x, y),
		Layered: Layered{l},
	}
}

func (ldp *LayeredPoint) GetLayer() int {
	if ldp == nil {
		return Undraw
	}
	return ldp.Layered.GetLayer()
}

func (ldp *LayeredPoint) Copy() LayeredPoint {
	ldp2 := LayeredPoint{}
	ldp2.Vector = ldp.Vector.Copy()
	ldp2.Layered = ldp.Layered
	return ldp2
}

func (ldp *LayeredPoint) ShiftX(x float64) {
	ldp.Vector.ShiftX(x)
}
func (ldp *LayeredPoint) ShiftY(y float64) {
	ldp.Vector.ShiftY(y)
}

func (ldp *LayeredPoint) SetPos(x, y float64) {
	ldp.Vector.SetPos(x, y)
}

func (ldp *LayeredPoint) GetDims() (int, int) {
	// We use 6,6 here because our polygon containment library has a bug where it
	// will misreport artificially small dimensions. This function is expected to
	// only be used to determine whether something is onscreen to be drawn.
	// Todo: write own polygon containment library
	return 6, 6
}

func (ldp *LayeredPoint) String() string {
	x := strconv.FormatFloat(ldp.X(), 'f', 2, 32)
	y := strconv.FormatFloat(ldp.Y(), 'f', 2, 32)
	l := strconv.Itoa(ldp.layer)
	return "X(): " + x + ", Y(): " + y + ", L: " + l
}
