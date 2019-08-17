package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntity struct {
	id   CID
	name string
}

func (t *testEntity) Init() CID {
	t.id = NextID(t)
	return t.id
}

func TestScanForEntity(t *testing.T) {
	type testCase struct {
		name       string
		spinup     func()
		check      func(interface{}) bool
		expectedID int
	}
	tcs := []testCase{
		{
			name: "Happy Path",
			spinup: func() {
				t := &testEntity{name: "nova"}
				t.Init()
				t = &testEntity{name: "celeste"}
				t.Init()
			},
			check: func(i interface{}) bool {
				if te, ok := i.(*testEntity); ok {
					return te.name == "celeste"
				}
				return false
			},
			expectedID: 2,
		},
		{
			name:   "Missing Entity",
			spinup: func() {},
			check: func(i interface{}) bool {
				if te, ok := i.(*testEntity); ok {
					return te.name == "celeste"
				}
				return false
			},
			expectedID: -1,
		},
	}

	// This test is not safe for concurrent running, as it
	// modifies and relies on the callers global.
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.spinup()
			id, _ := ScanForEntity(tc.check)
			require.Equal(t, tc.expectedID, id)
			ResetEntities()
		})
	}
}
