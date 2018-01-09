package shape

import (
	"testing"

	"github.com/oakmound/oak/alg/intgeom"
	"github.com/stretchr/testify/assert"
)

func TestHoles(t *testing.T) {
	shapes := []struct {
		sh   Shape
		out  [][]intgeom.Point
		w, h int
	}{
		{
			StrictRect([][]bool{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			}),
			[][]intgeom.Point{
				{intgeom.Point{1, 1}},
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
			[][]intgeom.Point{
				{intgeom.Point{2, 2}},
			},
			5, 5,
		},
	}
	for _, sh := range shapes {
		assert.Equal(t, sh.out, GetHoles(sh.sh, sh.w, sh.h))
	}
}
