package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

// Composite Types, distinct from Compound Types,
// Display all of their parts at the same time,
// and respect the positions and layers of their
// parts.
type CompositeMap map[string]Renderable

func NewCompositeMap(m map[string]Renderable) *CompositeMap {
	cm := CompositeMap(m)
	return &cm
}

func (cm *CompositeMap) Add(s string, r Renderable) {
	(*cm)[s] = r
}

func (cm *CompositeMap) Get(s string) Renderable {
	return (*cm)[s]
}

func (cm *CompositeMap) Draw(buff screen.Buffer) {
	for _, v := range *cm {
		v.Draw(buff)
	}
}
func (cm *CompositeMap) GetRGBA() *image.RGBA {
	return nil
}
func (cm *CompositeMap) ShiftX(x float64) {
	for _, v := range *cm {
		v.ShiftX(x)
	}
}
func (cm *CompositeMap) ShiftY(y float64) {
	for _, v := range *cm {
		v.ShiftY(y)
	}
}
func (cm *CompositeMap) SetPos(x, y float64) {
	for _, v := range *cm {
		v.SetPos(x, y)
	}
}
func (cm *CompositeMap) GetLayer() int {
	return 0
}
func (cm *CompositeMap) SetLayer(l int) {
	for _, v := range *cm {
		v.SetLayer(l)
	}
}
func (cm *CompositeMap) UnDraw() {
	for _, v := range *cm {
		v.UnDraw()
	}
}
