package shape

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRectangleOutline(t *testing.T) {
	type testCase struct {
		w, h           int
		expectedLength int
		expectedErr    error
	}
	tcs := []testCase{
		{
			w:              2,
			h:              2,
			expectedLength: 4,
		},
		{
			w:              3,
			h:              3,
			expectedLength: 8,
		},
		{
			w:              4,
			h:              4,
			expectedLength: 12,
		},
	}
	for _, tc := range tcs {
		out, err := Rectangle.Outline(tc.w, tc.h)
		if tc.expectedErr == nil {
			require.Nil(t, err)
		} else {
			require.Contains(t, tc.expectedErr.Error(), err.Error())
		}
		require.Equal(t, tc.expectedLength, len(out))
	}
}

func TestToOutline4(t *testing.T) {
	type testCase struct {
		w, h        int
		shape       Shape
		outlineLen  int
		outline4Len int
	}
	tcs := []testCase{
		{
			w:           10,
			h:           10,
			shape:       Heart,
			outlineLen:  24,
			outline4Len: 32,
		},
	}
	for _, tc := range tcs {
		out, _ := tc.shape.Outline(tc.w, tc.h)
		require.Equal(t, tc.outlineLen, len(out))

		out, _ = ToOutline4(tc.shape)(tc.w, tc.h)
		require.Equal(t, tc.outline4Len, len(out))
	}
}

func TestJustInOutline(t *testing.T) {
	type testCase struct {
		w, h        int
		in          func(int, int, ...int) bool
		expectedLen int
		shouldErr   bool
	}
	tcs := []testCase{
		{
			w: 10,
			h: 3,
			in: func(x, y int, sizes ...int) bool {
				return x > 5
			},
			expectedLen: 10,
		},
		{
			w: 3,
			h: 10,
			in: func(x, y int, sizes ...int) bool {
				return x > 5
			},
			shouldErr: true,
		},
		{
			w: 2,
			h: 2,
			in: func(x, y int, sizes ...int) bool {
				return x == 0 && y == 0
			},
			expectedLen: 1,
		},
	}
	for _, tc := range tcs {
		out, err := JustIn(tc.in).Outline(tc.w, tc.h)
		if tc.shouldErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		require.Equal(t, tc.expectedLen, len(out))
	}
}
