package ray

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/collision"
	"github.com/stretchr/testify/require"
)

func TestCasterScene(t *testing.T) {
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
	tree1, err := collision.NewTree(2, 20)
	require.Nil(t, err)
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
			require.Equal(t, len(tc.expected), len(out))
			for i, p := range out {
				require.Equal(t, tc.expected[i], p.Zone)
			}
			if tc.teardown != nil {
				tc.teardown()
			}
		})
	}
}
