package physics

type Attachable interface {
	Detach() Vector
	Attach(Attachable) Vector
	AttachX(Attachable) Vector
	AttachY(Attachable) Vector
	Vec() Vector
}

func (v Vector) Vec() Vector {
	return v
}

func (v Vector) Attach(a Attachable) Vector {
	v2 := a.Vec()
	v3 := Vector{v2.x, v2.y}
	return v3
}

func (v Vector) AttachX(a Attachable) Vector {
	v2 := a.Vec()
	v3 := Vector{v2.x, v.y}
	return v3
}

func (v Vector) AttachY(a Attachable) Vector {
	v2 := a.Vec()
	v3 := Vector{v.x, v2.y}
	return v3
}

func (v Vector) Detach() Vector {
	return NewVector(v.X(), v.Y())
}
