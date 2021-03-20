package render

import (
	"image"
	"testing"
)

func TestLogicFPS(t *testing.T) {
	initTestFont()
	lfps := NewLogicFPS(0.1, nil, 0, 0)
	lfps.Draw(image.NewRGBA(image.Rect(0, 0, 100, 100)), 10, 10)
	logicFPSBind(lfps.CID, nil)
	if lfps.fps == 0 {
		t.Fatalf("fps not set by binding")
	}
}
