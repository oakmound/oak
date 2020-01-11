package shape

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/intgeom"
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
				{
					intgeom.Point2{2, 2},
					intgeom.Point2{2, 3},
					intgeom.Point2{3, 3},
					intgeom.Point2{3, 2},
				},
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

func TestBorderHoles(t *testing.T) {
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
				{
					intgeom.Point2{2, 2},
					intgeom.Point2{0, 0},
					intgeom.Point2{0, 1},
					intgeom.Point2{0, 2},
					intgeom.Point2{0, 3},
					intgeom.Point2{0, 4},
					intgeom.Point2{1, 0},
					intgeom.Point2{2, 0},
					intgeom.Point2{3, 0},
					intgeom.Point2{4, 0},
					intgeom.Point2{1, 4},
					intgeom.Point2{2, 4},
					intgeom.Point2{3, 4},
					intgeom.Point2{4, 4},
					intgeom.Point2{4, 1},
					intgeom.Point2{4, 2},
					intgeom.Point2{4, 3},
				},
			},
			5, 5,
		},
	}
	for _, sh := range shapes {
		holes := GetBorderHoles(sh.sh, sh.w, sh.h)
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
