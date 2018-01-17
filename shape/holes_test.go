package shape

import (
	"testing"

	"github.com/oakmound/oak/alg/intgeom"
	"github.com/stretchr/testify/assert"
)

func TestHoles(t *testing.T) {
	shapes := []struct {
		sh   Shape
		out  [][]intgeom.Point2
		w, h int
	}{
		{
			StrictRect([][]bool{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			}),
			[][]intgeom.Point2{
				{intgeom.Point2{1, 1}},
			},
			3, 3,
		}, {
			StrictRect([][]bool{
				{false, false, false, false, false},
				{false, true, true, true, false},
				{false, true, false, true, false},
				{false, true, true, true, false},
				{false, false, false, false, false},
			}),
			[][]intgeom.Point2{
				{intgeom.Point2{2, 2}},
			},
			5, 5,
		}, {
			StrictRect([][]bool{
				{false, false, false, false, false, false},
				{false, true, true, true, true, false},
				{false, true, false, false, true, false},
				{false, true, false, false, true, false},
				{false, true, true, true, true, false},
				{false, false, false, false, false, false},
			}),
			[][]intgeom.Point2{
				{intgeom.Point2{2, 2},
					intgeom.Point2{2, 3},
					intgeom.Point2{3, 3},
					intgeom.Point2{3, 2}},
			},
			6, 6,
		},
	}
	for _, sh := range shapes {
		holes := GetHoles(sh.sh, sh.w, sh.h)
		found := map[intgeom.Point2]bool{}
		for _, group := range holes {
			for _, p1 := range group {
				found[p1] = true
			}
		}
		for _, col := range sh.out {
			for _, p := range col {
				assert.True(t, found[p])
				delete(found, p)
			}
		}
		assert.Empty(t, found)
	}
}
