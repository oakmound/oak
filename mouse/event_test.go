package mouse

import (
	"testing"

	"github.com/oakmound/oak/v4/collision"
)

func TestEventConversions(t *testing.T) {
	me := NewEvent(1.0, 1.0, ButtonLeft, Drag)
	s := me.ToSpace()
	Add(collision.NewUnassignedSpace(1.0, 1.0, .1, .1))
	if len(Hits(s)) == 0 {
		t.Fatalf("expected hits to catch unassigned space")
	}
}
