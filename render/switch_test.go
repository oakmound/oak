package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render/mod"
	"github.com/stretchr/testify/assert"
)

func TestCompoundFuncs(t *testing.T) {
	cmp := NewSwitch("red", map[string]Modifiable{
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
	assert.NotNil(t, cmp.Add("blue", NewColorBox(5, 5, color.RGBA{255, 255, 255, 255})))
	assert.Equal(t, cmp.GetRGBA().At(0, 0), color.RGBA{255, 255, 255, 255})
	assert.NotNil(t, cmp.Set("color5"))
	assert.Equal(t, cmp.GetSub("empty"), NewColorBox(5, 5, color.RGBA{0, 0, 0, 0}))

	cmp2 := cmp.Copy().(*Switch)
	w, h := cmp2.GetDims()
	assert.Equal(t, w, 5)
	assert.Equal(t, h, 5)
	assert.True(t, cmp2.IsStatic())
	assert.True(t, cmp2.IsInterruptable())
}

func TestSwitchPositioning(t *testing.T) {
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	cmp := NewSwitch("red", map[string]Modifiable{
		"red":  NewColorBox(5, 5, red),
		"blue": NewColorBox(5, 5, blue),
	})
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))

	cmp.Draw(rgba)
	assert.Equal(t, red, rgba.At(0, 0))
	assert.NotEqual(t, red, rgba.At(10, 0))

	rgba = image.NewRGBA(image.Rect(0, 0, 10, 10))
	assert.Nil(t, cmp.Set("blue"))
	cmp.ShiftPos(1, 1)
	cmp.DrawOffset(rgba, 0, 0)
	assert.Equal(t, blue, rgba.At(1, 1))
	assert.NotEqual(t, blue, rgba.At(0, 0))
	rgba = image.NewRGBA(image.Rect(0, 0, 10, 10))
	cmp.SetOffsets("red", physics.NewVector(5, 5))
	assert.Nil(t, cmp.Set("red"))
	cmp.Draw(rgba)
	assert.Equal(t, red, rgba.At(8, 8))
	assert.NotEqual(t, red, rgba.At(1, 1))

}

func TestSwitchModifiability(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{255, 0, 0, 255})
	cmp := NewSwitch("red", map[string]Modifiable{
		"red":  NewReverting(NewColorBox(5, 5, color.RGBA{255, 0, 0, 255})),
		"blue": NewReverting(NewColorBox(5, 5, color.RGBA{0, 0, 255, 255})),
	})

	assert.Equal(t, cb.GetRGBA(), cmp.GetRGBA())
	for _, mod := range neqmods {
		cmp.Modify(mod)
		assert.NotEqual(t, cb.GetRGBA(), cmp.GetRGBA())
		cmp.Revert(1)
		assert.Equal(t, cb.GetRGBA(), cmp.GetRGBA())
	}
	for _, mod := range eqmods {
		cmp.Modify(mod)
		assert.Equal(t, cb.GetRGBA(), cmp.GetRGBA())
	}
	cmp.RevertAll()
	assert.Equal(t, cb.GetRGBA(), cmp.GetRGBA())

	cmp.Filter(mod.Brighten(-100))
	assert.NotEqual(t, cb.GetRGBA(), cmp.GetRGBA())
	cmp.RevertAll()
	assert.Equal(t, cb.GetRGBA(), cmp.GetRGBA())

}
