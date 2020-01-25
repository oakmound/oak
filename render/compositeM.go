package render

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/render/mod"
)

// CompositeM Types display all of their parts at the same time,
// and respect the positions of their parts as relative to the
// position of the composite itself
type CompositeM struct {
	LayeredPoint
	rs []Modifiable
}

// NewCompositeM creates a CompositeM
func NewCompositeM(sl ...Modifiable) *CompositeM {
	cs := new(CompositeM)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.rs = sl
	return cs
}

// AppendOffset adds a new offset modifiable to the CompositeM
func (cs *CompositeM) AppendOffset(r Modifiable, p floatgeom.Point2) {
	r.SetPos(p.X(), p.Y())
	cs.Append(r)
}

// Append adds a renderable as is to the CompositeM
func (cs *CompositeM) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

// Prepend adds a new renderable to the front of the CompositeMR.
func (cs *CompositeM) Prepend(r Modifiable) {
	cs.rs = append([]Modifiable{r}, cs.rs...)
}

// SetIndex places a renderable at a certain point in the CompositeMs renderable slice
func (cs *CompositeM) SetIndex(i int, r Modifiable) {
	cs.rs[i] = r
}

// Slice creates a new CompositeM as a subslice of the existing CompositeM.
// No Modifiables will be copied, and the original will not be modified.
func (cs *CompositeM) Slice(start, end int) *CompositeM {
	if start < 0 {
		start = 0
	}
	if end > len(cs.rs) {
		end = len(cs.rs)
	}
	newRs := cs.rs[start:end]
	return &CompositeM{
		LayeredPoint: cs.LayeredPoint.Copy(),
		rs:           newRs,
	}
}

// Len returns the number of renderables in this CompositeM.
func (cs *CompositeM) Len() int {
	return len(cs.rs)
}

// AddOffset offsets all renderables in the CompositeM by a vector
func (cs *CompositeM) AddOffset(i int, p floatgeom.Point2) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X(), p.Y())
	}
}

// SetOffsets applies the initial offsets to the entire CompositeM
func (cs *CompositeM) SetOffsets(vs ...floatgeom.Point2) {
	for i, v := range vs {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(v.X(), v.Y())
		}
	}
}

// Get returns a renderable at the given index within the CompositeM
func (cs *CompositeM) Get(i int) Modifiable {
	return cs.rs[i]
}

// DrawOffset draws the CompositeM with some offset from its logical position
// (and therefore sub renderables logical positions).
func (cs *CompositeM) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

// Draw draws the CompositeM at its logical position
func (cs *CompositeM) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

// Undraw stops the CompositeM from being drawn
func (cs *CompositeM) Undraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.Undraw()
	}
}

// GetRGBA always returns nil from Composites
func (cs *CompositeM) GetRGBA() *image.RGBA {
	return nil
}

// Modify applies mods to the CompositeM
func (cs *CompositeM) Modify(ms ...mod.Mod) Modifiable {
	for _, r := range cs.rs {
		r.Modify(ms...)
	}
	return cs
}

// Filter filters each component part of this CompositeM by all of the inputs.
func (cs *CompositeM) Filter(fs ...mod.Filter) {
	for _, r := range cs.rs {
		r.Filter(fs...)
	}
}

// ToSprite converts the composite into a sprite by drawing each layer in order
// and overwriting lower layered pixels
func (cs *CompositeM) ToSprite() *Sprite {
	var maxW, maxH int
	for _, r := range cs.rs {
		x, y := int(r.X()), int(r.Y())
		w, h := r.GetDims()
		if x+w > maxW {
			maxW = x + w
		}
		if y+h > maxH {
			maxH = y + h
		}
	}
	sp := NewEmptySprite(cs.X(), cs.Y(), maxW, maxH)
	for _, r := range cs.rs {
		r.Draw(sp)
	}
	return sp
}

// Copy makes a new CompositeM with the same renderables
func (cs *CompositeM) Copy() Modifiable {
	cs2 := new(CompositeM)
	cs2.layer = cs.layer
	cs2.Vector = cs.Vector.Copy()
	cs2.rs = make([]Modifiable, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
	}
	return cs2
}
