package render

import (
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
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
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.rs = sl
	return cs
}

func (cs *Composite) AppendOffset(r Modifiable, v physics.Vector) {
	r.SetPos(v.X(), v.Y())
	cs.rs = append(cs.rs, r)
}

func (cs *Composite) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

func (cs *Composite) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

func (cs *Composite) AddOffset(i int, v physics.Vector) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(v.X(), v.Y())
	}
}

func (cs *Composite) SetOffsets(vs []physics.Vector) {
	for i, v := range vs {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(v.X(), v.Y())
		}
	}
}

func (cs *Composite) Get(i int) Modifiable {
	return cs.rs[i]
}

func (cs *Composite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}
func (cs *Composite) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}
func (cs *Composite) UnDraw() {
	cs.layer = Undraw
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

func (cs *Composite) Modify(ms ...Modification) Modifiable {
	for _, r := range cs.rs {
		r.Modify(ms...)
	}
	return cs
}
func (cs *Composite) Copy() Modifiable {
	cs2 := new(Composite)
	cs2.layer = cs.layer
	cs2.Vector = cs.Vector
	cs2.rs = make([]Modifiable, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
	}
	return cs2
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
	toPush []Renderable
	rs     []Renderable
}

func NewCompositeR(sl []Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.toPush = make([]Renderable, 0)
	cs.rs = sl
	return cs
}

func (cs *CompositeR) AppendOffset(r Renderable, v physics.Vector) {
	r.SetPos(v.X(), v.Y())
	cs.rs = append(cs.rs, r)
}

func (cs *CompositeR) Append(r Renderable) {
	cs.rs = append(cs.rs, r)
}

func (cs *CompositeR) Add(r Renderable, i int) Renderable {
	cs.toPush = append(cs.toPush, r)
	return r
}

func (cs *CompositeR) Replace(r1, r2 Renderable, i int) {
	cs.Add(r2, i)
	r1.UnDraw()
}

func (cs *CompositeR) AddOffset(i int, v physics.Vector) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(v.X(), v.Y())
	}
}

func (cs *CompositeR) PreDraw() {
	push := cs.toPush
	cs.toPush = []Renderable{}
	cs.rs = append(cs.rs, push...)
}

// CompositeRs cannot have their internal elements copied,
// as renderables cannot be copied.
func (cs *CompositeR) Copy() Addable {
	cs2 := new(CompositeR)
	cs2.LayeredPoint = cs.LayeredPoint
	cs2.rs = make([]Renderable, len(cs.rs))
	return cs2
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
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

func (cs *CompositeR) draw(world draw.Image, viewPos image.Point, screenW, screenH int) {
	realLength := len(cs.rs)
	for i := 0; i < realLength; i++ {
		r := cs.rs[i]
		for (r == nil || r.GetLayer() == Undraw) && realLength > i {
			cs.rs[i], cs.rs[realLength-1] = cs.rs[realLength-1], cs.rs[i]
			realLength--
			r = cs.rs[i]
		}
		if realLength == i {
			break
		}
		x := int(r.GetX())
		y := int(r.GetY())
		x2 := x
		y2 := y
		w, h := r.GetDims()
		x += w
		y += h
		if x > viewPos.X && y > viewPos.Y &&
			x2 < viewPos.X+screenW && y2 < viewPos.Y+screenH {

			if InDrawPolygon(x, y, x2, y2) {
				r.DrawOffset(world, float64(-viewPos.X), float64(-viewPos.Y))
			}
		}
	}
	cs.rs = cs.rs[0:realLength]
}

func (cs *CompositeR) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

func (cs *CompositeR) UnDraw() {
	cs.layer = Undraw
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
