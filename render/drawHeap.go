package render

import (
	"image/draw"
	"sync"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

// A RenderableHeap manages a set of renderables to be drawn in explicit layered
// order, using an internal heap to manage that order. It implements Stackable.
type RenderableHeap struct {
	layerHeap
	toPush   []Renderable
	toUndraw []Renderable
	static   bool
	addLock  sync.RWMutex
}

func newHeap(static bool) *RenderableHeap {
	rh := new(RenderableHeap)
	rh.rs = make([]Renderable, 0)
	rh.toPush = make([]Renderable, 0)
	rh.toUndraw = make([]Renderable, 0)
	rh.static = static
	rh.addLock = sync.RWMutex{}
	return rh
}

// NewDynamicHeap creates a renderable heap for drawing renderables by layer
// where the position of the viewport is taken into account to produce the drawn
// location of the renderable.
//
// Example:
// If drawing a Sprite at (100,100) with the viewport at (50,0), the sprite will
// appear at (50, 100).
func NewDynamicHeap() *RenderableHeap {
	return newHeap(false)
}

// NewStaticHeap creates a renderable heap for drawing renderables by layer
// where the position of renderable is absolute with regards to the viewport.
//
// Example:
// If drawing a Sprite at (100,100) with the viewport at (50,0), the sprite will
// appear at (100, 100).
func NewStaticHeap() *RenderableHeap {
	return newHeap(true)
}

func (rh *RenderableHeap) Clear() {
	*rh = *newHeap(rh.static)
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

// Replace adds a Renderable and removes an old one
func (rh *RenderableHeap) Replace(old, new Renderable, layer int) {
	new.SetLayer(layer)
	rh.addLock.Lock()
	rh.toPush = append(rh.toPush, new)
	rh.toUndraw = append(rh.toUndraw, old)
	rh.addLock.Unlock()
}

// PreDraw parses through renderables to be pushed
// and adds them to the drawheap.
func (rh *RenderableHeap) PreDraw() {
	rh.addLock.Lock()
	for _, r := range rh.toPush {
		if r != nil {
			rh.heapPush(r)
		}
	}
	for _, r := range rh.toUndraw {
		if r != nil {
			r.Undraw()
		}
	}
	rh.toPush = make([]Renderable, 0)
	rh.addLock.Unlock()
}

// Copy on a renderableHeap does not include any of its elements,
// as renderables cannot be copied.
func (rh *RenderableHeap) Copy() Stackable {
	return newHeap(rh.static)
}

func (rh *RenderableHeap) DrawToScreen(world draw.Image, viewPos *intgeom.Point2, screenW, screenH int) {
	newRh := &RenderableHeap{}
	if rh.static {
		var r Renderable
		// Undraws will all come first, loop to remove them
		for len(rh.rs) > 0 {
			r = rh.heapPop()
			if r.GetLayer() != Undraw {
				break
			}
		}
		for len(rh.rs) > 0 {
			r.Draw(world, 0, 0)
			newRh.heapPush(r)
			r = rh.heapPop()
		}
	} else {
		// TODO: test if we can remove these bounds checks (because draw.Draw already does them)
		vx := float64(-viewPos[0])
		vy := float64(-viewPos[1])
		for len(rh.rs) > 0 {
			r := rh.heapPop()
			if r.GetLayer() != Undraw {
				x2 := int(r.X())
				y2 := int(r.Y())
				w, h := r.GetDims()
				x := w + x2
				y := h + y2
				if x > viewPos[0] && y > viewPos[1] &&
					x2 < viewPos[0]+screenW && y2 < viewPos[1]+screenH {
					r.Draw(world, vx, vy)
				}
				newRh.heapPush(r)
			}
		}
	}
	rh.rs = newRh.rs
}

type layerHeap struct {
	rs []Renderable
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *layerHeap) heapPush(r Renderable) {
	if r == nil {
		return
	}
	h.rs = append(h.rs, r)
	h.up(len(h.rs) - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *layerHeap) heapPop() Renderable {
	n := len(h.rs) - 1
	h.rs[0], h.rs[n] = h.rs[n], h.rs[0]
	h.down(0, n)
	r := h.rs[len(h.rs)-1]
	h.rs = h.rs[0 : len(h.rs)-1]
	return r
}

func (h *layerHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.rs[i], h.rs[j] = h.rs[j], h.rs[i]
		j = i
	}
}

func (h *layerHeap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.rs[i], h.rs[j] = h.rs[j], h.rs[i]
		i = j
	}
	return i > i0
}

//Less returns whether a renderable at index i is at a lower layer than the one at index j
func (h *layerHeap) less(i, j int) bool {
	return h.rs[i].GetLayer() < h.rs[j].GetLayer()
}
