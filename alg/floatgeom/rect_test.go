package floatgeom

import (
	"testing"
	"testing/quick"
)

func TestRectConstructors(t *testing.T) {
	if NewRect2(0, 0, 1, 1) != (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}) {
		t.Fatalf("expected %v got %v", NewRect2(0, 0, 1, 1), Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}})
	}
	if NewRect2(1, 1, 0, 0) != (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}) {
		t.Fatalf("expected %v got %v", NewRect2(1, 1, 0, 0), Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}})
	}
	if NewRect2WH(0, 0, 1, 1) != (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}) {
		t.Fatalf("expected %v got %v", NewRect2WH(0, 0, 1, 1), (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}))
	}
	if NewRect2WH(1, 1, -1, -1) != (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}) {
		t.Fatalf("expected %v, got %v", NewRect2WH(1, 1, -1, -1), Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}})
	}
	if (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}) != NewBoundingRect2(Point2{0, 0}, Point2{0, 1}, Point2{1, 0}, Point2{1, 1}) {
		t.Fatalf("expected %v, got %v", (Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}), NewBoundingRect2(Point2{0, 0}, Point2{0, 1}, Point2{1, 0}, Point2{1, 1}))
	}

	if NewRect3(0, 0, 0, 1, 1, 1) != (Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}) {
		t.Fatalf("expected %v, got %v", NewRect3(0, 0, 0, 1, 1, 1), Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}})
	}
	if NewRect3(1, 1, 1, 0, 0, 0) != (Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}) {
		t.Fatalf("expected %v, got %v", NewRect3(1, 1, 1, 0, 0, 0), Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}})
	}
	if NewRect3WH(0, 0, 0, 1, 1, 1) != (Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}) {
		t.Fatalf("expected %v, got %v", NewRect3WH(0, 0, 0, 1, 1, 1), Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}})
	}
	if NewRect3WH(1, 1, 1, -1, -1, -1) != (Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}) {
		t.Fatalf("expected %v, got %v", NewRect3WH(1, 1, 1, -1, -1, -1), Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}})
	}
	if (Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}) != NewBoundingRect3(Point3{0, 0, 0}, Point3{0, .5, 1}, Point3{.5, 1, 0}, Point3{1, 0, .5}) {
		t.Fatalf("expected %v, got %v", Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}, NewBoundingRect3(Point3{0, 0, 0}, Point3{0, .5, 1}, Point3{.5, 1, 0}, Point3{1, 0, .5}))
	}
}

func TestRectAccess(t *testing.T) {
	r2 := NewRect2(0, 1, 2, 3)
	r3 := NewRect3(0, 1, 2, 3, 4, 5)
	if (4.0) != (r2.Area()) {
		t.Fatalf("expected %v got %v", 4.0, r2.Area())
	}
	if (27.0) != (r3.Space()) {
		t.Fatalf("expected %v got %v", 27.0, r3.Space())
	}
	if (2.0) != (r2.W()) {
		t.Fatalf("expected %v got %v", 2.0, r2.W())
	}
	if (2.0) != (r2.H()) {
		t.Fatalf("expected %v got %v", 2.0, r2.H())
	}
	if (3.0) != (r3.W()) {
		t.Fatalf("expected %v got %v", 3.0, r3.W())
	}
	if (3.0) != (r3.H()) {
		t.Fatalf("expected %v got %v", 3.0, r3.H())
	}
	if (3.0) != (r3.D()) {
		t.Fatalf("expected %v got %v", 3.0, r3.D())
	}
	if (1.0) != (r2.Midpoint(0)) {
		t.Fatalf("expected %v got %v", 1.0, r2.Midpoint(0))
	}
	if (2.0) != (r2.Midpoint(1)) {
		t.Fatalf("expected %v got %v", 2.0, r2.Midpoint(1))
	}
	if (1.5) != (r3.Midpoint(0)) {
		t.Fatalf("expected %v got %v", 1.5, r3.Midpoint(0))
	}
	if (2.5) != (r3.Midpoint(1)) {
		t.Fatalf("expected %v got %v", 2.5, r3.Midpoint(1))
	}
	if (3.5) != (r3.Midpoint(2)) {
		t.Fatalf("expected %v got %v", 3.5, r3.Midpoint(2))
	}
	if (8.0) != (r2.Perimeter()) {
		t.Fatalf("expected %v got %v", 8.0, r2.Perimeter())
	}
	if (36.0) != (r3.Margin()) {
		t.Fatalf("expected %v got %v", 36.0, r3.Margin())
	}
}

