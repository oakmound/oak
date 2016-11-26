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
	LayeredPoint
	rs []Modifiable
}

func NewComposite(sl []Modifiable) *Composite {
	cs := new(Composite)
	cs.rs = sl
	return cs
}

func (cs *Composite) AppendOffset(r Modifiable, p Point) {
	r.SetPos(p.X, p.Y)
	cs.rs = append(cs.rs, r)
}

func (cs *Composite) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

func (cs *Composite) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

func (cs *Composite) AddOffset(i int, p Point) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X, p.Y)
	}
}

func (cs *Composite) SetOffsets(ps []Point) {
	for i, p := range ps {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(p.X, p.Y)
		}
	}
}

func (cs *Composite) Get(i int) Modifiable {
	return cs.rs[i]
}

func (cs *Composite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X+xOff, cs.Y+yOff)
	}
}
func (cs *Composite) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X, cs.Y)
	}
}
func (cs *Composite) UnDraw() {
	cs.layer = -1
	for _, c := range cs.rs {
		c.UnDraw()
	}
}
func (cs *Composite) GetRGBA() *image.RGBA {
	return nil
}
func (cs *Composite) AlwaysDirty() bool {
	return true
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
	cs2.layer = cs.layer
	cs2.X = cs.X
	cs2.Y = cs.Y
	cs2.rs = make([]Modifiable, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
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
	s += "}"
	return s
}

type CompositeR struct {
	LayeredPoint
	rs []Renderable
}

func NewCompositeR(sl []Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.rs = sl
	return cs
}

func (cs *CompositeR) AppendOffset(r Modifiable, p Point) {
	r.SetPos(p.X, p.Y)
	cs.rs = append(cs.rs, r)
}

func (cs *CompositeR) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

func (cs *CompositeR) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

func (cs *CompositeR) AddOffset(i int, p Point) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X, p.Y)
	}
}

func (cs *CompositeR) SetOffsets(ps []Point) {
	for i, p := range ps {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(p.X, p.Y)
		}
	}
}

func (cs *CompositeR) Get(i int) Renderable {
	return cs.rs[i]
}

func (cs *CompositeR) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X+xOff, cs.Y+yOff)
	}
}
func (cs *CompositeR) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X, cs.Y)
	}
}
func (cs *CompositeR) UnDraw() {
	cs.layer = -1
	for _, c := range cs.rs {
		c.UnDraw()
	}
}
func (cs *CompositeR) GetRGBA() *image.RGBA {
	return nil
}

func (cs *CompositeR) AlwaysDirty() bool {
	return true
}

func (cs *CompositeR) String() string {
	s := "CompositeR{"
	for _, v := range cs.rs {
		s += v.String() + "\n"
	}
	s += "}"
	return s
}
