package render

// We manually define a LamdaHeap as it improves performance over using
// the stdlib's Heap interface.
type lambdaHeap struct {
	bh []Renderable
}

func (lh *lambdaHeap) Push(r Renderable) {
	if r == nil {
		return
	}
	lh.bh = append(lh.bh, r)
	lh.up(len(lh.bh) - 1)
}

func (lh *lambdaHeap) Pop() Renderable {
	n := len(lh.bh) - 1
	lh.bh[0], lh.bh[n] = lh.bh[n], lh.bh[0]

	lh.down(0, n)

	x := lh.bh[n]
	lh.bh = lh.bh[0:n]
	return x
}

func (lh *lambdaHeap) up(j int) {
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

func (lh *lambdaHeap) down(i, n int) {
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
