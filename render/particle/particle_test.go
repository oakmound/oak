package particle

import (
	"testing"

	"github.com/oakmound/oak/v2/render"
)

func TestParticle(t *testing.T) {
	var bp *baseParticle
	w, h := bp.GetDims()
	if w != 0 {
		t.Fatalf("expected 0 x, got %v", w)
	}
	if h != 0 {
		t.Fatalf("expected 0 y, got %v", h)
	}

	if bp.GetLayer() != render.Undraw {
		t.Fatalf("uninitialized particle was not set to the undraw layer")
	}

	bp = new(baseParticle)
	bp.setPID(100)
	if bp.pID != 100 {
		t.Fatalf("setPID failed, expected 100 got %v", bp.pID)
	}
}
