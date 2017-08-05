package render

import (
	"image/color"
	"testing"

	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"
)

func TestCompoundFuncs(t *testing.T) {
	cmp := NewCompound("red", map[string]Modifiable{
		"red":   NewColorBox(5, 5, color.RGBA{255, 0, 0, 255}),
		"blue":  NewColorBox(5, 5, color.RGBA{0, 0, 255, 255}),
		"green": NewColorBox(5, 5, color.RGBA{0, 255, 0, 255}),
		"empty": NewColorBox(5, 5, color.RGBA{0, 0, 0, 0}),
	})
	assert.Equal(t, cmp.GetRGBA().At(0, 0), color.RGBA{255, 0, 0, 255})
	assert.Equal(t, cmp.Get(), "red")
	assert.Nil(t, cmp.Set("blue"))
	assert.Equal(t, cmp.GetRGBA().At(0, 0), color.RGBA{0, 0, 255, 255})
	assert.Equal(t, cmp.Get(), "blue")
	cmp.Add("blue", NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}))
	assert.Equal(t, cmp.GetRGBA().At(0, 0), color.RGBA{255, 255, 255, 255})
	assert.NotNil(t, cmp.Set("color5"))
	assert.Equal(t, cmp.GetSub("empty"), NewColorBox(5, 5, color.RGBA{0, 0, 0, 0}))
	cmp.SetOffsets("Green", physics.NewVector(5, 5))
	cmp2 := cmp.Copy().(*Compound)
	w, h := cmp2.GetDims()
	assert.Equal(t, w, 5)
	assert.Equal(t, h, 5)
	assert.True(t, cmp2.IsStatic())
	assert.True(t, cmp2.IsInterruptable())
}