func TestRect2Contains(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true}

	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		if (e) != (r2.Contains(a)) {
			t.Fatalf("expected %v got %v", e, r2.Contains(a))
		}
	}
}

func TestRect3Contains(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true}

	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		if (e) != (r3.Contains(a)) {
			t.Fatalf("expected %v got %v", e, r3.Contains(a))
		}
	}
}

var (
	r2cases = []Rect2{
		NewRect2(1, 1, 2, 2),
		NewRect2(3, 3, 11, 11),
		NewRect2(11, 11, 12, 12),
	}
	r3cases = []Rect3{
		NewRect3(1, 1, 1, 2, 2, 2),
		NewRect3(3, 3, 3, 11, 11, 11),
		NewRect3(11, 11, 11, 12, 12, 12),
	}
)

func TestRect2Center(t *testing.T) {
	expected := []Point2{
		{1.5, 1.5},
		{7, 7},
		{11.5, 11.5}}
	for i, e := range expected {
		c := r2cases[i]
		if (e) != (c.Center()) {
			t.Fatalf("expected %v got %v", e, c.Center())
		}
	}
}

func TestRect3Center(t *testing.T) {
	expected := []Point3{
		{1.5, 1.5, 1.5},
		{7, 7, 7},
		{11.5, 11.5, 11.5},
	}
	for i, e := range expected {
		c := r3cases[i]
		if (e) != (c.Center()) {
			t.Fatalf("expected %v got %v", e, c.Center())
		}
	}
}

func TestRect2ContainsRect(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true, false, false}

	for i, e := range expected {
		c := r2cases[i]
		if (e) != (r2.ContainsRect(c)) {
			t.Fatalf("expected %v got %v", e, r2.ContainsRect(c))
		}
	}
}

func TestRect3ContainsRect(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true, false, false}

	for i, e := range expected {
		c := r3cases[i]
		if (e) != (r3.ContainsRect(c)) {
			t.Fatalf("expected %v got %v", e, r3.ContainsRect(c))
		}
	}
}

func TestRect2Intersects(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true, true, false}

	for i, e := range expected {
		c := r2cases[i]
		if (e) != (r2.Intersects(c)) {
			t.Fatalf("expected %v got %v", e, r2.Intersects(c))
		}
	}
	r2 = r2.Shift(Point2{3, 3})
	expectedUpdated := []bool{false, true, true}
	for i, e := range expectedUpdated {
		c := r2cases[i]
		if (e) != (r2.Intersects(c)) {
			t.Fatalf("expected %v got %v", e, r2.Intersects(c))
		}
	}
}

func TestRect3Intersects(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true, true, false}

	for i, e := range expected {
		c := r3cases[i]
		if (e) != (r3.Intersects(c)) {
			t.Fatalf("expected %v got %v", e, r3.Intersects(c))
		}
	}
	r3 = r3.Shift(Point3{3, 3, 3})
	expectedUpdated := []bool{false, true, true}
	for i, e := range expectedUpdated {
		c := r3cases[i]
		if (e) != (r3.Intersects(c)) {
			t.Fatalf("expected %v got %v", e, r3.Intersects(c))
		}
	}
}

func TestMaxRectDimensions(t *testing.T) {
	if (2) != (Rect2{}.MaxDimensions()) {
		t.Fatalf("expected %v got %v", 2, Rect2{}.MaxDimensions())
	}
	if (3) != (Rect3{}.MaxDimensions()) {
		t.Fatalf("expected %v got %v", 3, Rect3{}.MaxDimensions())
	}
}

