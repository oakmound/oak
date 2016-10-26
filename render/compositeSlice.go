package render

import (
	"image"
	"image/color"
	"image/draw"
)

type CompositeSlice struct {
	rs      []Modifiable
	offsets []Point
}

func NewCompositeSlice(sl []Modifiable) *CompositeSlice {
	cs := new(CompositeSlice)
	cs.rs = sl
	cs.offsets = make([]Point, len(sl))
	return cs
}

func (cs *CompositeSlice) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
	cs.offsets = append(cs.offsets, Point{})
}

func (cs *CompositeSlice) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

func (cs *CompositeSlice) AddOffset(i int, p Point) {
	cs.offsets[i] = p
}

func (cs *CompositeSlice) SetOffsets(ps []Point) {
	for i, p := range ps {
		cs.offsets[i] = p
	}
}

func (cs *CompositeSlice) Get(i int) Modifiable {
	return cs.rs[i]
}

func (cs *CompositeSlice) Draw(buff draw.Image) {
	for i, c := range cs.rs {
		switch t := c.(type) {
		case *CompositeSlice:
			t.Draw(buff)
			continue
		case *Reverting:
			t.updateAnimation()
		case *Animation:
			t.updateAnimation()
		case *Sequence:
			t.update()
		}
		img := c.GetRGBA()
		drawX := int(c.GetX()) + int(cs.offsets[i].X)
		drawY := int(c.GetY()) + int(cs.offsets[i].Y)
		ShinyDraw(buff, img, drawX, drawY)
	}
}
func (cs *CompositeSlice) GetRGBA() *image.RGBA {
	return nil
}
func (cs *CompositeSlice) ShiftX(x float64) {
	for _, v := range cs.rs {
		v.ShiftX(x)
	}
}
func (cs *CompositeSlice) ShiftY(y float64) {
	for _, v := range cs.rs {
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
	for _, v := range cs.rs {
		v.SetPos(x, y)
	}
}
func (cs *CompositeSlice) GetLayer() int {
	return 0
}
func (cs *CompositeSlice) SetLayer(l int) {
	for _, v := range cs.rs {
		v.SetLayer(l)
	}
}
func (cs *CompositeSlice) UnDraw() {
	for _, v := range cs.rs {
		v.UnDraw()
	}
}

func (cs *CompositeSlice) FlipX() Modifiable {
	for _, v := range cs.rs {
		v.FlipX()
	}
	return cs
}
func (cs *CompositeSlice) FlipY() Modifiable {
	for _, v := range cs.rs {
		v.FlipY()
	}
	return cs
}
func (cs *CompositeSlice) ApplyColor(c color.Color) Modifiable {
	for _, v := range cs.rs {
		v.ApplyColor(c)
	}
	return cs
}
func (cs *CompositeSlice) Copy() Modifiable {
	cs2 := new(CompositeSlice)
	cs2.rs = make([]Modifiable, len(cs.rs))
	cs2.offsets = make([]Point, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
		cs2.offsets[i] = cs.offsets[i]
	}
	return cs2
}
func (cs *CompositeSlice) FillMask(img image.RGBA) Modifiable {
	for _, v := range cs.rs {
		v.FillMask(img)
	}
	return cs
}
func (cs *CompositeSlice) ApplyMask(img image.RGBA) Modifiable {
	for _, v := range cs.rs {
		v.ApplyMask(img)
	}
	return cs
}
func (cs *CompositeSlice) Rotate(degrees int) Modifiable {
	for _, v := range cs.rs {
		v.Rotate(degrees)
	}
	return cs
}
func (cs *CompositeSlice) Scale(xRatio float64, yRatio float64) Modifiable {
	for _, v := range cs.rs {
		v.Scale(xRatio, yRatio)
	}
	return cs
}
func (cs *CompositeSlice) Fade(alpha int) Modifiable {
	for _, v := range cs.rs {
		v.Fade(alpha)
	}
	return cs
}

func (cs *CompositeSlice) String() string {
	s := "CompositeSlice{"
	for _, v := range cs.rs {
		s += v.String() + "\n"
	}
	return s
}
