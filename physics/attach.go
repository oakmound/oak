package physics

type Attachable interface {
	Detach() Vector
	Attach(Attachable, ...float64) Vector
	AttachX(Attachable, float64) Vector
	AttachY(Attachable, float64) Vector
	Vec() Vector
}

func (v Vector) Vec() Vector {
	return v
}

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

func (v Vector) AttachX(a Attachable, offX float64) Vector {
	v2 := a.Vec()
	v3 := Vector{v2.x, v.y, offX, 0}
	return v3
}

func (v Vector) AttachY(a Attachable, offY float64) Vector {
	v2 := a.Vec()
	v3 := Vector{v.x, v2.y, 0, offY}
	return v3
}

func (v Vector) Detach() Vector {
	return NewVector(v.X(), v.Y())
}
