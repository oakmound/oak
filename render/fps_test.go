package render

import (
	"image"
	"testing"
)

func TestDrawFPS(t *testing.T) {
	initTestFont()
	dfps := NewDrawFPS(0, nil, 0, 0)
	dfps.Draw(image.NewRGBA(image.Rect(0, 0, 100, 100)), 10, 10)
	if dfps.fps == 0 {
		t.Fatalf("fps not set by draw")
	}
}
