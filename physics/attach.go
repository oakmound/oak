package physics

// Attachable represents things that can be attached to Vectors
type Attachable interface {
	Detach() Vector
	Attach(Attachable, ...float64) Vector
	AttachX(Attachable, float64) Vector
	AttachY(Attachable, float64) Vector
	Vec() Vector
}

// Vec returns a vector itself
func (v Vector) Vec() Vector {
	return v
}

// Attach takes in something for this vector to attach to and a set of
// offsets. The resulting combined vector with the offsets is then returned,
// and needs to be assigned to the calling vector.
func (v Vector) Attach(a Attachable, offsets ...float64) Vector {
	xOff := 0.0
	yOff := 0.0
	if len(offsets) > 0 {
		xOff = offsets[0]
	}
	if len(offsets) > 1 {
		yOff = offsets[1]
	}
	v2 := a.Vec()
	v3 := Vector{v2.x, v2.y, xOff, yOff}
	return v3
}

// AttachX performs an attachment that only attaches on the X axis.
func (v Vector) AttachX(a Attachable, offX float64) Vector {
	v2 := a.Vec()
	v3 := Vector{v2.x, v.y, offX, 0}
	return v3
}

// AttachY performs an attachment that only attaches on the Y axis.
func (v Vector) AttachY(a Attachable, offY float64) Vector {
	v2 := a.Vec()
	v3 := Vector{v.x, v2.y, 0, offY}
	return v3
}

// Detach returns a vector no longer attached to anything. The returned vector
// needs to be assigned to the caller for the caller to be replaced (vectors
// do not use pointer receivers)
func (v Vector) Detach() Vector {
	return NewVector(v.X(), v.Y())
}
