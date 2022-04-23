package render

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/oakmound/oak/v4/physics"
	"github.com/oakmound/oak/v4/render/mod"
)

func TestCompoundFuncs(t *testing.T) {
	swtch := NewSwitch("red", map[string]Modifiable{
		"red":   NewColorBox(5, 5, color.RGBA{255, 0, 0, 255}),
		"blue":  NewColorBox(5, 5, color.RGBA{0, 0, 255, 255}),
		"green": NewColorBox(5, 5, color.RGBA{0, 255, 0, 255}),
		"empty": NewColorBox(5, 5, color.RGBA{0, 0, 0, 0}),
	})
	if swtch.GetRGBA().At(0, 0) != (color.RGBA{255, 0, 0, 255}) {
		t.Fatalf("red was not red")
	}
	if swtch.Get() != "red" {
		t.Fatalf("switch did not begin in red state")
	}
	if swtch.Set("blue") != nil {
		t.Fatalf("setting switch to blue should not fail")
	}
	if swtch.GetRGBA().At(0, 0) != (color.RGBA{0, 0, 255, 255}) {
		t.Fatalf("blue was not blue")
	}
	if swtch.Get() != ("blue") {
		t.Fatalf("set blue did not set key")
	}
	if swtch.Add("blue", NewColorBox(5, 5, color.RGBA{255, 255, 255, 255})) == nil {
		t.Fatalf("switch should return erroron duplicate Add")
	}
	if swtch.GetRGBA().At(0, 0) != (color.RGBA{255, 255, 255, 255}) {
		t.Fatalf("add blue did not also set blue on switch")
	}
	if swtch.Set("color5") == nil {
		t.Fatalf("switch should error setting to unknown string key")
	}
	if !reflect.DeepEqual(swtch.GetSub("empty"), NewColorBox(5, 5, color.RGBA{0, 0, 0, 0})) {
		t.Fatalf("empty renderable was not a 5x5 colorbox")
	}

	swtch2 := swtch.Copy().(*Switch)
	w, h := swtch2.GetDims()
	if w != 5 || h != 5 {
		t.Fatalf("get dims failed")
	}
	if !swtch2.IsStatic() {
		t.Fatalf("switch on colorbox should be static")
	}
	if !swtch2.IsInterruptable() {
		t.Fatalf("switch on colorbox should be interruptable")
	}
}

func TestSwitchPositioning(t *testing.T) {
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	cmp := NewSwitch("red", map[string]Modifiable{
		"red":  NewColorBox(5, 5, red),
		"blue": NewColorBox(5, 5, blue),
	})
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))

	cmp.Draw(rgba, 0, 0)
	if red != rgba.At(0, 0) {
		t.Fatalf("rgba was not red at 0,0")
	}
	if red == rgba.At(10, 0) {
		t.Fatalf("rgba was red at 10,0")
	}

	rgba = image.NewRGBA(image.Rect(0, 0, 10, 10))
	cmp.Set("blue")
	cmp.ShiftPos(1, 1)
	cmp.Draw(rgba, 0, 0)
	if blue != rgba.At(1, 1) {
		t.Fatalf("rgba was not blue at 1,1")
	}
	if blue == rgba.At(0, 0) {
		t.Fatalf("shifting rgba should not have drawn at 0,0")
	}
	rgba = image.NewRGBA(image.Rect(0, 0, 10, 10))
	cmp.SetOffsets("red", physics.NewVector(5, 5))
	cmp.Set("red")
	cmp.Draw(rgba, 0, 0)
	if red != rgba.At(8, 8) {
		t.Fatalf("red should be set at 8,8 with SetOffsets")
	}
	if red == rgba.At(1, 1) {
		t.Fatalf("set offsets should have unset 1,1")
	}

}

func TestSwitchModifiability(t *testing.T) {
	cb := NewColorBox(5, 5, color.RGBA{255, 0, 0, 255})
	cmp := NewSwitch("red", map[string]Modifiable{
		"red":  NewReverting(NewColorBox(5, 5, color.RGBA{255, 0, 0, 255})),
		"blue": NewReverting(NewColorBox(5, 5, color.RGBA{0, 0, 255, 255})),
	})

	if !reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
		t.Fatalf("rgba mismatch")
	}
	for _, mod := range neqmods {
		cmp.Modify(mod)
		if reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
			t.Fatalf("neq mod should result in changed rgba")
		}
		cmp.Revert(1)
		if !reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
			t.Fatalf("rgba mismatch")
		}
	}
	for _, mod := range eqmods {
		cmp.Modify(mod)
		if !reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
			t.Fatalf("rgba mismatch")
		}
	}
	cmp.RevertAll()
	if !reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
		t.Fatalf("rgba mismatch")
	}

	cmp.Filter(mod.Brighten(-100))
	if reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
		t.Fatalf("brighten filter should result in changed rgba")
	}
	cmp.RevertAll()
	if !reflect.DeepEqual(cb.GetRGBA(), cmp.GetRGBA()) {
		t.Fatalf("rgba mismatch")
	}

}
