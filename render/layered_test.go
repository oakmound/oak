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
}
