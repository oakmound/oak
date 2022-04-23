package render

import (
	"github.com/oakmound/oak/v4/physics"
)

const (
	// Undraw is a layer representing elements that should be undrawn, or removed
	// from the draw stack. This is exported in the rare case that there is
	// a need to use the default value for something else.
	Undraw = -1000
)

// Layered types know the order they should be drawn in relative to
// other layered types. Higher layers are drawn after lower layers,
// and so will appear on top of them. Drawn layers are anticipated to
// be all positive.
//
// Basic Implementing struct: Layer
type Layered interface {
	GetLayer() int
	SetLayer(l int)
	Undraw()
}

// A Layer object has a draw layer
type Layer struct {
	layer int
}

// GetLayer returns the layer of an object if it has one or else returns that the object needs to be undrawn
func (ld *Layer) GetLayer() int {
	if ld == nil {
		return Undraw
	}
	return ld.layer
}

// SetLayer sets an object's layer
func (ld *Layer) SetLayer(l int) {
	ld.layer = l
}

// Undraw sets that a Layer object should be undrawn
func (ld *Layer) Undraw() {
	ld.layer = Undraw
}

// A LayeredPoint is an object with a position Vector and a layer
type LayeredPoint struct {
	physics.Vector
	Layer
}

// NewLayeredPoint creates a new LayeredPoint at a given location and layer
func NewLayeredPoint(x, y float64, l int) LayeredPoint {
	return LayeredPoint{
		Vector: physics.NewVector(x, y),
		Layer:  Layer{l},
	}
}

// GetLayer returns the layer of this point. If this is nil,
// it will return Undraw
func (ldp *LayeredPoint) GetLayer() int {
	if ldp == nil {
		return Undraw
	}
	return ldp.Layer.GetLayer()
}

// Copy deep copies the LayeredPoint
func (ldp *LayeredPoint) Copy() LayeredPoint {
	ldp2 := LayeredPoint{}
	ldp2.Vector = ldp.Vector.Copy()
	ldp2.Layer = ldp.Layer
	return ldp2
}

// These functions are redefined because vector's internal
// functions return Vectors, and we don't want to return Vectors.

// ShiftX moves the LayeredPoint by the given x
func (ldp *LayeredPoint) ShiftX(x float64) {
	ldp.Vector.ShiftX(x)
}

// ShiftY moves the LayeredPoint by the given y
func (ldp *LayeredPoint) ShiftY(y float64) {
	ldp.Vector.ShiftY(y)
}

// SetPos sets the LayeredPoint's position to the given x, y
func (ldp *LayeredPoint) SetPos(x, y float64) {
	ldp.Vector.SetPos(x, y)
}

// GetDims returns the dimensions of this object. As a single point, LayeredPoint
// returns (1,1).
func (ldp *LayeredPoint) GetDims() (int, int) {
	return 1, 1
}
