package render

import "image/draw"

// Functional is a renderable which is defined by a function.
// Useful for more complex behaviors.
// the function F should not have any side effects,
// behavior when it does is undefined.
type Functional struct {
	F func() Renderable
}

func (f Functional) DrawOffset(buff draw.Image, x, y float64) {
	f.F().DrawOffset(buff, x, y)
}

func (f Functional) Draw(buff draw.Image) {
	f.F().Draw(buff)
}

func (f Functional) GetDims() (int,int) {
	return f.F().GetDims()
}

func (f Functional) GetLayer() int {
	return f.F().GetLayer()
}

func (f Functional) SetLayer(l int) {
	f.F().SetLayer(l)
}

func (f Functional) Undraw() {
	f.F().Undraw()
}

func (f Functional) X() float64 {
	return f.F().X()
}

func (f Functional) Y() float64 {
	return f.F().Y()
}

func (f Functional) ShiftX(x float64) {
	f.F().ShiftX(x)
}

func (f Functional) ShiftY(y float64) {
	f.F().ShiftY(y)
}

func (f Functional) SetPos(x, y float64) {
	f.F().SetPos(x, y)
}
