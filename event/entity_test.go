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

func TestGetEntityFails(t *testing.T) {
	entity := GetEntity(100)
	if entity != nil {
		t.Fatalf("expected nil entity, got %v", entity)
	}
}