func TestRect2GreaterOf(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []Rect2{
		NewRect2(0, 0, 10, 10),
		NewRect2(0, 0, 11, 11),
		NewRect2(0, 0, 12, 12),
	}

	for i, e := range expected {
		c := r2cases[i]
		if (e) != (r2.GreaterOf(c)) {
			t.Fatalf("expected %v got %v", e, r2.GreaterOf(c))
		}
	}
}

func TestRect3GreaterOf(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []Rect3{
		NewRect3(0, 0, 0, 10, 10, 10),
		NewRect3(0, 0, 0, 11, 11, 11),
		NewRect3(0, 0, 0, 12, 12, 12),
	}

	for i, e := range expected {
		c := r3cases[i]
		if (e) != (r3.GreaterOf(c)) {
			t.Fatalf("expected %v got %v", e, r3.GreaterOf(c))
		}
	}
}

func TestRect3ProjectZ(t *testing.T) {
	if err := quick.Check(projectZHolds, nil); err != nil {
		t.Error(err)
	}
}

func projectZHolds(x1, y1, z1, x2, y2, z2 float64) bool {
	r := NewRect3(x1, y1, z1, x2, y2, z2)
	expected := NewRect2(x1, y1, x2, y2)
	projected := r.ProjectZ()
	return expected == projected
}

func TestRect2Span(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		r := NewRect2WH(1, 1, 9, 9)
		p1 := r.Percentile(1.0)
		if p1 != r.Max {
			t.Errorf("Percentile(1.0) did not return max point: got %v expected %v", p1, r.Max)
		}
		p2 := r.Percentile(0.0)
		if p2 != r.Min {
			t.Errorf("Percentile(0.0) did not return min point: got %v expected %v", p2, r.Min)
		}
		const pollTries = 100
		for i := 0; i < pollTries; i++ {
			if !r.Contains(r.Poll()) {
				t.Fatalf("polled point did not lie within the creating rectangle")
			}
		}
		p3 := r.Clamp(Point2{0, 5})
		if p3 != (Point2{1, 5}) {
			t.Errorf("Clamp(0,5) did not return {1,5}: got %v", p3)
		}
		p4 := r.Clamp(Point2{2, 11})
		if p4 != (Point2{2, 10}) {
			t.Errorf("Clamp(2,11) did not return {2,10}: got %v", p4)
		}
		r2 := r.MulSpan(4)
		if r2 != NewRect2(4, 4, 40, 40) {
			t.Errorf("MulSpan did not return {4,4,40,40}: got %v", r2)
		}
	})
}

func TestRect3Span(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		r := NewRect3WH(1, 1, 1, 9, 9, 9)
		p1 := r.Percentile(1.0)
		if p1 != r.Max {
			t.Errorf("Percentile(1.0) did not return max point: got %v expected %v", p1, r.Max)
		}
		p2 := r.Percentile(0.0)
		if p2 != r.Min {
			t.Errorf("Percentile(0.0) did not return min point: got %v expected %v", p2, r.Min)
		}
		const pollTries = 100
		for i := 0; i < pollTries; i++ {
			if !r.Contains(r.Poll()) {
				t.Fatalf("polled point did not lie within the creating rectangle")
			}
		}
		p3 := r.Clamp(Point3{0, -1, 5})
		if p3 != (Point3{1, 1, 5}) {
			t.Errorf("Clamp(0,-1,5) did not return {1,1,5}: got %v", p3)
		}
		p4 := r.Clamp(Point3{20, 2, 11})
		if p4 != (Point3{10, 2, 10}) {
			t.Errorf("Clamp(20, 2,11) did not return {10,2,10}: got %v", p4)
		}
		r2 := r.MulSpan(4)
		if r2 != NewRect3(4, 4, 4, 40, 40, 40) {
			t.Errorf("MulSpan did not return {4,4,4,40,40,40}: got %v", r2)
		}
	})
}
