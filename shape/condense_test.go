package shape

import (
	"reflect"
	"sort"
	"testing"

	"github.com/oakmound/oak/v4/alg/intgeom"
)

func TestCondense(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name     string
		shape    Shape
		w, h     int
		expected []intgeom.Rect2
	}
	tcs := []testCase{
		{
			name:  "Single Rectangle",
			shape: Rectangle,
			w:     10,
			h:     10,
			expected: []intgeom.Rect2{
				intgeom.NewRect2WH(0, 0, 9, 9),
			},
		},
		{
			name: "Double Rectangle",
			shape: StrictRect([][]bool{
				{true, true, false, false, false},
				{true, true, false, true, true},
				{true, true, false, true, true},
				{false, false, false, true, true},
				{false, false, false, false, false},
			}),
			w: 5,
			h: 5,
			expected: []intgeom.Rect2{
				intgeom.NewRect2WH(0, 0, 2, 1),
				intgeom.NewRect2WH(1, 3, 2, 1),
			},
		},
		{
			name: "Double Rectangle 2",
			shape: StrictRect([][]bool{
				{true, false, false, false, false},
				{true, false, false, false, true},
				{true, false, false, false, true},
				{false, false, false, false, true},
				{false, false, false, false, false},
			}),
			w: 5,
			h: 5,
			expected: []intgeom.Rect2{
				intgeom.NewRect2WH(0, 0, 2, 0),
				intgeom.NewRect2WH(1, 4, 2, 0),
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := Condense(tc.shape, tc.w, tc.h)
			sort.Slice(out, func(i, j int) bool {
				if out[i].Min.X() != out[j].Min.X() {
					return out[i].Min.X() < out[j].Min.X()
				}
				return out[i].Min.Y() < out[j].Min.Y()
			})
			if !reflect.DeepEqual(out, tc.expected) {
				t.Fatalf("output rectangles did not match expected")
			}
		})
	}
}
