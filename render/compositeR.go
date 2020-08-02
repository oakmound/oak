package render

import (
	"image"
	"image/draw"
	"sync"

	"github.com/oakmound/oak/v2/alg/floatgeom"
)

// A CompositeR is equivalent to a CompositeM for Renderables instead of
// Modifiables. CompositeRs also implements Stackable.
type CompositeR struct {
	LayeredPoint
	toPush      []Renderable
	toUndraw    []Renderable
	rs          []Renderable
	predrawLock sync.Mutex
	DrawPolygon
}

// NewCompositeR creates a new CompositeR from a slice of renderables
func NewCompositeR(sl ...Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.toPush = make([]Renderable, 0)
	cs.toUndraw = make([]Renderable, 0)
	cs.rs = sl
	return cs
}

// AppendOffset adds a new renderable to CompositeR with an offset
func (cs *CompositeR) AppendOffset(r Renderable, p floatgeom.Point2) {
	r.SetPos(p.X(), p.Y())
	cs.Append(r)
}

// AddOffset adds an offset to a given renderable of the slice
func (cs *CompositeR) AddOffset(i int, p floatgeom.Point2) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X(), p.Y())
	}
}

// Append adds a new renderable to the end of the CompositeR.
func (cs *CompositeR) Append(r Renderable) {
	cs.rs = append(cs.rs, r)
}

// Prepend adds a new renderable to the front of the CompositeR.
func (cs *CompositeR) Prepend(r Renderable) {
	cs.rs = append([]Renderable{r}, cs.rs...)
}

// Len returns the number of renderables in this composite.
func (cs *CompositeR) Len() int {
	return len(cs.rs)
}

// SetIndex places a renderable at a certain point in the composites renderable slice
func (cs *CompositeR) SetIndex(i int, r Renderable) {
	cs.rs[i] = r
}

// SetOffsets sets all renderables in CompositeR to the passed in Vector positions positions
func (cs *CompositeR) SetOffsets(ps ...floatgeom.Point2) {
	for i, p := range ps {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(p.X(), p.Y())
		}
	}
}

// DrawOffset Draws the CompositeR with an offset from its logical location.
func (cs *CompositeR) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

// Draw draws the CompositeR at its logical location and therefore its consituent renderables as well
func (cs *CompositeR) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

// Undraw undraws the CompositeR and its consituent renderables
func (cs *CompositeR) Undraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.Undraw()
	}
}

// GetRGBA always returns nil from Composites
func (cs *CompositeR) GetRGBA() *image.RGBA {
	return nil
}

// Get returns renderable from a given index in CompositeR
func (cs *CompositeR) Get(i int) Renderable {
	return cs.rs[i]
}

// Add stages a renderable to be added to the Composite at the next PreDraw
func (cs *CompositeR) Add(r Renderable, _ ...int) Renderable {
	cs.predrawLock.Lock()
	cs.toPush = append(cs.toPush, r)
	cs.predrawLock.Unlock()
	return r
}

// Replace updates a renderable in the CompositeR to the new Renderable
func (cs *CompositeR) Replace(old, new Renderable, i int) {
	cs.predrawLock.Lock()
	cs.toPush = append(cs.toPush, new)
	cs.toUndraw = append(cs.toUndraw, old)
	cs.predrawLock.Unlock()

}

// PreDraw updates the CompositeR with the new renderables to add.
// This helps keep consistency and mitigates the threat of unsafe operations.
func (cs *CompositeR) PreDraw() {
	cs.predrawLock.Lock()
	push := cs.toPush
	cs.toPush = []Renderable{}
	cs.rs = append(cs.rs, push...)
	for _, r := range cs.toUndraw {
		r.Undraw()
	}
	cs.predrawLock.Unlock()
}

// Copy returns a new composite with the same length slice of renderables but no actual renderables...
// CompositeRs cannot have their internal elements copied,
// as renderables cannot be copied.
func (cs *CompositeR) Copy() Stackable {
	cs2 := new(CompositeR)
	cs2.LayeredPoint = cs.LayeredPoint
	cs2.rs = make([]Renderable, len(cs.rs))
	return cs2
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
		x := int(r.X())
		y := int(r.Y())
		x2 := x
		y2 := y
		w, h := r.GetDims()
		x += w
		y += h
		if x > viewPos.X && y > viewPos.Y &&
			x2 < viewPos.X+screenW && y2 < viewPos.Y+screenH {

			if cs.InDrawPolygon(x, y, x2, y2) {
				r.DrawOffset(world, float64(-viewPos.X), float64(-viewPos.Y))
			}
		}
	}
	cs.rs = cs.rs[0:realLength]
}
