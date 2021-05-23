package shape

import (
	"fmt"
	"testing"
)

func TestRectangleOutline(t *testing.T) {
	t.Parallel()
	type testCase struct {
		w, h           int
		expectedLength int
		shouldErr      bool
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
		tc := tc
		t.Run(fmt.Sprintf("%dx%d", tc.w, tc.h), func(t *testing.T) {
			t.Parallel()
			out, err := Rectangle.Outline(tc.w, tc.h)
			if !tc.shouldErr {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected non-nil error, got nil")
				}
			}
			if tc.expectedLength != len(out) {
				t.Fatalf("expected length of output %v, got %v", tc.expectedLength, len(out))
			}
		})
	}
}

func TestToOutline4(t *testing.T) {
	t.Parallel()
	type testCase struct {
		w, h           int
		shape          Shape
		outlineLength  int
		outline4Length int
	}
	tcs := []testCase{
		{
			w:              10,
			h:              10,
			shape:          Heart,
			outlineLength:  24,
			outline4Length: 32,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("%dx%d", tc.w, tc.h), func(t *testing.T) {
			t.Parallel()
			out, _ := tc.shape.Outline(tc.w, tc.h)
			if tc.outlineLength != len(out) {
				t.Fatalf("expected length of output %v, got %v", tc.outlineLength, len(out))
			}

			out, _ = ToOutline4(tc.shape)(tc.w, tc.h)
			if tc.outline4Length != len(out) {
				t.Fatalf("expected length of output 4 %v, got %v", tc.outline4Length, len(out))
			}
		})
	}
}

func TestJustInOutline(t *testing.T) {
	t.Parallel()
	type testCase struct {
		w, h           int
		in             func(int, int, ...int) bool
		expectedLength int
		shouldErr      bool
	}
	tcs := []testCase{
		{
			w: 10,
			h: 3,
			in: func(x, y int, sizes ...int) bool {
				return x > 5
			},
			expectedLength: 10,
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
			expectedLength: 1,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("%dx%d", tc.w, tc.h), func(t *testing.T) {
			t.Parallel()
			out, err := JustIn(tc.in).Outline(tc.w, tc.h)
			if !tc.shouldErr {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected non-nil error, got nil")
				}
			}
			if tc.expectedLength != len(out) {
				t.Fatalf("expected length of output %v, got %v", tc.expectedLength, len(out))
			}
		})
	}
}
