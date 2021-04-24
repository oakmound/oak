package ray

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/oakmound/oak/v2/alg"
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/collision"
)

func TestConeCasterSettings(t *testing.T) {
	c1 := NewConeCaster(
		CenterCone(false),
		ConeSpread(180),
		ConeRays(2),
	)
	if c1.CenterCone {
		t.Fatalf("center cone not set to false")
	}
	if c1.Rays != 2 {
		t.Fatalf("cone rays not set to 2")
	}
	if c1.ConeSpread != 180*alg.DegToRad {
		t.Fatalf("cone spread not set to 180")
	}

	c2 := NewConeCaster(
		ConeSpreadRadians(math.Pi/2),
		ConeRays(10),
	)

	if !c2.CenterCone {
		t.Fatalf("center cone not set to true by default")
	}
	if c2.Rays != 10 {
		t.Fatalf("cone rays not set to 10")
	}
	if c2.ConeSpread != math.Pi/2 {
		t.Fatalf("cone spread not set to pi/2")
	}
}

func TestConeCastZeroRays(t *testing.T) {
	const testCt = 100

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < testCt; i++ {
		c1 := NewConeCaster(
			CenterCone(false),
			ConeSpread(180),
			ConeRays(rand.Intn(100)*-1),
		)
		pts := c1.Cast(floatgeom.Point2{}, floatgeom.Point2{})
		if len(pts) != 0 {
			t.Fatalf("more than one point returned from negative rays cone cast")
		}
	}
}

func TestConeCasterScene(t *testing.T) {
	oldDefault := DefaultConeCaster
	defer SetDefaultConeCaster(oldDefault)
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
		coneOpts       []ConeCastOption
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
			},
			// TODO: This was really difficult to set up. This interface needs some
			// tutorial or example visualizing what different setups do
			coneOpts: []ConeCastOption{
				ConeRays(3),
				ConeSpread(90),
			},
			origin: floatgeom.Point2{11, 0},
			target: floatgeom.Point2{11, 1},
			expected: []*collision.Space{
				spaces["101-201"],
				spaces["100-200"],
				// TODO: add option so cone casters don't return duplicate spaces?
				spaces["100-200"],
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			if tc.teardown != nil {
				defer tc.teardown()
			}
			SetDefaultCaster(NewCaster(tc.opts...))
			out := NewConeCaster(tc.coneOpts...).CastTo(tc.origin, tc.target)
			fmt.Println(out)
			if len(out) != len(tc.expected) {
				t.Fatalf("expected length not matched: %v vs %v", len(out), len(tc.expected))
			}
			for i, p := range out {
				if !reflect.DeepEqual(tc.expected[i], p.Zone) {
					t.Fatalf("mismatch at index %d, %v vs %v", i, tc.expected[i], p.Zone)
				}
			}
		})
	}
}

func TestDefaultCasting(t *testing.T) {
	pts := Cast(floatgeom.Point2{}, floatgeom.Point2{})
	if len(pts) != 0 {
		t.Fatalf("Default cone cast did not return empty on empty tree")
	}
	pts = CastTo(floatgeom.Point2{}, floatgeom.Point2{})
	if len(pts) != 0 {
		t.Fatalf("Default cone cast to did not return empty on empty tree")
	}

	pts = ConeCast(floatgeom.Point2{}, floatgeom.Point2{})
	if len(pts) != 0 {
		t.Fatalf("Default cone cast did not return empty on empty tree")
	}
	pts = ConeCastTo(floatgeom.Point2{}, floatgeom.Point2{})
	if len(pts) != 0 {
		t.Fatalf("Default cone cast to did not return empty on empty tree")
	}
}
