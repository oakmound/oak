package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

const heapLoops = 2000

func TestDrawHeapLoop(t *testing.T) {
	h := NewHeap(false)
	h2 := NewHeap(true)

	type toAdd struct {
		r     Renderable
		layer int
	}

	toAdds := []toAdd{
		{EmptyRenderable(), 0},
		{NewEmptySprite(20, 20, 10, 10), 1},
		{NewColorBox(30, 30, color.RGBA{255, 255, 255, 255}), 2},
	}

	for _, a := range toAdds {
		h.Add(a.r, a.layer)
		h2.Add(a.r, a.layer)
	}

	world := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	viewPos := image.Point{0, 0}

	for i := 0; i < heapLoops; i++ {
		h.PreDraw()
		h2.PreDraw()
		h.draw(world, viewPos, 640, 480)
		h2.draw(world, viewPos, 640, 480)
	}
}

func TestDrawHeapFns(t *testing.T) {
	h := NewHeap(false)
	h.Push(nil)
	assert.Empty(t, h.rs)
	h.Push(EmptyRenderable())
	h = h.Copy().(*RenderableHeap)
	assert.Empty(t, h.rs)

	h.Replace(EmptyRenderable(), NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}), 10)
}
