package render

import (
	"image/color"
	"math"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v3/render/mod"
)

var (
	neqmods = []mod.Mod{
		mod.CutRound(.5, .5),
		mod.Rotate(10),
		mod.Scale(2, 2),
		mod.TrimColor(color.RGBA{255, 255, 255, 255}),
	}
	neqfilters = []mod.Filter{mod.Brighten(-100)}
	eqmods     = []mod.Mod{
		mod.FlipX,
		mod.FlipY,
	}
)

func TestRevertingMods(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{200, 0, 0, 255})
	rv := NewReverting(cb)
	if rv.GetRGBA() != cb.GetRGBA() {
		t.Fatalf("reverting did not call underlying GetRGBA")
	}
	for _, mod := range neqmods {
		rv.Modify(mod)
		if reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
			t.Fatalf("neq mod did not change rgba")
		}
		rv.Revert(1)
		if !reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
			t.Fatalf("revert did not match rgba")
		}
	}
	rv = NewReverting(cb)
	for _, mod := range eqmods {
		rv.Modify(mod)
		if !reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
			t.Fatalf("eq mods did not match rgba")
		}
	}
	rv.RevertAll()
	if !reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
		t.Fatalf("revert-all did not match rgba")
	}

	rv.Modify(mod.Scale(2, 2))
	rgba1 := rv.GetRGBA()
	rv = rv.Copy().(*Reverting)
	if !reflect.DeepEqual(rv.GetRGBA(), rgba1) {
		t.Fatalf("copied rgba did not match")
	}
	rv.RevertAndModify(1, mod.Scale(2, 2))
	if !reflect.DeepEqual(rv.GetRGBA(), rgba1) {
		t.Fatalf("revert and scaled rgba did not match")
	}

	rv.Revert(math.MaxInt32)
	rv.Revert(-math.MaxInt32)
	rv.RevertAndModify(math.MaxInt32)
	rv.RevertAndModify(-math.MaxInt32)
	// Assert nothing crashed
}

func TestRevertingFilters(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{200, 0, 0, 255})
	rv := NewReverting(cb)
	for _, f := range neqfilters {
		rv.Filter(f)
		if reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
			t.Fatalf("neq filter did not change rgba")
		}
		rv.Revert(1)
		if !reflect.DeepEqual(rv.GetRGBA(), cb.GetRGBA()) {
			t.Fatalf("revert did not match rgba")
		}
	}
	rv.Filter(mod.Brighten(-100))
	rgba1 := rv.GetRGBA()
	rv = rv.Copy().(*Reverting)
	if !reflect.DeepEqual(rv.GetRGBA(), rgba1) {
		t.Fatalf("copied rgba did not match")
	}
	if len(rv.rs) != 2 {
		t.Fatalf("expected 2 renderables in reverting stack, got %v", len(rv.rs))
	}
	rv2 := rv.RevertAndFilter(1, mod.Brighten(-100))
	rv = rv2.Copy().(*Reverting)
	if !reflect.DeepEqual(rv.GetRGBA(), rgba1) {
		t.Fatalf("copied rgba did not match")
	}
	if len(rv.rs) != 2 {
		t.Fatalf("expected 2 renderables in reverting stack, got %v", len(rv.rs))
	}

	rv.RevertAndFilter(math.MaxInt32)
	rv.RevertAndFilter(-math.MaxInt32)
}
func TestRevertingCascadeFns(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{200, 0, 0, 255})
	rv := NewReverting(cb)

	// NOP
	// (color box does not have these functions)
	if !rv.IsInterruptable() {
		t.Fatalf("colorbox is interruptable")
	}
	if !rv.IsStatic() {
		t.Fatalf("colorbox is static")
	}
	if rv.Set("Foo") != nil {
		t.Fatalf("colorbox.Set(foo) should return nil")
	}
	if rv.Get() != "" {
		t.Fatalf("get on a nil reverting should return an empty string")
	}
	rv.Pause()
	rv.Unpause()
	rv.SetTriggerID(0)
	rv.update()

	sq := NewSequence(1, cb, cb, cb)

	cmpd := NewSwitch("base",
		map[string]Modifiable{
			"base":  sq,
			"other": EmptyRenderable(),
		},
	)

	rv = NewReverting(cmpd)

	if sq.IsInterruptable() != rv.IsInterruptable() {
		t.Fatalf("sequence interruptable did not match reverting")
	}
	sq.Interruptable = !sq.Interruptable
	if sq.IsInterruptable() != rv.IsInterruptable() {
		t.Fatalf("sequence interruptable did not match reverting")
	}
	if sq.IsStatic() != rv.IsStatic() {
		t.Fatalf("sequence static did not match reverting")
	}

	rv.Pause()
	if sq.playing {
		t.Fatalf("reverting pause did not pause underlying sequence")
	}
	rv.Unpause()
	if !sq.playing {
		t.Fatalf("reverting unpause did not resume underlying sequence")
	}
	rv.SetTriggerID(1)
	if sq.CallerID != 1 {
		t.Fatalf("reverting cID did not set underlying squence cID")
	}
	rv.update()

	if rv.Set("other") != nil {
		t.Fatalf("switch.Set(foo) should return nil")
	}
	if rv.Set("notincompound") == nil {
		t.Fatalf("switch.Set(notincompound) should return an error")
	}

	if rv.Get() != "other" {
		t.Fatalf("set did not set key to other")
	}

	rv.Pause()
	if !sq.playing {
		t.Fatalf("reverting pause should not pause underlying sequence when switch is set to other")
	}
	rv.Unpause()
	if !sq.playing {
		t.Fatalf("reverting pause should not pause underlying sequence when switch is set to other")
	}
	rv.update()
}
