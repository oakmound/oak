package ray

import (
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/collision"
	"github.com/stretchr/testify/assert"
)

func TestRaycasts(t *testing.T) {
	collision.DefTree.Clear()
	vRange := floatrange.NewLinear(3, 359)
	tests := 100
	for i := 0; i < tests; i++ {
		p1 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		p2 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		assert.Empty(t, Cast(p1, p2))
		assert.Empty(t, CastTo(p1, p2))
		assert.Empty(t, ConeCast(p1, p2))
		assert.Empty(t, ConeCastTo(p1, p2))
	}
}
