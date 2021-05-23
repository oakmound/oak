package shape

import (
	"reflect"
	"testing"
)

func TestRect(t *testing.T) {
	shapes := []Shape{Square, Rectangle, Diamond, Circle, Checkered, Heart, JustIn(NotIn(Diamond.In))}

	w, h := 10, 10
	for k, s := range shapes {
		r := InToRect(s.In)(w, h)
		if !reflect.DeepEqual(r, s.Rect(w, h)) {
			t.Fatalf("Shape %v's InToRect did not match s.Rect", k)
		}
		for i := 0; i < w; i++ {
			for j := 0; j < h; j++ {
				if r[i][j] != s.In(i, j, w, h) {
					t.Fatalf("Shape %v's In at (%d,%d) did not match s.Rect", k, i, j)
				}
			}
		}
	}
}

func TestRectangleIn(t *testing.T) {
	if Rectangle.In(10, 10, 5, 5) {
		t.Fatal("10,10 should not be in a 5x5 rectangle")
	}
}

func TestStrictRect(t *testing.T) {
	sr := NewStrictRect(5, 5)
	for x := 0; x < 6; x++ {
		for y := 0; y < 6; y++ {
			if sr.In(x, y) {
				t.Fatalf("StrictRect.In was not completely false at %d,%d", x, y)
			}
		}
	}
	r := sr.Rect()
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			if r[x][y] {
				t.Fatalf("StrictRect.Rect was not completely false at %d,%d", x, y)
			}
		}
	}

	sr[3][3] = true

	out, err := sr.Outline()
	if err != nil {
		t.Fatalf("expected no error on outline, got %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 outline length, got %v", len(out))
	}
}
