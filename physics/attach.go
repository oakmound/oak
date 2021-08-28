package physics

// An Attachable can be attached to static or moving vectors.
type Attachable interface {
	Detach()
	Attach(Vecer, float64, float64)
	AttachX(Vecer, float64)
	AttachY(Vecer, float64)
	Vecer
}

// A Vecer can be converted into a Vector
type Vecer interface {
	Vec() Vector
}

// Vec returns a vector itself
func (v Vector) Vec() Vector {
	return v
}

// Attach takes in something for this vector to attach to and a set of
// offsets.
func (v *Vector) Attach(a Vecer, offX, offY float64) {
	v2 := a.Vec()
	v.x = v2.x
	v.y = v2.y
	v.offX = offX
	v.offY = offY
}

// AttachX performs an attachment that only attaches on the X axis.
func (v *Vector) AttachX(a Vecer, offX float64) {
	v2 := a.Vec()
	v.x = v2.x
	v.offX = offX
}

// AttachY performs an attachment that only attaches on the Y axis.
func (v *Vector) AttachY(a Vecer, offY float64) {
	v2 := a.Vec()
	v.y = v2.y
	v.offY = offY
}

// Detach modified a vector to no longer be attached to anything.
func (v *Vector) Detach() {
	v2 := NewVector(v.X(), v.Y())
	*v = v2
}
