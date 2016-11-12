package render

import (
	"image"
	"image/color"
	"image/draw"
)

// Composite Types, distinct from Compound Types,
// Display all of their parts at the same time,
// and respect the positions and layers of their
// parts.
type Composite struct {
	rs      []Modifiable
	offsets []Point
}

func NewComposite(sl []Modifiable) *Composite {
	cs := new(Composite)
	cs.rs = sl
	cs.offsets = make([]Point, len(sl))
	return cs
}

func (cs *Composite) AppendOffset(r Modifiable, p Point) {
	cs.rs = append(cs.rs, r)
	cs.offsets = append(cs.offsets, p)
}

func (cs *Composite) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
	cs.offsets = append(cs.offsets, Point{})
}

func (cs *Composite) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

func (cs *Composite) AddOffset(i int, p Point) {
	cs.offsets[i] = p
}

func (cs *Composite) SetOffsets(ps []Point) {
	for i, p := range ps {
		cs.offsets[i] = p
	}
}

func (cs *Composite) Get(i int) Modifiable {
	return cs.rs[i]
}

func (cs *Composite) Draw(buff draw.Image) {
	for i, c := range cs.rs {
		switch t := c.(type) {
		case *Composite:
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
		//if c.AlwaysDirty() || IsDirty(drawX, drawY) {
		ShinyDraw(buff, img, drawX, drawY)
		//}
	}
}
func (cs *Composite) GetRGBA() *image.RGBA {
	return nil
}
func (cs *Composite) ShiftX(x float64) {
	for _, v := range cs.rs {
		v.ShiftX(x)
	}
}
func (cs *Composite) ShiftY(y float64) {
	for _, v := range cs.rs {
		v.ShiftY(y)
	}
}
func (cs *Composite) AlwaysDirty() bool {
	return true
}
func (cs *Composite) GetX() float64 {
	return 0.0
}
func (cs *Composite) GetY() float64 {
	return 0.0
}

// This should be changed so that compositeSlice (and map)
// has a persistent concept of what it's smallest
// x and y are.
func (cs *Composite) SetPos(x, y float64) {
	for _, v := range cs.rs {
		v.SetPos(x, y)
	}
}
func (cs *Composite) GetLayer() int {
	return 0
}
func (cs *Composite) SetLayer(l int) {
	for _, v := range cs.rs {
		v.SetLayer(l)
	}
}
func (cs *Composite) UnDraw() {
	for _, v := range cs.rs {
		v.UnDraw()
	}
}

func (cs *Composite) FlipX() Modifiable {
	for _, v := range cs.rs {
		v.FlipX()
	}
	return cs
}
func (cs *Composite) FlipY() Modifiable {
	for _, v := range cs.rs {
		v.FlipY()
	}
	return cs
}
func (cs *Composite) ApplyColor(c color.Color) Modifiable {
	for _, v := range cs.rs {
		v.ApplyColor(c)
	}
	return cs
}
func (cs *Composite) Copy() Modifiable {
	cs2 := new(Composite)
	cs2.rs = make([]Modifiable, len(cs.rs))
	cs2.offsets = make([]Point, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
		cs2.offsets[i] = cs.offsets[i]
	}
	return cs2
}
func (cs *Composite) FillMask(img image.RGBA) Modifiable {
	for _, v := range cs.rs {
		v.FillMask(img)
	}
	return cs
}
func (cs *Composite) ApplyMask(img image.RGBA) Modifiable {
	for _, v := range cs.rs {
		v.ApplyMask(img)
	}
	return cs
}
func (cs *Composite) Rotate(degrees int) Modifiable {
	for _, v := range cs.rs {
		v.Rotate(degrees)
	}
	return cs
}
func (cs *Composite) Scale(xRatio float64, yRatio float64) Modifiable {
	for _, v := range cs.rs {
		v.Scale(xRatio, yRatio)
	}
	return cs
}
func (cs *Composite) Fade(alpha int) Modifiable {
	for _, v := range cs.rs {
		v.Fade(alpha)
	}
	return cs
}

func (cs *Composite) String() string {
	s := "Composite{"
	for _, v := range cs.rs {
		s += v.String() + "\n"
	}
	return s
}

type CompositeR struct {
	rs      []Renderable
	offsets []Point
	unDraw  bool
}

func NewCompositeR(sl []Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.rs = sl
	cs.offsets = make([]Point, len(sl))
	return cs
}

func (cs *CompositeR) AppendOffset(r Renderable, p Point) {
	cs.rs = append(cs.rs, r)
	cs.offsets = append(cs.offsets, p)
}

func (cs *CompositeR) Append(r Renderable) {
	cs.rs = append(cs.rs, r)
	cs.offsets = append(cs.offsets, Point{})
}

func (cs *CompositeR) Add(i int, r Renderable) {
	cs.rs[i] = r
}

func (cs *CompositeR) AddOffset(i int, p Point) {
	cs.offsets[i] = p
}

func (cs *CompositeR) SetOffsets(ps []Point) {
	for i, p := range ps {
		cs.offsets[i] = p
	}
}

func (cs *CompositeR) Get(i int) Renderable {
	return cs.rs[i]
}

func (cs *CompositeR) Draw(buff draw.Image) {
	for i, c := range cs.rs {
		switch t := c.(type) {
		case *CompositeR:
			t.Draw(buff)
			continue
		case *Text:
			t.Draw(buff)
			continue
		}
		img := c.GetRGBA()
		drawX := int(c.GetX()) + int(cs.offsets[i].X)
		drawY := int(c.GetY()) + int(cs.offsets[i].Y)
		ShinyDraw(buff, img, drawX, drawY)
	}
}
func (cs *CompositeR) GetRGBA() *image.RGBA {
	return nil
}
func (cs *CompositeR) ShiftX(x float64) {
	for _, v := range cs.rs {
		v.ShiftX(x)
	}
}
func (cs *CompositeR) ShiftY(y float64) {
	for _, v := range cs.rs {
		v.ShiftY(y)
	}
}
func (cs *CompositeR) AlwaysDirty() bool {
	return true
}
func (cs *CompositeR) GetX() float64 {
	return 0.0
}
func (cs *CompositeR) GetY() float64 {
	return 0.0
}

// This should be changed so that compositeSlice (and map)
// has a persistent concept of what it's smallest
// x and y are.
func (cs *CompositeR) SetPos(x, y float64) {
	for _, v := range cs.rs {
		v.SetPos(x, y)
	}
}
func (cs *CompositeR) GetLayer() int {
	if cs.unDraw {
		return -1
	}
	return 0
}
func (cs *CompositeR) SetLayer(l int) {
	for _, v := range cs.rs {
		v.SetLayer(l)
	}
}
func (cs *CompositeR) UnDraw() {
	for _, v := range cs.rs {
		v.UnDraw()
	}
	cs.unDraw = true
}

func (cs *CompositeR) String() string {
	s := "CompositeR{"
	for _, v := range cs.rs {
		s += v.String() + "\n"
	}
	return s
}
