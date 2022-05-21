package ray

import (
	"reflect"
	"testing"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/collision"
)

func TestCasterScene(t *testing.T) {
	SetDefaultCaster(&Caster{
		PointSize:    floatgeom.Point2{.1, .1},
		PointSpan:    1.0,
		CastDistance: 200,
		Tree:         collision.DefaultTree,
	})
	type testCase struct {
		name           string
		setup          func()
		teardown       func()
		opts           []CastOption
		defCaster      *Caster
		expected       []*collision.Space
		origin, target floatgeom.Point2
	}
	spaces := map[string]*collision.Space{
		"100-200": collision.NewFullSpace(0, 0, 10, 10, 100, 200),
		"101-201": collision.NewFullSpace(10, 10, 10, 10, 101, 201),
		"102-202": collision.NewFullSpace(20, 20, 10, 10, 102, 202),
	}
	tree1 := collision.NewTree()
	tcs := []testCase{
		{
			name: "Ignore Filters",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				LimitResults(2),
				IgnoreLabels(100),
				IgnoreIDs(201),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["102-202"],
			},
		}, {
			name: "Accept Filters",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				LimitResults(2),
				AcceptLabels(100),
				AcceptIDs(200),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["100-200"],
			},
		}, {
			name: "StopAtLabel",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				StopAtLabel(100),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["100-200"],
			},
		}, {
			name: "StopAtID",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				StopAtID(201),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["100-200"],
				spaces["101-201"],
			},
		}, {
			name: "StopAtNothing",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				StopAtID(203),
				StopAtLabel(103),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["100-200"],
				spaces["101-201"],
				spaces["102-202"],
			},
		}, {
			name: "Pierce",
			setup: func() {
				tree1.Add(spaces["100-200"])
				tree1.Add(spaces["101-201"])
				tree1.Add(spaces["102-202"])
			},
			teardown: func() {
				tree1.Delete(spaces["100-200"])
				tree1.Delete(spaces["101-201"])
				tree1.Delete(spaces["102-202"])
			},
			opts: []CastOption{
				Tree(tree1),
				CenterPoints(true),
				Distance(50),
				PointSize(floatgeom.Point2{.001, .001}),
				PointSpan(2),
				Pierce(2),
			},
			defCaster: NewCaster(),
			origin:    floatgeom.Point2{5, 5},
			target:    floatgeom.Point2{25, 25},
			expected: []*collision.Space{
				spaces["102-202"],
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			SetDefaultCaster(tc.defCaster)
			c := NewCaster(tc.opts...)
			out := c.CastTo(tc.origin, tc.target)
			if len(out) != len(tc.expected) {
				t.Fatalf("expected length not matched: %v vs %v", len(out), len(tc.expected))
			}
			for i, p := range out {
				if !reflect.DeepEqual(tc.expected[i], p.Zone) {
					t.Fatalf("mismatch at index %d, %v vs %v", i, tc.expected[i], p.Zone)
				}
			}
			if tc.teardown != nil {
				tc.teardown()
			}
		})
	}
}

func TestNewCasterDefaultTree(t *testing.T) {
	DefaultCaster.Tree = nil
	c := NewCaster()
	DefaultCaster.Tree = collision.DefaultTree
	if c.Tree == nil {
		t.Fatal("nil caster tree should have been set to default tree")
	}
}
