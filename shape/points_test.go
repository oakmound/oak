package shape

import (
	"testing"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

var (
	testPoints = NewPoints(
		intgeom.Point2{1, 1}, intgeom.Point2{2, 1}, intgeom.Point2{3, 1},
		intgeom.Point2{1, 2}, intgeom.Point2{3, 2},
		intgeom.Point2{1, 3}, intgeom.Point2{2, 3}, intgeom.Point2{3, 3},
	)
)

func TestPointsIn(t *testing.T) {
	if !testPoints.In(1, 3, 1, 1) {
		t.Fatalf("1,3 was not in testPoints")
	}
	if testPoints.In(10, 10, 1, 1) {
		t.Fatalf("10,10 was in testPoints")
	}
}

func TestPointsOutline(t *testing.T) {
	testOutline, _ := testPoints.Outline(4, 4)
	if (intgeom.Point2{3, 2}) != testOutline[3] {
		t.Fatalf("expected 3,2 at index 3 in outline, was %v", testOutline[3])
	}
}

func TestPointsRect(t *testing.T) {
	testRect := testPoints.Rect(4, 4)
	if testRect[0][0] {
		t.Fatalf("0,0 should not be set")
	}
	if testRect[2][2] {
		t.Fatalf("2,2 should not be set")
	}
	if !testRect[1][1] {
		t.Fatalf("1,1 should be set")
	}
}
