package render

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type CompositeSlice []Modifiable

func NewCompositeSlice(sl []Modifiable) *CompositeSlice {
	cs := CompositeSlice(sl)
	return &cs
}

func CompositeFilter(s Sprite) *CompositeSlice {
	rgba := s.GetRGBA()
	b := rgba.Bounds()
	if b.Max.X < dirtyWidth && b.Max.Y < dirtyHeight {
		panic("Invalid dimensioned sprite sent to composite filter")
	}
	cs := CompositeSlice{}
	for x := 0; x < b.Max.X; x += dirtyWidth {
		for y := 0; y < b.Max.Y; y += dirtyHeight {
			rgba2 := subImage(rgba, x, y, x+dirtyWidth, y+dirtyHeight)
			s2 := NewSprite(float64(x), float64(y), rgba2)
			cs = append(cs, s2)
		}
	}
	return &cs
}

func (cs *CompositeSlice) Append(r Modifiable) {
	*cs = append(*cs, r)
}

func (cs *CompositeSlice) Add(i int, r Modifiable) {
	(*cs)[i] = r
}

func (cs *CompositeSlice) Get(i int) Modifiable {
	return (*cs)[i]
}

func (cs *CompositeSlice) Draw(buff draw.Image) {
	for _, r := range *cs {
		x := int(r.GetX())
		y := int(r.GetY())
		if r.AlwaysDirty() || IsDirty(x, y) {
			r.Draw(buff)
		}
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
func (cs *CompositeSlice) AlwaysDirty() bool {
	return true
}
func (cs *CompositeSlice) GetX() float64 {
	return 0.0
}
func (cs *CompositeSlice) GetY() float64 {
	return 0.0
}

// This should be changed so that compositeSlice (and map)
// has a persistent concept of what it's smallest
// x and y are.
func (cs *CompositeSlice) SetPos(x, y float64) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	for _, v := range *cs {
		if minX > v.GetX() {
			minX = v.GetX()
		}
		if minY > v.GetY() {
			minY = v.GetY()
		}
	}
	for _, v := range *cs {
		v.SetPos(x+v.GetX()-minX, y+v.GetY()-minY)
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

func (cs *CompositeSlice) FlipX() Modifiable {
	for _, v := range *cs {
		v.FlipX()
	}
	return cs
}
func (cs *CompositeSlice) FlipY() Modifiable {
	for _, v := range *cs {
		v.FlipY()
	}
	return cs
}
func (cs *CompositeSlice) ApplyColor(c color.Color) Modifiable {
	for _, v := range *cs {
		v.ApplyColor(c)
	}
	return cs
}
func (cs *CompositeSlice) Copy() Modifiable {
	cs2 := CompositeSlice{}
	for _, v := range *cs {
		cs2 = append(cs2, v.Copy())
	}
	return &cs2
}
func (cs *CompositeSlice) FillMask(img image.RGBA) Modifiable {
	for _, v := range *cs {
		v.FillMask(img)
	}
	return cs
}
func (cs *CompositeSlice) ApplyMask(img image.RGBA) Modifiable {
	for _, v := range *cs {
		v.ApplyMask(img)
	}
	return cs
}
func (cs *CompositeSlice) Rotate(degrees int) Modifiable {
	for _, v := range *cs {
		v.Rotate(degrees)
	}
	return cs
}
func (cs *CompositeSlice) Scale(xRatio float64, yRatio float64) Modifiable {
	for _, v := range *cs {
		v.Scale(xRatio, yRatio)
	}
	return cs
}
func (cs *CompositeSlice) Fade(alpha int) Modifiable {
	for _, v := range *cs {
		v.Fade(alpha)
	}
	return cs
}
