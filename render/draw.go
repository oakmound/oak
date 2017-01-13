package render

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	"container/heap"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	//"runtime"
	"time"
)

var (
	rh                *LambdaHeap
	srh               *RenderableHeap
	toPushRenderables []Renderable
	toPushStatic      []Renderable
	preDrawBind       event.Binding
	resetHeap         bool
	EmptyRenderable   = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
	//EmptyRenderable   = new(Composite)
)

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
	rh = &LambdaHeap{}
	srh = &RenderableHeap{}
	heap.Init(srh)
}

// Drawing does not actually immediately draw a renderable,
// instead the renderable is added to a list of elements to
// be drawn next frame. This avoids issues where elements
// are added to the heap while it is being drawn.
func Draw(r Renderable, l int) Renderable {
	r.SetLayer(l)
	toPushRenderables = append(toPushRenderables, r)
	return r
}

func StaticDraw(r Renderable, l int) Renderable {
	r.SetLayer(l)
	toPushStatic = append(toPushStatic, r)
	return r
}

// PreDraw parses through renderables to be pushed
// and adds them to the drawheap.
func PreDraw() {
	if resetHeap == true {
		InitDrawHeap()
		resetHeap = false
	} else {
		for _, r := range toPushRenderables {
			if r == nil {
				dlog.Warn("A nil was added to the draw heap")
				continue
			}
			rh.Push(r)
			//heap.Push(rh, r)
		}
		for _, r := range toPushStatic {
			heap.Push(srh, r)
		}
	}
	toPushStatic = []Renderable{}
	toPushRenderables = []Renderable{}
}

// LoadSpriteAndDraw is shorthand for LoadSprite
// followed by Draw.
func LoadSpriteAndDraw(filename string, l int) *Sprite {
	s := LoadSprite(filename)
	return Draw(s, l).(*Sprite)
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x1, y1, x2, y2 float64, l int) {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	Draw(cb, l)
}

// DrawHeap takes every element in the heap
// and draws them as it removes them. It
// filters out elements who have the layer
// -1, reserved for elements to be undrawn.
func DrawHeap(b screen.Buffer, vx, vy, screenW, screenH int) {
	drawRenderableHeap(b, rh, vx, vy, screenW, screenH)
}

func DrawStaticHeap(b screen.Buffer) {
	newRh := &RenderableHeap{}
	for srh.Len() > 0 {
		rp := heap.Pop(srh)
		if rp != nil {
			r := rp.(Renderable)
			if r.GetLayer() != -1 {
				r.Draw(b.RGBA())
				heap.Push(newRh, r)
			}
		}
	}
	*srh = *newRh
}

func drawRenderableHeap(b screen.Buffer, rheap *LambdaHeap, vx, vy, screenW, screenH int) {
	newRh := &LambdaHeap{}
	for len(rheap.bh) > 0 {
		r := rheap.Pop()
		if r != nil {
			if r.GetLayer() != -1 {
				x := int(r.GetX())
				y := int(r.GetY())
				x2 := x
				y2 := y
				rgba := r.GetRGBA()
				if rgba != nil {
					max := rgba.Bounds().Max
					x += max.X
					y += max.Y
					// Artificial width and height added due to bug in polygon checking alg
				} else {
					x += 6
					y += 6
				}
				if x > vx && y > vy &&
					x2 < vx+screenW && y2 < vy+screenH {

					if InDrawPolygon(x, y, x2, y2) {
						//if r.AlwaysDirty() || IsDirty(x2, y2) {
						r.Draw(b.RGBA())
					}
					//}
				}
				newRh.Push(r)
			}
		}
	}
	*rheap = *newRh
	dirtyZones = [DirtyZonesX][DirtyZonesY]bool{}
}

// ShinyDraw performs a draw operation at -x, -y, because
// shiny/screen represents quadrant 4 as negative in both axes.
// draw.Over will merge two pixels at a given position based on their
// alpha channel.
func ShinyDraw(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Over)
}

// draw.Src will overwrite pixels beneath the given image regardless of
// the new image's alpha.
func ShinyOverwrite(buff screen.Buffer, img image.Image, x, y int) {
	draw.Draw(buff.RGBA(), buff.Bounds(),
		img, image.Point{-x, -y}, draw.Src)
}

// UndrawAfter will trigger a renderable's undraw function
// after a given time has passed
func UndrawAfter(r Renderable, t time.Duration) {
	go func(r Renderable, t time.Duration) {
		select {
		case <-time.After(t):
			r.UnDraw()
		}
	}(r, t)
}

// DrawForTime is a wrapper for Draw and UndrawAfter
func DrawForTime(r Renderable, l int, t time.Duration) {
	Draw(r, l)
	UndrawAfter(r, t)
}
