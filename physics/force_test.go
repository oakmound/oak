package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, *v.Force, 100.0)
	assert.Equal(t, *v2.Force, 10000.0)

	v3 := NewVector(100, 100).GetForce()
	assert.Equal(t, *v3.Force, 0.0)

	dm := &DeltaMass{
		Mass{100},
		NewVector(100, 0),
	}

	assert.NotNil(t, dm.SetMass(-10))
	dm.SetMass(10)

	dm2 := &DeltaMass{
		Mass{-10},
		NewVector(0, 0),
	}

	assert.NotNil(t, Push(v3, dm2))
	assert.Nil(t, Push(v3, dm))

	dm2.Freeze()

	assert.Nil(t, Push(v3, dm2))

	// Todo: test that pushing results in expected changes
}
