package physics

import (
	"testing"

	"github.com/oakmound/oak/v4/alg"
)

func TestVectorFuncs(t *testing.T) {
	// Constructors
	v := NewVector(1, 1)
	v2 := AngleVector(45)
	v3 := MaxVector(v, v2)
	VectorsEqual(t, v, v3)
	v3 = MaxVector(v2, v)
	VectorsEqual(t, v, v3)

	x, y := 1.0, 1.0
	v20 := PtrVector(&x, &y)
	VectorsEqual(t, v20, v)

	v21 := NewVector32(1.0, 1.0)
	VectorsEqual(t, v, v21)

	// Copy behavior
	v4 := v.Copy()
	VectorsEqual(t, v, v4)
	v5 := Vector{}
	v6 := NewVector(0, 0)
	VectorsEqual(t, v5.Copy(), v6)

	// Magnitude
	if !alg.F64eq(v.Magnitude(), v2.X()+v2.Y()) {
		t.Fatalf("magnitude mismatch (1)")
	}
	if v2.Magnitude() != 1.0 {
		t.Fatalf("magnitude mismatch (2)")
	}

	// Normalize
	v7 := v.Normalize()
	if !alg.F64eq(v7.X(), v2.X()) {
		t.Fatalf("normalize mismatch")
	}
	v8 := v6.Normalize()
	VectorsEqual(t, v6, v8)

	// Zero
	v9 := v4.Zero()
	VectorsEqual(t, v9, v6)

	// Add, Scale
	v10 := NewVector(1, 1).Add(NewVector(1, 1))
	v11 := NewVector(2, 2)
	VectorsEqual(t, v10, v11)
	v10.Scale(.5)
	if v10.X() != 1.0 {
		t.Fatalf("scale mismatch")
	}

	// Rotate / Angle
	v12 := AngleVector(45)
	v13 := AngleVector(90)
	v14 := v12.Rotate(45)
	VectorsEqual(t, v13, v14)
	if v14.Angle() != 90.0 {
		t.Fatalf("angle mismatch")
	}

	// Dot product
	if v14.Dot(v13) != 1.0 {
		t.Fatalf("dot mismatch")
	}

	// Distance
	v15 := NewVector(0, 0)
	v16 := NewVector(0, 10)
	if v15.Distance(v16) != 10.0 {
		t.Fatalf("distance mismatch")
	}

	// Getters
	x, y = v16.GetPos()
	x3 := v16.X()
	y3 := v16.Y()
	x4 := *v16.Xp()
	y4 := *v16.Yp()
	if x != x3 {
		t.Fatalf("x vs x pointer mismatch (1)")
	}
	if x != x4 {
		t.Fatalf("x vs x pointer mismatch (2)")
	}
	if y != y3 {
		t.Fatalf("y vs y pointer mismatch (1)")
	}
	if y != y4 {
		t.Fatalf("y vs y pointer mismatch (2)")
	}

	// Setters
	v17 := v16.SetPos(0, 0)
	VectorsEqual(t, v17, NewVector(0, 0))
	v18 := v17.SetX(1)
	v18 = v18.SetY(1)
	VectorsEqual(t, v18, NewVector(1, 1))
	v19 := v18.ShiftX(-1)
	v19 = v19.ShiftY(-1)
	VectorsEqual(t, v19, NewVector(0, 0))
}

func TestVectorSub(t *testing.T) {
	v := NewVector(0, 0)
	v2 := NewVector(10, 9)
	v.Sub(v2)
	VectorsEqual(t, v, NewVector(-10, -9))
}

func VectorsEqual(t *testing.T, v1, v2 Vector) {
	t.Helper()
	if v1.X() != v2.X() {
		t.Fatalf("x mismatch: %v vs %v", v1.X(), v2.X())
	}
	if v1.Y() != v2.Y() {
		t.Fatalf("y mismatch: %v vs %v", v1.Y(), v2.Y())
	}
	if v1.offX != v2.offX {
		t.Fatalf("xOff mismatch: %v vs %v", v1.offX, v2.offX)
	}
	if v1.offY != v2.offY {
		t.Fatalf("yOff mismatch: %v vs %v", v1.offY, v2.offY)
	}
}
