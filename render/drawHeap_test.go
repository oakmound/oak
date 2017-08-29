package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/collision"
	"github.com/stretchr/testify/assert"
)

const heapLoops = 2000

func TestDrawHeapLoop(t *testing.T) {
	initTestFont()
	h := NewHeap(false)
	h2 := NewHeap(true)

	type toAdd struct {
		r     Renderable
		layer int
	}

	n := 3

	toAdds := []toAdd{
		{EmptyRenderable(), 0},
		{NewEmptySprite(20, 20, 10, 10), 1},
		{NewColorBox(30, 30, color.RGBA{255, 255, 255, 255}), 2},
		{NewSequence([]Modifiable{
			NewColorBox(5, 5, color.RGBA{0, 0, 0, 255}),
			NewColorBox(6, 6, color.RGBA{20, 0, 0, 255}),
			NewColorBox(7, 7, color.RGBA{40, 0, 0, 255}),
			NewColorBox(8, 9, color.RGBA{60, 0, 0, 255}),
		}, 10), 3},
		{NewScrollBox(
			[]Renderable{
				NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}),
				DefFont().NewStrText("test", 11, 11),
			},
			50,
			50,
			20,
			20,
		), 4},
		{DefFont().NewStrText("fire", 15, 15), 5},
		{DefFont().NewIntText(&n, 15, 15), 6},
		{DefFont().NewText(collision.NewUnassignedSpace(0, 0, 10, 10), 15, 15), 7},
	}
	// Scrollbox test that shouldn't really be here
	toAdds[4].r.(*ScrollBox).AddRenderable(DefFont().NewStrText("test", 14, 14))

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
