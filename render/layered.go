package render

import (
	"strconv"

	"github.com/oakmound/oak/physics"
)

const (
	// Undraw is a constant used to undraw elements
	Undraw = -1000
)

//A Layered object is one with a layer
type Layered struct {
	layer int
}

//GetLayer returns the layer of an object if it has one or else returns that the object needs to be undrawn
func (ld *Layered) GetLayer() int {
	if ld == nil {
		return Undraw
	}
	return ld.layer
}

//SetLayer sets an object that has a layer to the given layer
func (ld *Layered) SetLayer(l int) {
	ld.layer = l
}

//UnDraw sets that a layered object should be undrawn
func (ld *Layered) UnDraw() {
	ld.layer = Undraw
}

//A LayeredPoint is an object with a position Vector and a layer
type LayeredPoint struct {
	physics.Vector
	Layered
}

//NewLayeredPoint creates a new LayeredPoint at a given location and layer
func NewLayeredPoint(x, y float64, l int) LayeredPoint {
	return LayeredPoint{
		Vector:  physics.NewVector(x, y),
		Layered: Layered{l},
	}
}

// GetLayer returns the layer of this point. If this is nil,
// it will return Undraw
func (ldp *LayeredPoint) GetLayer() int {
	if ldp == nil {
		return Undraw
	}
	return ldp.Layered.GetLayer()
}

// Copy deep copies the LayeredPoint
func (ldp *LayeredPoint) Copy() LayeredPoint {
	ldp2 := LayeredPoint{}
	ldp2.Vector = ldp.Vector.Copy()
	ldp2.Layered = ldp.Layered
	return ldp2
}

// These functions are redefined because vector's internal
// functions return Vectors, and we don't want to return Vectors.

// ShiftX moves the LayeredPoint by the given x
func (ldp *LayeredPoint) ShiftX(x float64) {
	ldp.Vector.ShiftX(x)
}

//ShiftY moves the LayeredPoint by the given y
func (ldp *LayeredPoint) ShiftY(y float64) {
	ldp.Vector.ShiftY(y)
}

//SetPos sets the LayeredPoint's position to the given x, y
func (ldp *LayeredPoint) SetPos(x, y float64) {
	ldp.Vector.SetPos(x, y)
}

//GetDims returns a static small amount so that polygon containment does not throw errors
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
