package render

import (
	"strconv"

	"fmt"

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

func NewLayeredPoint(x, y float64, l int) LayeredPoint {
	return LayeredPoint{
		Vector:  physics.NewVector(x, y),
		Layered: Layered{l},
	}
}

func (ldp *LayeredPoint) Copy() LayeredPoint {
	ldp2 := LayeredPoint{}
	ldp2.Vector = ldp.Vector.Copy()
	ldp2.Layered = ldp.Layered
	return ldp2
}

func (ldp *LayeredPoint) ShiftX(x float64) {
	fmt.Println(x)
	fmt.Println(ldp.Vector)
	ldp.Vector = ldp.Vector.ShiftX(x)
}
func (ldp *LayeredPoint) ShiftY(y float64) {
	ldp.Vector = ldp.Vector.ShiftY(y)
}

func (ldp *LayeredPoint) SetPos(x, y float64) {
	ldp.Vector = ldp.Vector.SetPos(x, y)
}

func (ldp *LayeredPoint) GetDims() (int, int) {
	return 6, 6
}

func (ldp *LayeredPoint) String() string {
	x := strconv.FormatFloat(ldp.X(), 'f', 2, 32)
	y := strconv.FormatFloat(ldp.Y(), 'f', 2, 32)
	l := strconv.Itoa(ldp.layer)
	return "X(): " + x + ", Y(): " + y + ", L: " + l
}
