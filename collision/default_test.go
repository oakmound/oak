package collision

import (
	"testing"
)

func TestDefaultFns(t *testing.T) {
	Clear()
	s := NewUnassignedSpace(0, 0, 10, 10)
	Add(s)
	Remove(s)
	if len(Hits(NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("empty tree hits check should have been empty")
	}

	Add(s)
	err := ShiftSpace(3, 3, s)
	if err != nil {
		t.Fatalf("shift space failed: %v", err)
	}
	if len(Hits(NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("missing hits check should have been empty")
	}
	err = UpdateSpace(3, 3, 10, 10, s)
	if err != nil {
		t.Fatalf("update space failed: %v", err)
	}
	if len(Hits(NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("missing hits check should have been empty")
	}
	err = s.Update(0, 0, 10, 10)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if len(Hits(NewUnassignedSpace(1, 1, 1, 1))) == 0 {
		t.Fatalf("valid hits check should not have been empty")
	}

	Clear()
	if len(Hits(NewUnassignedSpace(1, 1, 1, 1))) != 0 {
		t.Fatalf("cleared hits check should have been empty")
	}

	s = NewLabeledSpace(0, 0, 2, 2, Label(2))
	Add(s)
	if HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)) != nil {
		t.Fatalf("hit label check (1) should have been nil")
	}
	err = s.SetDim(10, 10)
	if err != nil {
		t.Fatalf("SetDim failed: %v", err)
	}
	if HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)) == nil {
		t.Fatalf("hit label check (2) should not have been nil")
	}
	s.UpdateLabel(Label(1))
	if HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)) != nil {
		t.Fatalf("hit label check (3) should have been nil")
	}

}
