package event

import (
	"testing"
)

func TestGetEntityFails(t *testing.T) {
	entity := GetEntity(100)
	if entity != nil {
		t.Fatalf("expected nil entity, got %v", entity)
	}
}
