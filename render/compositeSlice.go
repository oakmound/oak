package render

import (
	"image"
	"image/draw"
)

type CompositeSlice []Renderable

func NewCompositeSlice(sl []Renderable) *CompositeSlice {
	cs := CompositeSlice(sl)
	return &cs
}

func (cs *CompositeSlice) Append(r Renderable) {
	*cs = append(*cs, r)
}

func (cs *CompositeSlice) Add(i int, r Renderable) {
	(*cs)[i] = r
}

func (cs *CompositeSlice) Get(i int) Renderable {
	return (*cs)[i]
}

func (cs *CompositeSlice) Draw(buff draw.Image) {
	for _, v := range *cs {
		v.Draw(buff)
	}
}
func (cs *CompositeSlice) GetRGBA() *image.RGBA {
	return nil
}
func (cs *CompositeSlice) ShiftX(x float64) {
	for _, v := range *cs {
		v.ShiftX(x)
	}
}
func (cs *CompositeSlice) ShiftY(y float64) {
	for _, v := range *cs {
		v.ShiftY(y)
	}
}
func (cs *CompositeSlice) SetPos(x, y float64) {
	for _, v := range *cs {
		v.SetPos(x, y)
	}
}
func (cs *CompositeSlice) GetLayer() int {
	return 0
}
func (cs *CompositeSlice) SetLayer(l int) {
	for _, v := range *cs {
		v.SetLayer(l)
	}
}
func (cs *CompositeSlice) UnDraw() {
	for _, v := range *cs {
		v.UnDraw()
	}
}
