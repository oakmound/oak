package render

import (
	"image"
	"testing"

	"github.com/oakmound/oak/v4/event"
)

func TestLogicFPS(t *testing.T) {
	lfps := NewLogicFPS(0, nil, 0, 0)
	lfps.Draw(image.NewRGBA(image.Rect(0, 0, 100, 100)), 10, 10)
	logicFPSBind(lfps, event.EnterPayload{})
	if lfps.fps == 0 {
		t.Fatalf("fps not set by binding")
	}
}
