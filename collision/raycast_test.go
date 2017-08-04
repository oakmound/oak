package collision

import (
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
)

func TestRaycasts(t *testing.T) {
	Clear()
	// First, nothing is in the tree, so make sure we get nothing
	vRange := floatrange.NewLinear(3, 359)
	for i := 0; i < 100; i++ {
		assert.Empty(t, RayCast(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll()))
		p := RayCastSingle(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), []event.CID{5, 6})
		assert.True(t, p.IsNil())
		p = RayCastSingleLabels(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), Label(3))
		assert.True(t, p.IsNil())
		p = RayCastSingleIgnoreLabels(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), Label(3))
		assert.True(t, p.IsNil())
		p = RayCastSingleIgnore(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), []event.CID{5, 6}, Label(3))
		assert.True(t, p.IsNil())
		assert.Empty(t, ConeCast(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll()))
		assert.Empty(t, ConeCastSingle(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), []event.CID{5, 6}))
		assert.Empty(t, ConeCastSingleLabels(vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), vRange.Poll(), Label(3)))
	}
}
