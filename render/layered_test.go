package render

import (
	"testing"
)

func TestLayeredNils(t *testing.T) {
	var ld *Layer
	if ld.GetLayer() != Undraw {
		t.Fatalf("nil layer should be undrawn")
	}
	var ldp *LayeredPoint
	if ldp.GetLayer() != Undraw {
		t.Fatalf("nil layered point should be undrawn")
	}
	w, h := ldp.GetDims()
	if w != 1 || h != 1 {
		t.Fatalf("GetDims faailed")
	}
}
