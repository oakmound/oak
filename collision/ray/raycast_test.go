package ray

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/alg/range/floatrange"
	"github.com/oakmound/oak/v2/collision"
)

func TestEmptyRaycasts(t *testing.T) {
	t.Skip()
	collision.DefaultTree.Clear()
	vRange := floatrange.NewLinear(3, 359)
	tests := 100
	for i := 0; i < tests; i++ {
		p1 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		p2 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		if len(Cast(p1, p2)) != 0 {
			t.Fatalf("cast found a point in the empty tree")
		}
		if len(CastTo(p1, p2)) != 0 {
			t.Fatalf("cast to found a point in the empty tree")
		}
		if len(ConeCast(p1, p2)) != 0 {
			t.Fatalf("cone cast found a point in the empty tree")
		}
		if len(ConeCastTo(p1, p2)) != 0 {
			t.Fatalf("cone cast to found a point in the empty tree")
		}
	}
}
