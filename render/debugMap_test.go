package render

import (
	"image/color"
	"testing"
)

func TestDebugIdentity(t *testing.T) {
	r := NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	UpdateDebugMap("r", r)
	r2, ok := GetDebugRenderable("r")
	if !ok {
		t.Fatalf("debug renderable not set")
	}
	if r != r2 {
		t.Fatalf("got renderable did not match set renderable")
	}
	_, ok = GetDebugRenderable("doesn't exist")
	if ok {
		t.Fatalf("debug renderable string should not have existed")
	}
}
