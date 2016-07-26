package render

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
