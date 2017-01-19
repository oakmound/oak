package render

import "container/heap"

type RenderableHeap []Renderable

// Satisfying the Heap interface
func (h RenderableHeap) Len() int           { return len(h) }
func (h RenderableHeap) Less(i, j int) bool { return h[i].GetLayer() < h[j].GetLayer() }
func (h RenderableHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *RenderableHeap) Push(x interface{}) {
	if x == nil {
		return
	}
	*h = append(*h, x.(Renderable))
}

func (h_p *RenderableHeap) Pop() interface{} {
	h := *h_p
	n := len(h)
	x := h[n-1]
	*h_p = h[0 : n-1]
	return x
}

// ResetDrawHeap sets a flag to clear the drawheap
// at the next predraw phase
func ResetDrawHeap() {
	resetHeap = true
}

func InitDrawHeap() {
	rh = &RenderableHeap{}
	srh = &RenderableHeap{}
	heap.Init(rh)
	heap.Init(srh)
}

// We manually define a LamdaHeap as it improves performance over using
// the stdlib's Heap interface.
type LambdaHeap struct {
	bh []Renderable
}

func (lh *LambdaHeap) Push(r Renderable) {
	if r == nil {
		return
	}
	lh.bh = append(lh.bh, r)
	lh.up(len(lh.bh) - 1)
}

func (lh *LambdaHeap) Pop() Renderable {
	n := len(lh.bh) - 1
	lh.bh[0], lh.bh[n] = lh.bh[n], lh.bh[0]

	lh.down(0, n)

	x := lh.bh[n]
	lh.bh = lh.bh[0:n]
	return x
}

func (lh *LambdaHeap) up(j int) {
	h := lh.bh
	var i int
	for {
		i = (j - 1) / 2 // parent
		if i == j || !(h[j].GetLayer() < h[i].GetLayer()) {
			break
		}
		lh.bh[i], lh.bh[j] = h[j], h[i]
		h = lh.bh
		j = i
	}
}

func (lh *LambdaHeap) down(i, n int) {
	h := lh.bh
	for {
		j1 := 2*i + 1
		if j1 >= n { // j1 < 0 after int overflow, ignored
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !(h[j1].GetLayer() < h[j2].GetLayer()) {
			j = j2 // = 2*i + 2  // right child
		}
		if !(h[j].GetLayer() < h[i].GetLayer()) {
			break
		}
		lh.bh[i], lh.bh[j] = h[j], h[i]
		h = lh.bh
		i = j
	}
}
