package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriangulateConvex(t *testing.T) {
	type testCase struct {
		in  []int
		out [][3]int
	}
	testCases := []testCase{
		{
			[]int{0, 1},
			[][3]int{},
		},
		{
			[]int{0},
			[][3]int{},
		},
		{
			[]int{0, 1, 2},
			[][3]int{{0, 1, 2}},
		},
		{
			[]int{0, 1, 2, 3, 4},
			[][3]int{{0, 1, 2}, {0, 2, 3}, {0, 3, 4}},
		},
	}
	for _, tc := range testCases {
		out := TriangulateConvex(tc.in)
		assert.Equal(t, tc.out, out)
	}
}
