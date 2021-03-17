package mouse

import (
	"testing"

	"github.com/oakmound/oak/v2/collision"
)

func TestDefaultFunctions(t *testing.T) {
	Clear()
	s := collision.NewUnassignedSpace(0, 0, 10, 10)
	Add(s)
	Remove(s)
	if len(Hits(collision.NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("expected empty tree to have no contents")
	}

	Add(s)
	if ShiftSpace(3, 3, s) != nil {
		t.Fatalf("shift space failed")
	}
	if len(Hits(collision.NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("hit away from space should not collide with space")
	}

	if UpdateSpace(0, 0, 10, 10, s) != nil {
		t.Fatalf("update space failed")
	}
	if len(Hits(collision.NewUnassignedSpace(1, 1, 1, 1))) == 0 {
		t.Fatalf("hit on space should collide")
	}

	Clear()
	if len(Hits(collision.NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("expected cleared tree to have no contents")
	}

	s = collision.NewLabeledSpace(0, 0, 10, 10, collision.Label(2))
	Add(s)
	if HitLabel(collision.NewUnassignedSpace(1, 1, 1, 1), collision.Label(2)) == nil {
		t.Fatalf("hit label missed")
	}
}
