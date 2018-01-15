package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render/mod"
	"github.com/stretchr/testify/assert"
)

func TestComposite(t *testing.T) {
	cmp := NewComposite(
		NewColorBox(5, 5, color.RGBA{255, 0, 0, 255}),
		NewColorBox(3, 3, color.RGBA{0, 255, 0, 255}),
	)
	cmp.Append(NewColorBox(6, 6, color.RGBA{0, 0, 255, 255}))
	cmp.SetIndex(0, NewColorBox(4, 4, color.RGBA{255, 255, 255, 255}))
	cmp.SetOffsets(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{1, 1},
		floatgeom.Point2{2, 2},
	)
	cmp.AddOffset(0, floatgeom.Point2{4, 4})
	cb := NewColorBox(4, 4, color.RGBA{255, 255, 255, 255})
	cb.SetPos(4, 4)
	assert.Equal(t, cb, cmp.Get(0))
	assert.Nil(t, cmp.GetRGBA())

	cmp.Draw(image.NewRGBA(image.Rect(0, 0, 10, 10)))
	cmp.DrawOffset(image.NewRGBA(image.Rect(0, 0, 10, 10)), 5, 5)
	cmp.Undraw()
	assert.Equal(t, Undraw, cmp.GetLayer())
	cmp2 := cmp.Copy()
	cmp2.Filter(mod.Brighten(-100))
	cmp2.Modify(mod.Scale(.5, .5))
	cmp3 := NewComposite(
		NewColorBox(2, 2, color.RGBA{0, 0, 0, 255}),
		NewColorBox(1, 1, color.RGBA{0, 0, 0, 255}),
		NewColorBox(3, 3, color.RGBA{0, 0, 0, 255}),
	)
	cmp3.Undraw()
	cmp3.SetOffsets(
		floatgeom.Point2{4, 4},
		floatgeom.Point2{1, 1},
		floatgeom.Point2{2, 2},
	)
	assert.Equal(t, cmp2, cmp3)

	cmp3.Prepend(nil)
	assert.Equal(t, 4, cmp3.Len())
}

func TestCompositeR(t *testing.T) {
	cmp := NewCompositeR(
		NewColorBox(5, 5, color.RGBA{255, 0, 0, 255}),
		NewColorBox(3, 3, color.RGBA{0, 255, 0, 255}),
	)
	cmp.SetIndex(0, NewColorBox(4, 4, color.RGBA{255, 255, 255, 255}))
	cmp.SetOffsets(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{1, 1},
	)
	cmp.AppendOffset(
		NewColorBox(6, 6, color.RGBA{0, 0, 255, 255}),
		floatgeom.Point2{2, 2},
	)
	cmp.AddOffset(0, floatgeom.Point2{4, 4})
	cb := NewColorBox(4, 4, color.RGBA{255, 255, 255, 255})
	cb.SetPos(4, 4)
	assert.Equal(t, cb, cmp.Get(0))
	assert.Nil(t, cmp.GetRGBA())

	cmp.Draw(image.NewRGBA(image.Rect(0, 0, 10, 10)))
	cmp.DrawOffset(image.NewRGBA(image.Rect(0, 0, 10, 10)), 5, 5)
	cmp.draw(image.NewRGBA(image.Rect(0, 0, 3, 3)), image.Point{0, 0}, 100, 100)
	cmp.Undraw()
	assert.Equal(t, Undraw, cmp.GetLayer())

	cmp.Replace(cmp.Get(0), NewColorBox(1, 1, color.RGBA{0, 0, 128, 255}), 0)

	cmp2 := cmp.Copy().(*CompositeR)

	assert.Equal(t, len(cmp.rs), len(cmp2.rs))

	cmp2.draw(image.NewRGBA(image.Rect(0, 0, 3, 3)), image.Point{0, 0}, 100, 100)

	cmp2.Prepend(nil)
	assert.Equal(t, 1, cmp2.Len())
}
