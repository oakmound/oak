package render

import (
	"image/color"
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/fileutil"
	"github.com/stretchr/testify/assert"
)

func DrawExample() {
	// We haven't modified the draw stack, so it contains a single draw heap.
	// Draw a Color Box
	Draw(NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}), 3)
	// Draw a Gradient Box above that color box
	Draw(NewHorizontalGradientBox(5, 5, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}), 4)
}

func TestDrawHelpers(t *testing.T) {
	r, err := LoadSpriteAndDraw("Not a sprite", 0)
	assert.Nil(t, r)
	assert.NotNil(t, err)

	fileutil.BindataDir = AssetDir
	fileutil.BindataFn = Asset
	r, err = LoadSpriteAndDraw(filepath.Join("16", "jeremy.png"), 0)
	assert.NotNil(t, r)
	assert.Nil(t, err)

	r, err = DrawColor(color.RGBA{255, 255, 255, 255}, 0, 0, 10, 10, 0, 0)
	assert.NotNil(t, r)
	assert.Nil(t, err)

	GlobalDrawStack.Push(&CompositeR{})
	GlobalDrawStack.PreDraw()

	_, err = DrawColor(color.RGBA{255, 255, 255, 255}, 0, 0, 10, 10, 0, 3)
	assert.NotNil(t, err)

	err = DrawForTime(NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 4, 0)
	assert.NotNil(t, err)

	err = DrawForTime(NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 0, 0)
	assert.Nil(t, err)
}
