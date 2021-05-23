package particle

import (
	"testing"
)

func TestBaseGenerator(t *testing.T) {
	bg := new(BaseGenerator)
	bg.setDefaults()
	bg.ShiftX(10)
	bg.ShiftY(10)
	if bg.X() != 10 {
		t.Fatalf("expected 10 x, got %v", bg.X())
	}
	if bg.Y() != 10 {
		t.Fatalf("expected 10 y, got %v", bg.Y())
	}
}
