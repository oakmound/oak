package render

import (
	"container/heap"
	"image"
	"image/draw"
	"sync"
)

// A RenderableHeap manages a set of renderables to be drawn in explicit layered
// order, using an internal heap to manage that order.
type RenderableHeap struct {
	rs      []Renderable
	toPush  []Renderable
	static  bool
	addLock sync.RWMutex
}

// NewHeap creates a new renderableHeap. The static boolean represents whether
// this heap exists relative to the viewport or not-- if true, an element at 40,40
// will always be at 40,40. If false, when the viewport moves, the element will
// move opposite the direction of the viewport.
func NewHeap(static bool) *RenderableHeap {
	rh := new(RenderableHeap)
	rh.rs = make([]Renderable, 0)
	rh.toPush = make([]Renderable, 0)
	rh.static = static
	rh.addLock = sync.RWMutex{}
	return rh
}

//Add stages a new Renderable to add to the heap
func (rh *RenderableHeap) Add(r Renderable, layers ...int) Renderable {
	if len(layers) > 0 {
		r.SetLayer(layers[0])
	}
	rh.addLock.Lock()
	rh.toPush = append(rh.toPush, r)
	rh.addLock.Unlock()
	return r
}

//Replace adds a Renderable and removes an old one
func (rh *RenderableHeap) Replace(r1, r2 Renderable, layer int) {
	rh.Add(r2, layer)
	r1.Undraw()
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
	if r == nil {
		return
	}
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
	rh.addLock.Lock()
	for _, r := range rh.toPush {
		if r != nil {
			heap.Push(rh, r)
		}
	}
	rh.toPush = make([]Renderable, 0)
	rh.addLock.Unlock()
}

// Copy on a renderableHeap does not include any of its elements,
// as renderables cannot be copied.
func (rh *RenderableHeap) Copy() Stackable {
	return NewHeap(rh.static)
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
	} else {
		vx := float64(-viewPos.X)
		vy := float64(-viewPos.Y)
		for rh.Len() > 0 {
			intf := heap.Pop(rh)
			if intf != nil {
				r := intf.(Renderable)
				if r.GetLayer() != Undraw {
					x2 := int(r.X())
					y2 := int(r.Y())
					w, h := r.GetDims()
					x := w + x2
					y := h + y2
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
	rh.rs = newRh.rs
}
