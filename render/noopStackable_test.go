package render

import (
	"image"
	"testing"
)

func TestNoopStackable(t *testing.T) {
	noop := NoopStackable{}
	// these calls are noops
	noop.PreDraw()
	noop.Replace(nil, nil, -142)
	noop.draw(nil, image.Point{}, -124, 23)
	r := noop.Add(nil, 01, 124, 04, 2)
	if r != nil {
		t.Fatalf("expected nil renderable from Add, got %v", r)
	}
	noop2 := noop.Copy()
	if noop2 != noop {
		t.Fatalf("expected equal noop stackables")
	}
}
