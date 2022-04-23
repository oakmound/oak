package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/v4/alg/intgeom"
)

const heapLoops = 2000

func TestDrawHeapLoop(t *testing.T) {
	h := NewDynamicHeap()
	h2 := NewStaticHeap()

	type toAdd struct {
		r     Renderable
		layer int
	}

	n := 3

	toAdds := []toAdd{
		{EmptyRenderable(), 0},
		{NewEmptySprite(20, 20, 10, 10), 1},
		{NewColorBox(30, 30, color.RGBA{255, 255, 255, 255}), 2},
		{NewSequence(10,
			NewColorBox(5, 5, color.RGBA{0, 0, 0, 255}),
			NewColorBox(6, 6, color.RGBA{20, 0, 0, 255}),
			NewColorBox(7, 7, color.RGBA{40, 0, 0, 255}),
			NewColorBox(8, 9, color.RGBA{60, 0, 0, 255})), 3},
		{DefaultFont().NewText("fire", 15, 15), 5},
		{DefaultFont().NewIntText(&n, 15, 15), 6},
	}
	for _, a := range toAdds {
		h.Add(a.r, a.layer)
		h2.Add(a.r, a.layer)
	}

	world := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	viewPos := intgeom.Point2{0, 0}

	for i := 0; i < heapLoops; i++ {
		h.PreDraw()
		h2.PreDraw()
		h.DrawToScreen(world, &viewPos, 640, 480)
		h2.DrawToScreen(world, &viewPos, 640, 480)
	}
}

func TestDrawHeapFns(t *testing.T) {
	h := NewDynamicHeap()
	h.heapPush(nil)
	if len(h.rs) != 0 {
		t.Fatalf("expected zero renderables in heap")
	}
	h.heapPush(EmptyRenderable())
	h = h.Copy().(*RenderableHeap)
	if len(h.rs) != 0 {
		t.Fatalf("expected zero renderables in copied heap")
	}

	h.Replace(EmptyRenderable(), NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}), 10)
	h.PreDraw()
}
