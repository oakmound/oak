package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"container/heap"
	"golang.org/x/exp/shiny/screen"
	"image/color"
)

var (
	rh                *RenderableHeap
	toPushRenderables []Renderable
	postDrawBind      event.Binding
	bindingInit       bool
)

type RenderableHeap []Renderable

func (h RenderableHeap) Len() int           { return len(h) }
func (h RenderableHeap) Less(i, j int) bool { return h[i].GetLayer() < h[j].GetLayer() }
func (h RenderableHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *RenderableHeap) Push(x interface{}) {
	*h = append(*h, x.(Renderable))
}

func (h_p *RenderableHeap) Pop() interface{} {
	h := *h_p
	n := len(h)
	x := h[n-1]
	*h_p = h[0 : n-1]
	return x
}

func ResetDrawHeap() {
	InitDrawHeap()
}

func InitDrawHeap() {
	rh = &RenderableHeap{}
	heap.Init(rh)
	if bindingInit == false {
		postDrawBind, _ = event.GlobalBind(PostDraw, "PostDraw")
		bindingInit = true
	}
}

func Draw(r Renderable, l int) Renderable {
	// Bind to PostDraw if this causes synchronization issues with DrawHeap
	r.SetLayer(l)
	toPushRenderables = append(toPushRenderables, r)
	return r
}

func PostDraw(no int, nothing interface{}) error {
	for _, r := range toPushRenderables {
		heap.Push(rh, r)
	}
	toPushRenderables = []Renderable{}
	return nil
}

// For testing rectangle spaces
func DrawColor(c color.Color, x1, y1, x2, y2 float64, l int) {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	Draw(cb, l)
}

func LoadSpriteAndDraw(filename string, l int) *Sprite {
	s := LoadSprite(filename)
	return Draw(s, l).(*Sprite)
}

func DrawHeap(b screen.Buffer) {
	newRh := &RenderableHeap{}
	for rh.Len() > 0 {
		r := heap.Pop(rh).(Renderable)
		if r.GetLayer() != -1 {
			r.Draw(b)
			heap.Push(newRh, r)
		}
	}
	rh = newRh
}
