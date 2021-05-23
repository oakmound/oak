package render

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/render/mod"
)

func TestComposite(t *testing.T) {
	cmp := NewCompositeM(
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
	if !reflect.DeepEqual(cmp.Get(0), cb) {
		t.Fatalf("color box did not match color box in composite")
	}
	if cmp.GetRGBA() != nil {
		t.Fatalf("composite rgba was not nil")
	}

	cmp.Draw(image.NewRGBA(image.Rect(0, 0, 10, 10)), 0, 0)
	cmp.Undraw()
	if cmp.GetLayer() != Undraw {
		t.Fatalf("composite layer was not set to Undraw when undrawn")
	}
	cmp2 := cmp.Copy()
	cmp2.Filter(mod.Brighten(-100))
	cmp2.Modify(mod.Scale(.5, .5))
	cmp3 := NewCompositeM(
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
	if !reflect.DeepEqual(cmp2, cmp3) {
		t.Fatalf("composites did not match")
	}

	cmp3.Prepend(nil)
	if cmp3.Len() != 4 {
		t.Fatalf("composite length was not increased by prepend")
	}

	cmp4 := NewCompositeM(
		NewColorBox(2, 2, color.RGBA{0, 0, 0, 255}),
		NewColorBox(1, 1, color.RGBA{0, 0, 0, 255}),
		NewColorBox(3, 3, color.RGBA{0, 0, 0, 255}),
	)

	cSprite := cmp4.ToSprite()
	if (color.RGBA{0, 0, 0, 255}) != cSprite.At(1, 1).(color.RGBA) {
		t.Fatalf("composite did not combine 0,0,0,0 3 times into 0,0,0,0")
	}
}

func TestCompositeM_Slice(t *testing.T) {
	cmp := NewCompositeM(
		EmptyRenderable(),
		EmptyRenderable(),
		EmptyRenderable(),
		EmptyRenderable(),
		EmptyRenderable(),
	)
	cmp2 := cmp.Slice(0, 2)
	if len(cmp2.rs) != 2 {
		t.Fatalf("composite slice did not reduce rs count")
	}
	cmp2 = cmp.Slice(-1, 1000)
	if len(cmp2.rs) != 5 {
		t.Fatalf("composite slice did not adjust when given invalid input")
	}
}
