package shape

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/stretchr/testify/require"
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
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := Condense(tc.shape, tc.w, tc.h)
			require.Equal(t, tc.expected, out)
		})
	}
}
