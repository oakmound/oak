package render

import (
	"container/heap"
	"image"
	"image/color"
	"image/draw"

	"time"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/timing"
)

var (
	rh                *RenderableHeap
	srh               *RenderableHeap
	toPushRenderables []Renderable
	toPushStatic      []Renderable
	preDrawBind       event.Binding
	resetHeap         bool
	EmptyRenderable   = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
	//EmptyRenderable   = new(Composite)
)

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
	i := 0
	// defer func() {
	// 	if x := recover(); x != nil {
	// 		dlog.Error("Invalid Memory Address in toPushRenderables")
	// 		// This does not work-- all addresses following the bad address
	// 		// at i are also bad
	// 		//toPushRenderables = toPushRenderables[i+1:]
	// 		toPushRenderables = []Renderable{}
	// 	}
	// }()
	if resetHeap == true {
		InitDrawHeap()
		resetHeap = false
	} else {
		for _, r := range toPushRenderables {
			if r != nil {
				heap.Push(rh, r)
			}
			i++
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
func DrawHeap(target *image.RGBA, viewPos image.Point, screenW, screenH int) {
	drawRenderableHeap(target, rh, viewPos, screenW, screenH)
}

func DrawStaticHeap(target *image.RGBA) {
	newRh := &RenderableHeap{}
	for srh.Len() > 0 {
		rp := heap.Pop(srh)
		if rp != nil {
			r := rp.(Renderable)
			if r.GetLayer() != -1 {
				r.Draw(target)
				heap.Push(newRh, r)
			}
		}
	}
	*srh = *newRh
}

func drawRenderableHeap(target *image.RGBA, rheap *RenderableHeap, viewPos image.Point, screenW, screenH int) {
	newRh := &RenderableHeap{}
	for rheap.Len() > 0 {
		intf := heap.Pop(rheap)
		if intf != nil {
			r := intf.(Renderable)
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
				if x > viewPos.X && y > viewPos.Y &&
					x2 < viewPos.X+screenW && y2 < viewPos.Y+screenH {

					if InDrawPolygon(x, y, x2, y2) {
						r.Draw(target)
					}
				}
				heap.Push(newRh, r)
			}
		}
	}
	*rheap = *newRh
}

// UndrawAfter will trigger a renderable's undraw function
// after a given time has passed
func UndrawAfter(r Renderable, t time.Duration) {
	go func(r Renderable, t time.Duration) {
		timing.DoAfter(t, func() {
			r.UnDraw()
		})
	}(r, t)
}

// DrawForTime is a wrapper for Draw and UndrawAfter
func DrawForTime(r Renderable, l int, t time.Duration) {
	Draw(r, l)
	UndrawAfter(r, t)
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
func ShinyOverwrite(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Src)
}
