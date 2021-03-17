package event

import (
	"testing"
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

	// This test is not safe for t.Parallel, as it
	// modifies and relies on the callers global.
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.spinup()
			id, _ := ScanForEntity(tc.check)
			if tc.expectedID != id {
				t.Fatalf("expected id %v, got %v", tc.expectedID, id)
			}
			ResetEntities()
		})
	}
}

func TestGetEntityFails(t *testing.T) {
	entity := GetEntity(100)
	if entity != nil {
		t.Fatalf("expected nil entity, got %v", entity)
	}
}
