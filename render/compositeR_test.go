package render

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
)

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
	if !reflect.DeepEqual(cmp.Get(0), cb) {
		t.Fatalf("color box did not match color box in composite")
	}
	if cmp.GetRGBA() != nil {
		t.Fatalf("composite rgba was not nil")
	}

	cmp.Draw(image.NewRGBA(image.Rect(0, 0, 10, 10)), 0, 0)
	cmp.DrawToScreen(image.NewRGBA(image.Rect(0, 0, 3, 3)), intgeom.Point2{0, 0}, 100, 100)
	cmp.Undraw()
	if cmp.GetLayer() != Undraw {
		t.Fatalf("composite layer was not set to Undraw when undrawn")
	}

	cmp.Replace(cmp.Get(0), NewColorBox(1, 1, color.RGBA{0, 0, 128, 255}), 0)
	cmp.PreDraw()

	cmp2 := cmp.Copy().(*CompositeR)

	if len(cmp.rs) != len(cmp2.rs) {
		t.Fatalf("composite length mismatch post copy")
	}

	cmp2.DrawToScreen(image.NewRGBA(image.Rect(0, 0, 3, 3)), intgeom.Point2{0, 0}, 100, 100)

	cmp2.Prepend(nil)
	if cmp2.Len() != 1 {
		t.Fatalf("composite length was not increased by prepend")
	}
}

func TestCompositeR_Add(t *testing.T) {
	cmp := &CompositeR{}
	cmp.Add(EmptyRenderable())
	if len(cmp.toPush) != 1 {
		t.Fatalf("add did not add to composite r's toPush: %v", len(cmp.toPush))
	}
}
