package physics

import (
	"testing"

	"github.com/oakmound/oak/alg"
	"github.com/stretchr/testify/assert"
)

func TestVectorFuncs(t *testing.T) {
	// Constructors
	v := NewVector(1, 1)
	v2 := AngleVector(45)
	v3 := MaxVector(v, v2)
	assert.Equal(t, v, v3)
	v3 = MaxVector(v2, v)
	assert.Equal(t, v, v3)

	// Copy behavior
	v4 := v.Copy()
	assert.Equal(t, v, v4)
	v5 := Vector{}
	v6 := NewVector(0, 0)
	assert.Equal(t, v5.Copy(), v6)

	// Magnitude
	assert.True(t, alg.F64eq(v.Magnitude(), v2.X()+v2.Y()))
	assert.Equal(t, v2.Magnitude(), 1.0)

	// Normalize
	v7 := v.Normalize()
	assert.True(t, alg.F64eq(v7.X(), v2.X()))
	v8 := v6.Normalize()
	assert.Equal(t, v6, v8)

	// Zero
	v9 := v4.Zero()
	assert.Equal(t, v9, v6)

	// Add, Scale
	v10 := NewVector(1, 1).Add(NewVector(1, 1))
	v11 := NewVector(2, 2)
	assert.Equal(t, v10, v11)
	v10.Scale(.5)
	assert.Equal(t, v10.X(), 1.0)

	// Rotate / Angle
	v12 := AngleVector(45)
	v13 := AngleVector(90)
	v14 := v12.Rotate(45)
	assert.Equal(t, v13, v14)
	assert.Equal(t, v14.Angle(), 90.0)

	// Dot product
	assert.Equal(t, v14.Dot(v13), 1.0)

	// Distance
	v15 := NewVector(0, 0)
	v16 := NewVector(0, 10)
	assert.Equal(t, v15.Distance(v16), 10.0)

	// Getters
	x, y := v16.GetPos()
	x2 := v16.GetX()
	x3 := v16.X()
	y2 := v16.GetY()
	y3 := v16.Y()
	x4 := *v16.Xp()
	y4 := *v16.Yp()
	assert.Equal(t, x, x2)
	assert.Equal(t, x, x3)
	assert.Equal(t, x, x4)
	assert.Equal(t, y, y2)
	assert.Equal(t, y, y3)
	assert.Equal(t, y, y4)

	// Setters
	v17 := v16.SetPos(0, 0)
	assert.Equal(t, v17, NewVector(0, 0))
	v18 := v17.SetX(1)
	v18 = v18.SetY(1)
	assert.Equal(t, v18, NewVector(1, 1))
	v19 := v18.ShiftX(-1)
	v19 = v19.ShiftY(-1)
	assert.Equal(t, v19, NewVector(0, 0))
}

func TestVectorSub(t *testing.T) {
	v := NewVector(0, 0)
	v2 := NewVector(10, 9)
	v.Sub(v2)
	assert.Equal(t, v, NewVector(-10, -9))
}
