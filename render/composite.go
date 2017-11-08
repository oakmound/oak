package render

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/physics"
)

// Composite Types, distinct from Compound Types,
// Display all of their parts at the same time,
// and respect the positions and layers of their
// parts.
type Composite struct {
	LayeredPoint
	rs []Modifiable
}

//NewComposite creates a Composite
func NewComposite(sl []Modifiable) *Composite {
	cs := new(Composite)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.rs = sl
	return cs
}

//AppendOffset adds a new offset modifiable to the composite
func (cs *Composite) AppendOffset(r Modifiable, v physics.Vector) {
	r.SetPos(v.X(), v.Y())
	cs.rs = append(cs.rs, r)
}

//Append adds a renderable as is to the composite
func (cs *Composite) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

//Add places a renderable at a certain point in the composites renderable slice
func (cs *Composite) Add(i int, r Modifiable) {
	cs.rs[i] = r
}

//AddOffset offsets all renderables in the composite by a vector
func (cs *Composite) AddOffset(i int, v physics.Vector) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(v.X(), v.Y())
	}
}

//SetOffsets applies the initial offsets to the entire Composite
func (cs *Composite) SetOffsets(vs []physics.Vector) {
	for i, v := range vs {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(v.X(), v.Y())
		}
	}
}

//Get returns a renderable at the given index within the composite
func (cs *Composite) Get(i int) Modifiable {
	return cs.rs[i]
}

//DrawOffset draws the Composite with some offset from its logical position (and therefore sub renderables logical positions).
func (cs *Composite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

//Draw draws the Composite at its logical position
func (cs *Composite) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

//UnDraw stops the composite from being drawn
func (cs *Composite) UnDraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.UnDraw()
	}
}

//GetRGBA does not work on a composite and therefore returns nil
func (cs *Composite) GetRGBA() *image.RGBA {
	return nil
}

//AlwaysDirty shows that the Composite always needs updating
func (cs *Composite) AlwaysDirty() bool {
	return true
}

//Modify applies modifications to the composite
func (cs *Composite) Modify(ms ...Modification) Modifiable {
	for _, r := range cs.rs {
		r.Modify(ms...)
	}
	return cs
}

//Copy makes a new Composite with the same renderables
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

func (cs *Composite) ToSprite() *Sprite {
	w, h := 0.0, 0.0
	for _, v := range cs.rs {
		x := v.GetX()
		y := v.GetY()
		w2, h2 := v.GetDims()
		if float64(w2)+x > w {
			w = float64(w2) + x
		}
		if float64(h2)+y > h {
			h = float64(h2) + y
		}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	cs.Draw(rgba)
	return NewSprite(cs.X(), cs.Y(), rgba)
}

//CompositeR keeps track of a set of renderables at a location
type CompositeR struct {
	LayeredPoint
	toPush []Renderable
	rs     []Renderable
}

//NewCompositeR creates a new CompositeR from a slice of renderables
func NewCompositeR(sl []Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.toPush = make([]Renderable, 0)
	cs.rs = sl
	return cs
}

//AppendOffset adds a new renderable to CompositeR with an offset
func (cs *CompositeR) AppendOffset(r Renderable, v physics.Vector) {
	r.SetPos(v.X(), v.Y())
	cs.rs = append(cs.rs, r)
}

//Append adds a new renderable to CompositeR
func (cs *CompositeR) Append(r Renderable) {
	cs.rs = append(cs.rs, r)
}

//Add stages a renderable to be added to CompositeR at a give position in the slice
func (cs *CompositeR) Add(r Renderable, i int) Renderable {
	cs.toPush = append(cs.toPush, r)
	return r
}

//Replace updates a renderable in the CompositeR to the new Renderable
func (cs *CompositeR) Replace(r1, r2 Renderable, i int) {
	cs.Add(r2, i)
	r1.UnDraw()
}

//AddOffset adds an offset to a given renderable of the slice
func (cs *CompositeR) AddOffset(i int, v physics.Vector) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(v.X(), v.Y())
	}
}

//PreDraw updates the CompositeR with the new renderables to add. This helps keep consistency and mitigates the threat of unsafe operations.
func (cs *CompositeR) PreDraw() {
	push := cs.toPush
	cs.toPush = []Renderable{}
	cs.rs = append(cs.rs, push...)
}

// Copy returns a new composite with the same length slice of renderables but no actual renderables...
// CompositeRs cannot have their internal elements copied,
// as renderables cannot be copied.
func (cs *CompositeR) Copy() Addable {
	cs2 := new(CompositeR)
	cs2.LayeredPoint = cs.LayeredPoint
	cs2.rs = make([]Renderable, len(cs.rs))
	return cs2
}

//SetOffsets sets all renderables in CompositeR to the passed in Vector positions positions
func (cs *CompositeR) SetOffsets(ps []physics.Vector) {
	for i, p := range ps {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(p.X(), p.Y())
		}
	}
}

//Get returns renderable at given location in CompositeR
func (cs *CompositeR) Get(i int) Renderable {
	return cs.rs[i]
}

//DrawOffset Draws the CompositeR with an offset from its logical location.
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

//Draw draws the CompositeR at its logical location and therefore its consituent renderables as well
func (cs *CompositeR) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

//UnDraw undraws the CompositeR and therefore its consituent renderables as well
func (cs *CompositeR) UnDraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.UnDraw()
	}
}

//GetRGBA does not work on composites and returns nil
func (cs *CompositeR) GetRGBA() *image.RGBA {
	return nil
}

//AlwaysDirty notes that CompositeR is alwaysdirty
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
