package oak

import (
	"testing"
)

func TestAspectRatio(t *testing.T) {
	c1 := blankScene(t)
	c1.SetAspectRatio(2)
	c1.ChangeWindow(10, 10)
	w := c1.windowRect.Max.X - c1.windowRect.Min.X
	h := c1.windowRect.Max.Y - c1.windowRect.Min.Y
	if w != 10 {
		t.Fatalf("height was not 10, got %v", w)
	}
	if h != 5 {
		t.Fatalf("height was not 5, got %v", h)
	}
	c1.ChangeWindow(10, 2)
	w = c1.windowRect.Max.X - c1.windowRect.Min.X
	h = c1.windowRect.Max.Y - c1.windowRect.Min.Y
	if w != 4 {
		t.Fatalf("height was not 4, got %v", w)
	}
	if h != 2 {
		t.Fatalf("height was not 2, got %v", h)
	}
}
