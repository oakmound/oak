package collision

import (
	"testing"
)

func TestReactiveSpace(t *testing.T) {
	Clear()
	var triggered bool
	rs1 := NewEmptyReactiveSpace(NewUnassignedSpace(0, 0, 10, 10))
	if rs1 == nil {
		t.Fatalf("reactive space was nil after creation")
	}
	rs2 := NewReactiveSpace(NewUnassignedSpace(5, 5, 10, 10), map[Label]OnHit{
		Label(1): OnIDs(func(id1, id2 int) {
			triggered = true
		}),
	})
	Add(NewLabeledSpace(6, 6, 1, 1, Label(1)))
	<-rs2.CallOnHits()
	if !triggered {
		t.Fatalf("CallOnHits did not trigger reactive space's callback")
	}
	triggered = false

	rs2.Clear()
	<-rs2.CallOnHits()
	if triggered {
		t.Fatalf("CallOnHits triggered reactive space's callback after it was cleared")
	}

	rs1.Add(Label(1), func(*Space, *Space) {
		triggered = true
	})
	<-rs1.CallOnHits()
	if !triggered {
		t.Fatalf("CallOnHits did not trigger reactive space's callback (2)")
	}

	rs1.Remove(Label(1))
	triggered = false
	<-rs1.CallOnHits()
	if triggered {
		t.Fatalf("CallOnHits triggered reactive space's callback after it was removed")
	}
}
