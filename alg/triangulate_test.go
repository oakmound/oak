package alg

import (
	"reflect"
	"strconv"
	"testing"
)

func TestTriangulateConvex(t *testing.T) {
	t.Parallel()
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
	for i, tc := range testCases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			out := TriangulateConvex(tc.in)
			if !reflect.DeepEqual(tc.out, out) {
				t.Fatalf("expected %v got %v", tc.out, out)
			}
		})
	}
}
