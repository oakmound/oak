package render

import (
	"image/color"
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/v2/fileutil"
)

func ExampleDraw() {
	// We haven't modified the draw stack, so it contains a single draw heap.
	// Draw a Color Box
	Draw(NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}), 3)
	// Draw a Gradient Box above that color box
	Draw(NewHorizontalGradientBox(5, 5, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}), 4)
}

func TestDrawHelpers(t *testing.T) {
	r, err := LoadSpriteAndDraw("Not a sprite", 0)
	if err == nil {
		t.Fatalf("load sprite should have failed")
	}
	if r != nil {
		t.Fatalf("failed load sprite should give nil renderable")
	}

	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	r, err = LoadSpriteAndDraw(filepath.Join("16", "jeremy.png"), 0)
	if err != nil {
		t.Fatalf("load sprite should not have failed")
	}
	if r == nil {
		t.Fatalf("load sprite should not give nil renderable")
	}

	r, err = DrawColor(color.RGBA{255, 255, 255, 255}, 0, 0, 10, 10, 0, 0)
	if err != nil {
		t.Fatalf("draw color should not have failed")
	}
	if r == nil {
		t.Fatalf("draw color should not give nil renderable")
	}
	GlobalDrawStack.Push(&CompositeR{})
	GlobalDrawStack.PreDraw()

	_, err = DrawColor(color.RGBA{255, 255, 255, 255}, 0, 0, 10, 10, 3, 0)
	if err == nil {
		t.Fatalf("draw color to invalid layer should fail")
	}

	err = DrawForTime(NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 0, 4)
	if err == nil {
		t.Fatalf("draw time to invalid layer should fail")
	}

	err = DrawForTime(NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 0, 0)
	if err != nil {
		t.Fatalf("draw time should not have failed")
	}

	_, err = DrawPoint(color.RGBA{100, 100, 100, 255}, 0, 0, 0)
	if err != nil {
		t.Fatalf("draw color should not have failed")
	}
}
