package physics

import (
	"testing"
)

type DeltaMass struct {
	Mass
	delta Vector
}

func (dm *DeltaMass) GetDelta() Vector {
	return dm.delta
}

func TestForce(t *testing.T) {
	v := NewForceVector(NewVector(100, 0), 100)
	v2 := DefaultForceVector(NewVector(100, 0), 100)
	if *v.Force != 100.0 {
		t.Fatalf("force vector created with 100 force did not have 100 force")
	}
	if *v2.Force != 10000.0 {
		t.Fatalf("mass to force vector did not scale")
	}

	v3 := NewVector(100, 100).GetForce()
	if *v3.Force != 0.0 {
		t.Fatalf("Non-force vector had force")
	}

	dm := &DeltaMass{
		Mass{100},
		NewVector(100, 0),
	}

	if dm.SetMass(-10) == nil {
		t.Fatalf("set mass to negative did not fail")
	}
	if dm.SetMass(10) != nil {
		t.Fatalf("set mass failed")
	}

	dm2 := &DeltaMass{
		Mass{-10},
		NewVector(0, 0),
	}

	if Push(v3, dm2) == nil {
		t.Fatalf("pushing negative delta mass did not fail")
	}
	if Push(v3, dm) != nil {
		t.Fatalf("pushing positive delta mass failed")
	}

	dm2.Freeze()

	if Push(v3, dm2) != nil {
		t.Fatalf("pushing frozen delta mass failed")
	}

	// Todo: test that pushing results in expected changes
}
