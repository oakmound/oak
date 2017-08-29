package render

import (
	"image/color"
	"math"
	"testing"

	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
)

var (
	neqmods = []Modification{
		Brighten(10),
		CutRound(.5, .5),
		Fade(100),
		ApplyColor(color.RGBA{0, 255, 0, 255}),
		ColorBalance(100, 0, 0),
		ApplyMask(*NewColorBox(5, 5, color.RGBA{100, 100, 200, 200}).GetRGBA()),
		Rotate(10),
		Scale(2, 2),
		TrimColor(color.RGBA{255, 255, 255, 255}),
	}
	eqmods = []Modification{
		FlipX,
		FlipY,
	}
)

func TestRevertingMods(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{200, 0, 0, 255})
	rv := NewReverting(cb)
	assert.Equal(t, rv.GetRGBA(), cb.GetRGBA())
	for _, mod := range neqmods {
		rv.Modify(mod)
		assert.NotEqual(t, rv.GetRGBA(), cb.GetRGBA())
		rv.Revert(1)
		assert.Equal(t, rv.GetRGBA(), cb.GetRGBA())
	}
	rv = NewReverting(cb)
	for _, mod := range eqmods {
		rv.Modify(mod)
		assert.Equal(t, rv.GetRGBA(), cb.GetRGBA())
	}
	rv.RevertAll()
	assert.Equal(t, rv.GetRGBA(), cb.GetRGBA())

	rv.Modify(Scale(2, 2))
	rgba1 := rv.GetRGBA()
	rv = rv.Copy().(*Reverting)
	assert.Equal(t, rv.GetRGBA(), rgba1)
	rv.RevertAndModify(1, Scale(2, 2))
	assert.Equal(t, rv.GetRGBA(), rgba1)

	rv.Revert(math.MaxInt32)
	rv.Revert(-math.MaxInt32)
	rv.RevertAndModify(math.MaxInt32)
	rv.RevertAndModify(-math.MaxInt32)
	// Assert nothing went wrong with ^^
}

func TestRevertingCascadeFns(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{200, 0, 0, 255})
	rv := NewReverting(cb)

	// NOP
	// (color box does not have these functions)
	assert.True(t, rv.IsInterruptable())
	assert.True(t, rv.IsStatic())
	rv.Set("Foo")
	rv.Pause()
	rv.Unpause()
	rv.SetTriggerID(0)
	rv.update()

	sq := NewSequence([]Modifiable{
		cb, cb, cb,
	}, 1)

	cmpd := NewCompound("base",
		map[string]Modifiable{
			"base":  sq,
			"other": EmptyRenderable(),
		},
	)

	rv = NewReverting(cmpd)

	assert.Equal(t, sq.IsInterruptable(), rv.IsInterruptable())
	sq.Interruptable = !sq.Interruptable
	assert.Equal(t, sq.IsInterruptable(), rv.IsInterruptable())
	assert.Equal(t, sq.IsStatic(), rv.IsStatic())

	rv.Pause()
	assert.Equal(t, sq.playing, false)
	rv.Unpause()
	assert.Equal(t, sq.playing, true)
	rv.SetTriggerID(1)
	assert.Equal(t, sq.cID, event.CID(1))
	rv.update()

	rv.Set("other")

	rv.Pause()
	assert.Equal(t, sq.playing, true)
	rv.Unpause()
	assert.Equal(t, sq.playing, true)
	rv.update()
}
