package render

import (
	"container/heap"
	"image"
	"image/draw"

	"github.com/oakmound/oak/dlog"
)

//The RenderableHeap type is set up to manage a set of renderables to prevent any unsafe operations
// and allow for distinct updates between draw cycles
type RenderableHeap struct {
	rs     []Renderable
	toPush []Renderable
	static bool
}

//NewHeap creates a new renderableHeap
func NewHeap(static bool) *RenderableHeap {
	rh := new(RenderableHeap)
	rh.rs = make([]Renderable, 0)
	rh.toPush = make([]Renderable, 0)
	rh.static = static
	return rh
}

//Add stages a new Renderable to add to the heap
func (rh *RenderableHeap) Add(r Renderable, layer int) Renderable {
	r.SetLayer(layer)
	rh.toPush = append(rh.toPush, r)
	return r
}

//Replace adds a Renderable and removes an old one
func (rh *RenderableHeap) Replace(r1, r2 Renderable, layer int) {
	rh.Add(r2, layer)
	r1.UnDraw()
}

// Satisfying the Heap interface
//Len gets the length of the current heap
func (rh *RenderableHeap) Len() int { return len(rh.rs) }

//Less returns whether a renderable at index i is at a lower layer than the one at index j
func (rh *RenderableHeap) Less(i, j int) bool { return rh.rs[i].GetLayer() < rh.rs[j].GetLayer() }

//Swap moves two locations
func (rh *RenderableHeap) Swap(i, j int) { rh.rs[i], rh.rs[j] = rh.rs[j], rh.rs[i] }

//Push adds to the renderable heap
func (rh *RenderableHeap) Push(r interface{}) {
	defer func() {
		if x := recover(); x != nil {
			dlog.Error("Invalid Memory address pushed to Draw Heap")
		}
	}()
	if r == nil {
		return
	}
	// This can cause a 'name offset base pointer out of range' error
	// Maybe having incrementing sizes instead of appending could help that?
	rh.rs = append(rh.rs, r.(Renderable))
}

//Pop pops from the heap
func (rh *RenderableHeap) Pop() interface{} {
	n := len(rh.rs)
	x := rh.rs[n-1]
	rh.rs = rh.rs[0 : n-1]
	return x
}

// PreDraw parses through renderables to be pushed
// and adds them to the drawheap.
func (rh *RenderableHeap) PreDraw() {
	defer func() {
		if x := recover(); x != nil {
			dlog.Error("Invalid Memory Address in Draw heap")
			// This does not work-- all addresses following the bad address
			// at i are also bad
			//toPushRenderables = toPushRenderables[i+1:]
			rh.toPush = []Renderable{}
		}
	}()
	l := len(rh.toPush)
	for i := 0; i < l; i++ {
		r := rh.toPush[i]
		if r != nil {
			heap.Push(rh, r)
		}
	}
	rh.toPush = rh.toPush[l:]
}

// Copy on a renderableHeap does not include any of its elements,
// as renderables cannot be copied.
func (rh *RenderableHeap) Copy() Addable {
	rh2 := new(RenderableHeap)
	rh2.static = rh.static
	return rh2
}

func (rh *RenderableHeap) draw(world draw.Image, viewPos image.Point, screenW, screenH int) {
	newRh := &RenderableHeap{}
	if rh.static {
		for rh.Len() > 0 {
			rp := heap.Pop(rh)
			if rp != nil {
				r := rp.(Renderable)
				if r.GetLayer() != Undraw {
					r.Draw(world)
					heap.Push(newRh, r)
				}
			}
		}
		newRh.static = true
	} else {
		vx := float64(-viewPos.X)
		vy := float64(-viewPos.Y)
		for rh.Len() > 0 {
			intf := heap.Pop(rh)
			if intf != nil {
				r := intf.(Renderable)
				if r.GetLayer() != Undraw {
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
							r.DrawOffset(world, vx, vy)
						}
					}
					heap.Push(newRh, r)
				}
			}
		}
	}
	newRh.toPush = rh.toPush
	*rh = *newRh
}
