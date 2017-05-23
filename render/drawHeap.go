package render

import (
	"container/heap"
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

type RenderableHeap struct {
	rs     []Renderable
	toPush []Renderable
	static bool
}

func NewHeap(static bool) *RenderableHeap {
	rh := new(RenderableHeap)
	rh.rs = make([]Renderable, 0)
	rh.toPush = make([]Renderable, 0)
	rh.static = static
	return rh
}

func (rh *RenderableHeap) Add(r Renderable, layer int) Renderable {
	r.SetLayer(layer)
	rh.toPush = append(rh.toPush, r)
	return r
}

// Satisfying the Heap interface
func (h *RenderableHeap) Len() int           { return len(h.rs) }
func (h *RenderableHeap) Less(i, j int) bool { return h.rs[i].GetLayer() < h.rs[j].GetLayer() }
func (h *RenderableHeap) Swap(i, j int)      { h.rs[i], h.rs[j] = h.rs[j], h.rs[i] }

func (h *RenderableHeap) Push(x interface{}) {
	if x == nil {
		return
	}
	h.rs = append(h.rs, x.(Renderable))
}

func (h *RenderableHeap) Pop() interface{} {
	n := len(h.rs)
	x := h.rs[n-1]
	h.rs = h.rs[0 : n-1]
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
	//fmt.Println("Static?", rh.static)
}

// Copying a renderableHeap does not include any of its elements,
// as renderables cannot be copied.
func (rh *RenderableHeap) Copy() Addable {
	rh2 := new(RenderableHeap)
	rh2.static = rh.static
	return rh2
}

func (rh *RenderableHeap) draw(world draw.Image, viewPos image.Point, screenW, screenH int) {
	newRh := &RenderableHeap{}
	if rh.static {
		//fmt.Println("There is a static draw heap")
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
		//fmt.Println("Heap length", len(rh.rs))
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
