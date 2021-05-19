package render

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

func TestDrawStack(t *testing.T) {
	GlobalDrawStack.PreDraw()
	if len(GlobalDrawStack.as) != 1 {
		t.Fatalf("global draw stack did not have one length initially")
	}
	SetDrawStack(
		NewDynamicHeap(),
		NewStaticHeap(),
	)
	if len(GlobalDrawStack.as) != 2 {
		t.Fatalf("global draw stack did not have two length after reset")
	}
	GlobalDrawStack.Pop()
	GlobalDrawStack.PreDraw()
	if len(GlobalDrawStack.as) != 1 {
		t.Fatalf("global draw stack did not have one length after pop")
	}
}

func TestDrawStack_Draw(t *testing.T) {
	_, err := Draw(nil)
	if err == nil {
		t.Fatalf("draw(nil) should have failed")
	}
	cb := NewColorBox(10, 10, color.RGBA{0, 0, 255, 255})
	Draw(cb)
	GlobalDrawStack.Clear()
	SetDrawStack(
		NewDynamicHeap(),
		NewStaticHeap(),
	)
	Draw(cb)
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))
	GlobalDrawStack.PreDraw()
	GlobalDrawStack.DrawToScreen(rgba, intgeom.Point2{0, 0}, 10, 10)
	if !reflect.DeepEqual(rgba, cb.GetRGBA()) {
		t.Fatalf("rgba mismatch")
	}
}
