package interfaces

type Point interface {
	GetX() float64
	GetY() float64
	SetPos(x, y float64)
	GetLogicPos() (float64, float64)
	SetLogicPos(x, y float64)
}

type Doodad interface {
	event.Entity
	Point
	GetRenderable() render.Renderable
	GetID() event.CID
	Destroy()
}

type Solid interface {
	Doodad
	SetDim(w, h float64)
	GetLogicDim(w, h float64)
	SetLogicDim(w, h float64)
	GetSpace() *collision.Space
	SetSpace(s *collision.Space)
}
