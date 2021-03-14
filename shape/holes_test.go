package shape

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/intgeom"
)

func TestHoles(t *testing.T) {
	t.Parallel()
	shapes := []struct {
		name string
		sh   Shape
		out  [][]intgeom.Point2
		w, h int
	}{
		{
			name: "3x3",
			sh: StrictRect([][]bool{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			}),
			out: [][]intgeom.Point2{
				{intgeom.Point2{1, 1}},
			},
			w: 3,
			h: 3,
		}, {
			name: "3x3 with extra border",
			sh: StrictRect([][]bool{
				{false, false, false, false, false},
				{false, true, true, true, false},
				{false, true, false, true, false},
				{false, true, true, true, false},
				{false, false, false, false, false},
			}),
			out: [][]intgeom.Point2{
				{intgeom.Point2{2, 2}},
			},
			w: 5,
			h: 5,
		}, {
			name: "4x4 with extra border",
			sh: StrictRect([][]bool{
				{false, false, false, false, false, false},
				{false, true, true, true, true, false},
				{false, true, false, false, true, false},
				{false, true, false, false, true, false},
				{false, true, true, true, true, false},
				{false, false, false, false, false, false},
			}),
			out: [][]intgeom.Point2{
				{
					{2, 2},
					{2, 3},
					{3, 3},
					{3, 2},
				},
			},
			w: 6,
			h: 6,
		},
	}
	for _, sh := range shapes {
		sh := sh
		t.Run(sh.name, func(t *testing.T) {
			t.Parallel()
			holes := GetHoles(sh.sh, sh.w, sh.h)
			found := map[intgeom.Point2]bool{}
			for _, group := range holes {
				for _, p1 := range group {
					found[p1] = true
				}
			}
			for _, col := range sh.out {
				for _, p := range col {
					if !found[p] {
						t.Fatalf("unexpected point found at %v", p)
					}
					delete(found, p)
				}
			}
			if len(found) != 0 {
				t.Fatalf("not all points found")
			}
		})
	}
}

func TestBorderHoles(t *testing.T) {
	shapes := []struct {
		name string
		sh   Shape
		out  [][]intgeom.Point2
		w, h int
	}{
		{
			name: "3x3",
			sh: StrictRect([][]bool{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			}),
			out: [][]intgeom.Point2{
				{{1, 1}},
			},
			w: 3,
			h: 3,
		}, {
			name: "3x3 with extra border",
			sh: StrictRect([][]bool{
				{false, false, false, false, false},
				{false, true, true, true, false},
				{false, true, false, true, false},
				{false, true, true, true, false},
				{false, false, false, false, false},
			}),
			out: [][]intgeom.Point2{
				{
					{2, 2},
					{0, 0},
					{0, 1},
					{0, 2},
					{0, 3},
					{0, 4},
					{1, 0},
					{2, 0},
					{3, 0},
					{4, 0},
					{1, 4},
					{2, 4},
					{3, 4},
					{4, 4},
					{4, 1},
					{4, 2},
					{4, 3},
				},
			},
			w: 5,
			h: 5,
		},
	}
	for _, sh := range shapes {
		sh := sh
		t.Run(sh.name, func(t *testing.T) {
			t.Parallel()
			holes := GetBorderHoles(sh.sh, sh.w, sh.h)
			found := map[intgeom.Point2]bool{}
			for _, group := range holes {
				for _, p1 := range group {
					found[p1] = true
				}
			}
			for _, col := range sh.out {
				for _, p := range col {
					if !found[p] {
						t.Fatalf("unexpected point found at %v", p)
					}
					delete(found, p)
				}
			}
			if len(found) != 0 {
				t.Fatalf("not all points found")
			}
		})
	}
}
